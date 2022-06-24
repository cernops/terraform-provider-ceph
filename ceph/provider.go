package ceph

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"config_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to the ceph config",
			},
			"entity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cephx entity to use to connect to Ceph (i.e.: client.admin).",
			},
			"cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Ceph cluster to use.",
				Default:     "ceph",
			},
			"keyring": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The actual keyring (not a path to a file) to use to connect to Ceph.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The actual key (not a path to a file) to use to connect to Ceph.",
			},
			"mon_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List of mon to connect to Ceph.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ceph_wait_online": dataSourceWaitOnline(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ceph_auth": resourceAuth(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		ConfigPath: d.Get("config_path").(string),
		Entity:     d.Get("entity").(string),
		Cluster:    d.Get("cluster").(string),
		Keyring:    d.Get("keyring").(string),
		Key:        d.Get("key").(string),
		MonHost:    d.Get("mon_host").(string),
	}

	return config, nil
}
