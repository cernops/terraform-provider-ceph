package ceph

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAuth() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows you to get information about a ceph client.",
		ReadContext: dataSourceAuthRead,

		Schema: map[string]*schema.Schema{
			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The entity name (i.e.: client.admin)",
			},

			"caps": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "The caps of the entity",
			},

			"keyring": {
				Type: schema.TypeString,

				Computed:    true,
				Description: "The cephx keyring of the entity",
			},

			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cephx key of the entity",
			},
		},
	}
}

func dataSourceAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*Config).GetCephConnection()
	if err != nil {
		return diag.Errorf("Unable to connect to Ceph: %s", err)
	}
	entity := d.Get("entity").(string)

	command, err := json.Marshal(map[string]interface{}{
		"prefix": "auth get",
		"format": "json",
		"entity": entity,
	})
	if err != nil {
		return diag.Errorf("Error data_source_auth unable to create get JSON command: %s", err)
	}

	buf, _, err := conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error data_source_auth on get command: %s", err)
	}

	var authResponses []authResponse
	err = json.Unmarshal(buf, &authResponses)
	if err != nil {
		return diag.Errorf("Error data_source auth unmarshal on response: %s", err)
	}

	d.SetId(entity)
	return setAuthResourceData(d, authResponses)
}
