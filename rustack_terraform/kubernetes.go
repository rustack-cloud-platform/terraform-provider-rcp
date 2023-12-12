package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextGetKubernetes() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the kubernetes",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the kubernetes",
		},
	})
}

func (args *Arguments) injectContextKubernetesById() {
	args.merge(Arguments{
		"kubernetes_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the kubernetes",
		},
	})
}

func (args *Arguments) injectCreateKubernetes() {
	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 100),
			),
			Description: "name of the kubernetes",
		},
		"platform": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "the number of virtual cpus",
		},
		"node_cpu": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 128),
			Description:  "the number of virtual cpus",
		},
		"node_ram": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "memory of the kubernetes in gigabytes",
		},
		"floating": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "enable floating ip for the kubernetes",
		},
		"floating_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "floating ip for the kubernetes. May be omitted",
		},
		"node_disk_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "size in gb for the vms disk attached to kubernetes.",
		},
		"nodes_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "count of vms attached to kubernetes",
		},
		"user_public_key_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "pub key id for vms attached to kubernetes.",
		},
		"node_storage_profile_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "storage_profile_id for vms disks attached to kubernetes.",
		},
		"vms": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			MinItems:    1,
			Description: "List of Vms connected to the kubernetes",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"dashboard_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Kubernetes dashboard url",
		},
		"tags": newTagNamesResourceSchema("tags of the Kubernetes"),
	})
}

func (args *Arguments) injectResultKubernetes() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the kubernetes",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the kubernetes",
		},
		"node_cpu": {
			Type:     schema.TypeInt,
			Computed: true,

			Description: "the number of virtual cpus",
		},
		"node_ram": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "memory of the kubernetes in gigabytes",
		},
		"template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Template",
		},
		"floating": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "enable floating ip for the kubernetes",
		},
		"floating_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "floating ip for the kubernetes. May be omitted",
		},
		"node_disk_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "size in gb for the vms disk attached to kubernetes.",
		},
		"nodes_count": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "count of vms attached to kubernetes",
		},
		"user_public_key_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "pub key id for vms attached to kubernetes.",
		},
		"node_storage_profile_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "storage_profile_id for vms disks attached to kubernetes.",
		},
		"vms": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "List of Vms connected to the kubernetes",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"dashboard_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kubernetes dashboard url",
		},
	})
}

func (args *Arguments) injectResultListKubernetes() {
	s := Defaults()
	s.injectResultKubernetes()

	args.merge(Arguments{
		"kubernetess": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
