package rustack_terraform

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func resourceRustackPaasService() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceRustackPaasServiceRead,
		CreateContext: resourceRustackPaasServiceCreate,
		DeleteContext: resourceRustackPaasServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of Paas Service at ESU",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of Paas Service at ESU",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of Project",
			},
			"paas_service_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "id of Paas Template",
			},
			"paas_service_inputs": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "inputs of Paas Service as JSON object",
			},
			// "paas_service_inputs": {
			// 	Type:        schema.TypeMap,
			// 	Required:    true,
			// 	ForceNew:    true,
			// 	StateFunc:   func(i interface{}) string { return "" },
			// 	Description: "inputs of Paas Service",
			// },
		},
	}
}

func resourceRustackPaasServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).rustackManager()
	service, err := manager.GetPaasService(d.Get("id").(string))
	if err != nil {
		if err.(*rustack.RustackApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error getting Paas Service: %s", err)
		}
	}
	d.Set("id", service.ID)
	d.Set("name", service.Name)
	d.Set("project_id", service.Project.ID)
	d.Set("paas_service_id", service.PaasServiceID)
	inputsString, err := json.Marshal(service.Inputs)
	if err != nil {
		return diag.Errorf("Error marshalling Paas Service inputs: %s", err)
	}
	d.Set("paas_service_inputs", inputsString)
	d.SetId(service.ID)
	return nil
}

func resourceRustackPaasServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	var inputs map[string]interface{}
	inputsString := d.Get("paas_service_inputs").(string)
	if err := json.Unmarshal([]byte(inputsString), &inputs); err != nil {
		return diag.Errorf("Error parsing Paas Service inputs: %s", err)
	}
	service := &rustack.PaasService{
		Name: d.Get("name").(string),
		Project: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID: d.Get("project_id").(string),
		},
		PaasServiceID: d.Get("paas_service_id").(int),
		Inputs:        inputs,
	}
	if err := manager.CreatePaasService(service); err != nil {
		return diag.Errorf("Error creating Paas Service: %s", err)
	}
	d.Set("id", service.ID)
	return resourceRustackPaasServiceRead(ctx, d, meta)
}

func resourceRustackPaasServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	err := manager.DeletePaasService(d.Get("id").(string))
	if err != nil {
		return diag.Errorf("Error deleting Paas Service: %s", err)
	}
	return nil
}
