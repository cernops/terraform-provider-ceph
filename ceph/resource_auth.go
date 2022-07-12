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
		Description:   "This resource allows you to create a ceph client and retrieve his key and/or keyring.",
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
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The caps of the entity",
			},

			"keyring": {
				Type: schema.TypeString,

				Computed:    true,
				Sensitive:   true,
				Description: "The cephx keyring of the entity",
			},

			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The cephx key of the entity",
			},
		},
	}
}

const clientKeyringFormat = `[%s]
	key = %s
`

func setAuthResourceData(d *schema.ResourceData, authResponses []authResponse) diag.Diagnostics {
	if len(authResponses) == 0 {
		return diag.Errorf("No data returned by ceph auth command")
	}
	if err := d.Set("key", authResponses[0].Key); err != nil {
		return diag.Errorf("Unable to set key: %s", err)
	}

	keyring := fmt.Sprintf(clientKeyringFormat, authResponses[0].Entity, authResponses[0].Key)
	if err := d.Set("keyring", keyring); err != nil {
		return diag.Errorf("Unable to set keyring: %s", err)
	}
	if err := d.Set("caps", authResponses[0].Caps); err != nil {
		return diag.Errorf("Unable to set caps: %s", err)
	}

	return nil
}

func toCapsArray(caps map[string]interface{}) []string {
	var ret []string

	for key, val := range caps {
		ret = append(ret, key)
		ret = append(ret, val.(string))
	}

	return ret
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
		"caps":   toCapsArray(d.Get("caps").(map[string]interface{})),
	})
	if err != nil {
		return diag.Errorf("Error resource_auth unable to create get-or-create JSON command: %s", err)
	}

	buf, _, err := conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error resource_auth on get-or-create command: %s", err)
	}

	var authResponses []authResponse
	err = json.Unmarshal(buf, &authResponses)
	if err != nil {
		return diag.Errorf("Error resource_auth unmarshal on get-or-create response: %s", err)
	}

	d.SetId(entity)
	return setAuthResourceData(d, authResponses)
}

func resourceAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*Config).GetCephConnection()
	if err != nil {
		return diag.Errorf("Unable to connect to Ceph: %s", err)
	}
	entity := d.Id()

	command, err := json.Marshal(map[string]interface{}{
		"prefix": "auth get",
		"format": "json",
		"entity": entity,
	})
	if err != nil {
		return diag.Errorf("Error resource_auth unable to create get JSON command: %s", err)
	}

	buf, _, err := conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error resource_auth on get command: %s", err)
	}

	var authResponses []authResponse
	err = json.Unmarshal(buf, &authResponses)
	if err != nil {
		return diag.Errorf("Error resource_auth unmarshal on get response: %s", err)
	}

	if err := d.Set("entity", entity); err != nil {
		return diag.Errorf("Unable to set entity: %s", err)
	}
	return setAuthResourceData(d, authResponses)
}

func resourceAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(*Config).GetCephConnection()
	if err != nil {
		return diag.Errorf("Unable to connect to Ceph: %s", err)
	}
	entity := d.Get("entity").(string)

	command, err := json.Marshal(map[string]interface{}{
		"prefix": "auth caps",
		"format": "json",
		"entity": entity,
		"caps":   toCapsArray(d.Get("caps").(map[string]interface{})),
	})
	if err != nil {
		return diag.Errorf("Error resource_auth unable to create caps JSON command: %s", err)
	}

	_, _, err = conn.MonCommand(command)
	if err != nil {
		return diag.Errorf("Error resource_auth on caps command: %s", err)
	}

	return resourceAuthRead(ctx, d, meta)
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
