package shanda_config

import (
	"encoding/json"
	"log"
	"os"

	reclog "github.com/alibaba/pairec/log"
	"go.uber.org/zap"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

type Config struct {
	LoggerCfg *zap.Config `json:"logger_cfg" yaml:"logger_cfg"`
}

func getApolloInfo() *config.AppConfig {
	var ip string
	var secret string
	namespace := "a012_infer.json"

	env := os.Getenv(ENV_APOLLO_URL)
	if env == "" {
		ip = "http://dev-tanka-sg-apollo-config.aws.tankaapps.com:80"
		secret = "d4d279eb0f5040d7b9fb3ab79c6d5d1f"
		namespace = "a012_infer_local.json"
	} else {
		// this golang app is using .json config instead of .properties which is used by python
		// the ip should use http://dev-tanka-sg-apollo-config.aws.tankaapps.com:80 instead of
		// http://dev-tanka-sg-apollo-admin.aws.tankaapps.com:80, which is used by python.
		ip = os.Getenv(ENV_APOLLO_URL)
		secret = os.Getenv(ENV_APOLLO_SECRET)
		namespace = os.Getenv(ENV_APOLLO_NAMESPACE) + ".json"
	}

	c := &config.AppConfig{
		AppID:          "solab-ai",
		Cluster:        "default",
		IP:             ip,
		NamespaceName:  namespace,
		IsBackupConfig: true,
		Secret:         secret,
	}
	return c
}

func LoadConfig() (*Config, error) {
	c := getApolloInfo()

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		log.Fatalf("LoadFromApollo StartWithConfig err: %s", err.Error())
	}

	cache := client.GetConfigCache(c.NamespaceName)
	content, err := cache.Get("content")
	if err != nil {
		log.Fatalf("LoadFromApollo GetConfigCache err: %s", err.Error())
	}

	contentBytes := []byte(content.(string))
	var cfg Config
	err = json.Unmarshal(contentBytes, &cfg)
	if err != nil {
		log.Fatalf("LoadFromApollo Unmarshal err: %s", err.Error())
	}

	reclog.SetConfig(cfg.LoggerCfg)

	return &cfg, nil
}
