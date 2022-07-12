package ceph

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWaitOnline() *schema.Resource {
	return &schema.Resource{
		Description: "This dummy resource is waiting to Ceph to be online at creation time for up to 1 hour. " +
			"This is useful for example on a boostrap procedure.",
		CreateContext: resourceWaitOnlineCreate,
		ReadContext:   resourceWaitOnlineRead,
		DeleteContext: resourceWaitOnlineDummy,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "That's a workaround to actually have an id, set this to something unique (i.e.: the cluster name).",
			},

			"online": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the cluster is online, only checked at creationg time (always true)",
			},
		},
	}
}

func resourceWaitOnlineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	log.Printf("[DEBUG] Ceph starting ceph_wait_online")

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := config.GetCephConnection()
		if err == nil {
			log.Printf("[DEBUG] Ceph online on ceph_wait_online")
			d.SetId(d.Get("cluster_name").(string))
			if err := d.Set("online", true); err != nil {
				return resource.NonRetryableError(fmt.Errorf("Unable to set online: %s", err))
			}
			return nil
		}

		log.Printf("[DEBUG] Cannot connect to Ceph on ceph_wait_online: %s", err)
		return resource.RetryableError(err)
	})

	return diag.FromErr(err)
}

func resourceWaitOnlineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := d.Set("cluster_name", d.Id()); err != nil {
		return diag.Errorf("Unable to set cluster_name: %s", err)
	}
	return nil
}

func resourceWaitOnlineDummy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
