package rustack_terraform

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func (args *Arguments) injectContextPublicKeyById() {
	args.merge(Arguments{
		"pub_key_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the Public Key",
		},
	})
}

func (args *Arguments) injectContextGetPublicKey() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the Public Key",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the Public Key",
		},
	})
}

func (args *Arguments) injectResultPublicKey() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Public Key",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Public Key",
		},
		"fingerprint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "fingerprint of the  Public Key",
		},
		"public_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "public_key of the Public Key instance",
		},
	})
}

func (args *Arguments) injectResultListPublicKeys() {
	s := Defaults()
	s.injectResultPublicKey()

	args.merge(Arguments{
		"public_keys": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
