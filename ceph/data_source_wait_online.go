package ceph

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	maxDuration   = time.Hour
	sleepDuration = time.Minute
)

func dataSourceWaitOnline() *schema.Resource {
	return &schema.Resource{
		Description: "This dummy resource is waiting to Ceph to be online for up to 1 hour. " +
			"This is useful for example on a boostrap procedure.",
		ReadContext: dataSourceWaitOnlineRead,
		Schema:      map[string]*schema.Schema{},
	}
}

func dataSourceWaitOnlineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	startTime := time.Now()

	for time.Since(startTime) < maxDuration {
		_, err := config.GetCephConnection()
		log.Printf("[DEBUG] Cannot connect to Ceph on ceph_wait_online: %s", err)
		if err == nil {
			log.Printf("[DEBUG] Ceph online on ceph_wait_online")
			return nil
		}
		time.Sleep(sleepDuration)
	}

	return diag.Errorf("Error time out after trying to connect to Ceph many times")
}
