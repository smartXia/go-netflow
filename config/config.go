package config

import "github.com/rfyiamcool/go-netflow/utils"

type (
	LogConfig struct {
		Level    string `json:"level" yaml:"level"`
		Filepath string `json:"filePath" yaml:"filePath"`
	}

	AgentConfig struct {
		// If root dir is not explicitly provided, we use the dir
		// that contains this app config.
		RootDir        string `json:"rootDir" yaml:"rootDir"`
		ServerEndpoint string `json:"serverEndpoint" yaml:"serverEndpoint"`
		DeviceId       string `json:"deviceId" yaml:"deviceId"`
		BizType        string `json:"bizType" yaml:"bizType"`
		AppKey         string `json:"appKey" yaml:"appKey"`
		AppSecret      string `json:"appSecret" yaml:"appSecret"`
	}

	PppoeConfig struct {
		PingAddr string `json:"pingAddr" yaml:"pingAddr"`
		FrpsAddr string `json:"frpsAddr" yaml:"frpsAddr"`
		FrpcRoot string `json:"frpcroot" yaml:"frpcroot"`
	}

	MonitorConfig struct {
		CollectdRddPath string `json:"collectdRddPath" yaml:"collectdRddPath"`
		ReportInterval  int64  `json:"reportInterval" yaml:"reportInterval"`
		ReportBucket    int    `json:"reportBucket" yaml:"reportBucket"`
	}

	// super-agent app config
	Config struct {
		Log                  LogConfig     `json:"log" yaml:"log"`
		Agent                AgentConfig   `json:"agent" yaml:"agent"`
		Pppoe                PppoeConfig   `json:"pppoe" yaml:"pppoe"`
		MonitorConfig        MonitorConfig `json:"monitor" yaml:"monitor"`
		MockedServerConfPath string        `json:"mockedServerConfPath" yaml:"mockedServerConfPath"`
		DeviceIdPath         string        `json:"deviceIdPath" yaml:"deviceIdPath"`
		Nethogs              string        `json:"nethogs" yaml:"nethogs"`
	}
)

func GetConfig(filePath string) Config {
	var config Config
	if err := utils.UnmarshalConfigFromFile(filePath, &config); err != nil {
		panic(err)
	}
	if config.Agent.ServerEndpoint == "" {
		panic("server endpoint path must be provided")
	}

	if config.DeviceIdPath == "" {
		panic("deviceId filePath must be provided")
	}

	return config
}
