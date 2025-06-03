package config

import (
	"go-zrbc/pkg/xlog"

	jsoniter "github.com/json-iterator/go"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
)

// InitConfigWithJson 初始化Json配置
func InitConfigWithJson(dataID, group string, v interface{}) {
	configClient := initNacosClient(dataID, group)
	if configClient == nil {
		xlog.Errorf("configClient nil")
		panic("configClient nil")
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group})
	if err != nil {
		xlog.Errorf("GetConfig error: %v", err.Error())
		panic("GetConfig error")
	}
	xlog.Infof("config content: %v", content)

	err = jsoniter.UnmarshalFromString(content, v)
	if err != nil {
		xlog.Errorf("UnmarshalFromString error: %v", err.Error())
		panic("UnmarshalFromString error")
	}

	// watch
	configClient.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			xlog.Infof("config changed group: %v, dataID: %v, data: %v", group, dataId, data)
			err := jsoniter.UnmarshalFromString(data, v)
			if err != nil {
				xlog.Errorf("UnmarshalFromString error: %v", err.Error())
			}
		},
	})
}

// InitConfigWithYaml 初始化Yaml配置
func InitConfigWithYaml(dataID, group string, v interface{}) {
	configClient := initNacosClient(dataID, group)
	if configClient == nil {
		xlog.Errorf("configClient nil")
		panic("configClient nil")
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group})
	if err != nil {
		xlog.Errorf("GetConfig error: %v", err.Error())
		panic("GetConfig error")
	}
	xlog.Infof("config content: %v", content)

	err = yaml.Unmarshal([]byte(content), v)
	if err != nil {
		xlog.Errorf("UnmarshalYaml error: %v", err.Error())
		panic("UnmarshalYaml error")
	}

	// watch
	configClient.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			xlog.Infof("config changed group: %v, dataID: %v, data: %v", group, dataId, data)
			err := yaml.Unmarshal([]byte(data), v)
			if err != nil {
				xlog.Errorf("UnmarshalYaml error: %v", err.Error())
			}
		},
	})
}

func initNacosClient(dataID, group string) config_client.IConfigClient {
	xlog.Infof("init config")
	// create ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "nacos-headless.config.svc.cluster.local",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// create ClientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// create config client
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		xlog.Errorf("initConfig error: %v", err.Error())
		panic("initConfig error")
	}

	return configClient
}
