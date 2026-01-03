package singleton

import (
	"fmt"
	"sync"

	"github.com/hashicorp/consul/api"
)

var (
	consulClient *api.Client
	consulOnce   sync.Once
)

func GetConsulClient() (*api.Client, error) {

	var err error
	config := GetGlobalConfig()

	consulOnce.Do(func() {
		consulConfig := api.DefaultConfig()
		consulConfig.Address = config.ConsulAddress

		client, clientErr := api.NewClient(consulConfig)
		if clientErr != nil {
			err = fmt.Errorf("could not create consul client: %v", clientErr)
			return
		}

		// 3. Verify connection with a simple Agent check
		_, agentErr := client.Agent().Self()
		if agentErr != nil {
			err = fmt.Errorf("consul agent unreachable at %s: %v", config.ConsulAddress, agentErr)
			return
		}

		consulClient = client
	})

	return consulClient, err
}
