package ceph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type authResponse struct {
	Entity string            `json:"entity"`
	Key    string            `json:"key"`
	Caps   map[string]string `json:"caps"`
}

func resourceAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthCreate,
		ReadContext:   resourceAuthRead,
		UpdateContext: resourceAuthUpdate,
		DeleteContext: resourceAuthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The entity name (i.e.: client.admin)",
			},

			"caps": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The caps wanted for the entity",
			},

			"keyring": {
				Type:        schema.TypeString,
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

const clientKeyringFormat = `
[%s]
%s
`

func setResourceData(d *schema.ResourceData, authResponse authResponse) diag.Diagnostics {
	if err := d.Set("key", authResponse.Key); err != nil {
		return diag.Errorf("Unable to set key: %s", err)
	}

	keyring := fmt.Sprintf(clientKeyringFormat, authResponse.Entity, authResponse.Key)
	if err := d.Set("keyring", keyring); err != nil {
		return diag.Errorf("Unable to set keyring: %s", err)
	}

	return nil
}

func resourceAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*Config).GetCephConnection()
	if err != nil {
		return diag.Errorf("Unable to connect to Ceph: %s", err)
	}
	entity := d.Get("entity").(string)

	command, err := json.Marshal(map[string]interface{}{
		"prefix": "auth get-or-create",
		"format": "json",
		"entity": entity,
	})
	if err != nil {
		return diag.Errorf("Unable resource_auth unable to create get-or-create JSON command: %s", err)
	}

	buf, _, err := conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error resource_auth on get-or-create command: %s", err)
	}

	var authResponse authResponse
	err = json.Unmarshal(buf, &authResponse)
	if err != nil {
		return diag.Errorf("Error unmarshal on get-or-create response: %s", err)
	}

	d.SetId(entity)
	return setResourceData(d, authResponse)
}

func resourceAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Unable resource_auth unable to create get JSON command: %s", err)
	}

	buf, _, err := conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error resource_auth on get command: %s", err)
	}

	var authResponse authResponse
	err = json.Unmarshal(buf, &authResponse)
	if err != nil {
		return diag.Errorf("Error unmarshal on get-or-create response: %s", err)
	}

	return setResourceData(d, authResponse)
}

func resourceAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceAuthCreate(ctx, d, meta)
}

func resourceAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*Config).GetCephConnection()
	if err != nil {
		return diag.Errorf("Unable to connect to Ceph: %s", err)
	}
	entity := d.Get("entity").(string)

	command, err := json.Marshal(map[string]interface{}{
		"prefix": "auth rm",
		"format": "json",
		"entity": entity,
	})
	if err != nil {
		return diag.Errorf("Unable resource_auth unable to create delete JSON command: %s", err)
	}

	_, _, err = conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error resource_auth on rm command: %s", err)
	}

	return nil
}
