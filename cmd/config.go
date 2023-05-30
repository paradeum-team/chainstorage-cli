package cmd

type CliConfig struct {
	IpfsGateway      string `toml:"ipfsGateway"`
	GgcscmdPath      string `toml:"ggcscmdPath"`
	UseHTTPSProtocol string `toml:"useHttpsProtocol"`
	BucketPrefix     string `toml:"bucketPrefix"`
	ListOffset       int    `toml:"listOffset"`
	CleanTmpData     bool   `toml:"cleanTmpData"`
	MaxRetries       int    `toml:"maxRetries"`
	RetryDelay       int    `toml:"retryDelay"`
}

type SdkConfig struct {
	DefaultRegion            string `toml:"defaultRegion"`
	TimeZone                 string `toml:"timeZone"`
	ChainStorageApiEndpoint  string `toml:"chainStorageApiEndpoint"`
	CarFileWorkPath          string `toml:"carFileWorkPath"`
	CarFileShardingThreshold int    `toml:"carFileShardingThreshold"`
	ChainStorageApiToken     string `toml:"chainStorageApiToken"`
	HTTPRequestUserAgent     string `toml:"httpRequestUserAgent"`
	HTTPRequestOvertime      int    `toml:"httpRequestOvertime"`
}

type LoggerConfig struct {
	LogPath      string `toml:"logPath"`
	Mode         string `toml:"mode"`
	Level        string `toml:"level"`
	IsOutPutFile bool   `toml:"isOutPutFile"`
	MaxAgeDay    int    `toml:"maxAgeDay"`
	RotationTime int    `toml:"rotationTime"`
	UseJSON      bool   `toml:"useJson"`
	LoggerFile   string `toml:"loggerFile"`
}

type CscConfig struct {
	Cli    CliConfig    `toml:"cli"`
	Sdk    SdkConfig    `toml:"sdk"`
	Logger LoggerConfig `toml:"logger"`
}
