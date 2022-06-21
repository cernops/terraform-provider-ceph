package ceph

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ceph/go-ceph/rados"
)

type Config struct {
	ConfigPath string
	Username   string
	Cluster    string
	Keyring    string
	MonHost    string
}

func (config Config) GetCephConnection() (*rados.Conn, error) {
	var conn *rados.Conn
	var err error

	if config.Username != "" && config.Cluster != "" {
		conn, err = rados.NewConnWithClusterAndUser(config.Cluster, config.Username)
	} else if config.Username != "" {
		conn, err = rados.NewConnWithUser(config.Username)
	} else {
		conn, err = rados.NewConn()
	}
	if err != nil {
		conn.Shutdown()
		return nil, err
	}

	configPath := config.ConfigPath
	if config.Keyring != "" {
		if config.MonHost == "" {
			conn.Shutdown()
			return nil, fmt.Errorf("Error creating Ceph connection: keyring specified while mon_host is not")
		}

		keyringFile, err := ioutil.TempFile("", "terraform-provider-ceph")
		if err != nil {
			return nil, err
		}
		defer os.Remove(keyringFile.Name())
		_, err = keyringFile.WriteString(config.Keyring)
		if err != nil {
			return nil, err
		}

		cephConfigFile, err := ioutil.TempFile("", "terraform-provider-ceph")
		if err != nil {
			return nil, err
		}
		defer os.Remove(cephConfigFile.Name())
		cephConfig := fmt.Sprintf(`
[global]
mon host = %s
keyring = %s
`, config.MonHost, keyringFile.Name())
		_, err = cephConfigFile.WriteString(cephConfig)
		if err != nil {
			return nil, err
		}
	}

	if configPath != "" {
		err = conn.ReadConfigFile(configPath)
	} else {
		err = conn.ReadDefaultConfigFile()
	}
	if err != nil {
		conn.Shutdown()
		return nil, err
	}

	err = conn.Connect()
	if err != nil {
		conn.Shutdown()
	}
	return conn, err
}
