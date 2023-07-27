package rustack_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func dataSourceRustackPublicKey() *schema.Resource {
	args := Defaults()
	args.injectResultPublicKey()
	args.injectContextAccountById()
	args.injectContextGetPublicKey() // override name

	return &schema.Resource{
		ReadContext: dataSourceRustackPublicKeyRead,
		Schema:      args,
	}
}

func dataSourceRustackPublicKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting PublicKey: %s", err)
	}
	var targetPublicKey *rustack.PubKey
	if target == "id" {
		targetPublicKey, err = manager.GetPublicKey(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting PublicKey: %s", err)
		}
	} else {
		targetPublicKey, err = GetPubKeyByName(d, manager)
		if err != nil {
			return diag.Errorf("Error getting PublicKey: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":   targetPublicKey.ID,
		"name": targetPublicKey.Name,
		"public_key": targetPublicKey.Fingerprint,
		"fingerprint": targetPublicKey.PublicKey,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetPublicKey.ID)
	return nil
}
