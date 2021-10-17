package rustack_terraform

import (
	"fmt"

	"github.com/pilat/rustack-go/rustack"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token            string
	APIEndpoint      string
	TerraformVersion string
}

type CombinedConfig struct {
	manager *rustack.Manager
}

func (c *CombinedConfig) rustackManager() *rustack.Manager { return c.manager }

func (c *Config) Client() (*CombinedConfig, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	manager := rustack.NewManager(c.Token)
	manager.Logger = logger
	manager.BaseURL = c.APIEndpoint
	manager.UserAgent = fmt.Sprintf("Terraform/%s", c.TerraformVersion)

	return &CombinedConfig{
		manager: manager,
	}, nil
}
