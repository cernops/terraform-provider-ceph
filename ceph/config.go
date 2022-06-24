package ceph

import (
	"io/ioutil"
	"os"

	"github.com/ceph/go-ceph/rados"
)

type Config struct {
	ConfigPath string
	Entity     string
	Cluster    string
	Keyring    string
	Key        string
	MonHost    string

	RadosConn *rados.Conn
}

func (config *Config) GetCephConnection() (*rados.Conn, error) {
	var conn *rados.Conn
	var err error

	if config.RadosConn != nil {
		return config.RadosConn, nil
	}

	if config.Entity != "" {
		conn, err = rados.NewConnWithClusterAndUser(config.Cluster, config.Entity)
	} else {
		conn, err = rados.NewConn()
	}
	if err != nil {
		return nil, err
	}

	if config.ConfigPath != "" {
		if err = conn.ReadConfigFile(config.ConfigPath); err != nil {
			return nil, err
		}
	} else {
		conn.ReadDefaultConfigFile() //nolint:golint,errcheck
	}

	if config.MonHost != "" {
		if err = conn.SetConfigOption("mon_host", config.MonHost); err != nil {
			return nil, err
		}
	}

	if config.Key != "" {
		if err = conn.SetConfigOption("key", config.Key); err != nil {
			return nil, err
		}
	}

	if config.Keyring != "" {
		keyringFile, err := ioutil.TempFile("", "terraform-provider-ceph")
		if err != nil {
			return nil, err
		}
		defer os.Remove(keyringFile.Name())
		if err = conn.SetConfigOption("keyring", keyringFile.Name()); err != nil {
			conn.Shutdown()
			return nil, err
		}
		if _, err = keyringFile.WriteString(config.Keyring); err != nil {
			return nil, err
		}
	}

	if err = conn.Connect(); err == nil {
		config.RadosConn = conn
	}
	return conn, err
}
