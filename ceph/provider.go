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
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cephx username to use to connect to Ceph (i.e.: client.admin).",
			},
			"cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Ceph cluster to use.",
			},
			"keyring": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The actual keyring (not a path to a file), to use to connect to Ceph. " +
					"Using this ignore `config_path` and you must also specify `mon_host`",
			},
			"mon_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List of mon to connect to Ceph. This is only used with `keyring`, otherwise it is ignored.",
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
		Username:   d.Get("username").(string),
		Cluster:    d.Get("cluster").(string),
		Keyring:    d.Get("keyring").(string),
		MonHost:    d.Get("mon_host").(string),
	}

	return config, nil
}
