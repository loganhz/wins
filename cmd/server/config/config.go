package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/rancher/wins/cmd/cmds"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/rancher/system-agent/pkg/config"
	"github.com/rancher/wins/pkg/csiproxy"
	"github.com/rancher/wins/pkg/defaults"
	"github.com/rancher/wins/pkg/tls"
)

func DefaultConfig() *Config {
	return &Config{
		Listen: defaults.NamedPipeName,
		Proxy:  defaults.ProxyPipeName,
		WhiteList: WhiteListConfig{
			ProcessPaths: []string{},
			ProxyPorts:   []int{},
		},
		Upgrade: UpgradeConfig{
			Mode:         "watching",
			WatchingPath: defaults.UpgradeWatchingPath,
		},
		TLSConfig: &tls.Config{
			CertFilePath: defaults.CertPath,
		},
	}
}

type Config struct {
	Debug       bool                `yaml:"debug" json:"debug"`
	Listen      string              `yaml:"listen" json:"listen"`
	Proxy       string              `yaml:"proxy" json:"proxy"`
	WhiteList   WhiteListConfig     `yaml:"white_list" json:"white_list"`
	Upgrade     UpgradeConfig       `yaml:"upgrade" json:"upgrade"`
	SystemAgent *config.AgentConfig `yaml:"systemagent" json:"systemagent"`
	CSIProxy    *csiproxy.Config    `yaml:"csi-proxy" json:"csi-proxy"`
	TLSConfig   *tls.Config         `yaml:"tls-config" json:"tls-config"`
}

func (c *Config) ValidateTLSConfig() error {
	if b, err := ioutil.ReadFile(c.TLSConfig.CertFilePath); b == nil || err != nil {
		return errors.Wrapf(err, "failed to read certificate from %s", c.TLSConfig.CertFilePath)
	}

	if c.TLSConfig.CertFilePath != defaults.CertPath {
		// load non-default certificate file
		_ = csiproxy.Config{
			Config: tls.Config{
				CertFilePath: c.TLSConfig.CertFilePath,
			},
		}
	}

	if *c.TLSConfig.Insecure {
		// set insecure flag for all subsequent CSI Proxy functions
		_ = csiproxy.Config{
			Config: tls.Config{
				Insecure: cmds.BoolAddr(true),
			},
		}
	} else {
		_ = csiproxy.Config{
			Config: tls.Config{
				Insecure: cmds.BoolAddr(false),
			},
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Listen) == "" {
		return errors.New("listen could not be blank")
	}

	// validate white list field
	if err := c.WhiteList.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate white list field")
	}

	// validate upgrade field
	if err := c.Upgrade.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate upgrade field")
	}

	// validate

	return nil
}

type WhiteListConfig struct {
	ProcessPaths []string `yaml:"process_paths" json:"processPaths"`
	ProxyPorts   []int    `yaml:"proxy_ports" json:"proxyPorts"`
}

func (c *WhiteListConfig) Validate() error {
	// process path
	for _, processPath := range c.ProcessPaths {
		if strings.TrimSpace(processPath) == "" {
			return errors.New("could not accept blank path as process white list")
		}
	}
	for _, proxyPort := range c.ProxyPorts {
		if proxyPort < 0 || proxyPort > 0xFFFF {
			return errors.New("could not accept invalid port number in proxy ports")
		}
	}
	return nil
}

type UpgradeConfig struct {
	Mode         string `yaml:"mode" json:"mode"`
	WatchingPath string `yaml:"watching_path" json:"watchingPath"`
}

func (c *UpgradeConfig) Validate() error {
	switch strings.TrimSpace(c.Mode) {
	case "watching":
		if strings.TrimSpace(c.WatchingPath) == "" {
			return errors.New("could not accept blank path as watching path")
		}
	case "none", "":
	default:
		return errors.Errorf("could not accept %q as upgrade mode", c.Mode)
	}

	return nil
}

func (c *UpgradeConfig) IsWatchingMode() bool {
	return c.Mode == "watching"
}

func LoadConfig(path string, v *Config) error {
	if v == nil {
		return errors.New("config could not be nil")
	}

	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "could not load config")
		}
		return nil
	} else if stat.IsDir() {
		return errors.New("could not load config from directory")
	}

	if err := DecodeConfig(path, v); err != nil {
		return errors.Wrap(err, "could not decode config")
	}

	return v.Validate()
}

func DecodeConfig(path string, v *Config) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, v)
}
