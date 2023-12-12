package rustack_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func marshalTagNames(tags []rustack.Tag) []interface{} {
	convertedTags := make([]interface{}, len(tags))
	for i, tag := range tags {
		convertedTags[i] = tag.Name
	}
	return convertedTags
}

func unmarshalTagNames(tags interface{}) []rustack.Tag {
	tagList := tags.(*schema.Set).List()
	resultTags := make([]rustack.Tag, len(tagList))
	for i, tag := range tagList {
		resultTags[i] = rustack.Tag{Name: tag.(string)}
	}
	return resultTags
}

func newTagNamesResourceSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:        schema.TypeString,
			Description: "name of the Tag",
		},
		Description: description,
	}
}

func marshalTags(tags []rustack.Tag) []map[string]interface{} {
	convertedTags := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		convertedTags[i]["id"] = map[string]interface{}{"id": tag.ID, "name": tag.Name}
	}
	return convertedTags
}

func newTagsDatasourceSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: Arguments{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "id of the Tag",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "name of the Tag",
				},
			},
		},
		Description: description,
	}
}
