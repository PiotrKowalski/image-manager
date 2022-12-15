package remote

import (
	"fmt"
	"github.com/PiotrKowalski/image-manager/pkg/config"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"time"
)

const (
	ProviderName        = "consul"
	KV_STORE_KEY string = "KV_STORE_KEY"
)

type consul struct {
	isWatching bool
	watcher    chan interface{}
}

func NewConsulRemoteConfigProvider() SecretRemoteConfigProvider {
	return consul{
		isWatching: false,
		watcher:    make(chan interface{}),
	}
}

func (c consul) LoadStoreConfig() error {
	secret, err := config.ReadEnvString(KV_STORE_KEY)
	if err != nil {
		return err
	}
	host, err := config.ReadEnvString(KV_STORE_HOST)
	if err != nil {
		return err
	}
	port, err := config.ReadEnvString(KV_STORE_PORT)
	if err != nil {
		return err
	}

	log.Println(host, port, secret)

	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Reading config from remote %v", addr)
	err = viper.AddRemoteProvider(ProviderName, addr, secret)
	if err != nil {
		log.Println(err)
		return err
	}

	viper.SetConfigType("json") // Need to explicitly set this to json
	err = viper.ReadRemoteConfig()
	if err != nil {
		log.Printf("unable to read remote config: %s", err)
	}

	return nil
}

func (c consul) StartRemoteWatch() error {
	if c.isWatching {
		return nil
	}

	go func() {
		for {
			select {
			case <-c.watcher:
				c.isWatching = false
				return
			case <-time.After(5 * time.Second):
				err := viper.WatchRemoteConfig()
				if err != nil {
					log.Printf("unable to read remote config: %s", err)
					continue
				}

			}

		}

	}()

	return nil
}

func (c consul) StopRemoteWatch() error {
	if !c.isWatching {
		return nil
	}

	c.watcher <- 0

	return nil
}
