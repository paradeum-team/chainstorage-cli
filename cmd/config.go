package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

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
	CarVersion               int    `toml:"carVersion"`
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

// region Config show

//var bucketListCmd = &cobra.Command{
//	Use:     "ls",
//	Short:   "ls",
//	Long:    "List links from object or bucket",
//	Example: "gcscmd ls [--Offset=<Offset>]",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		//cmd.Help()
//		//fmt.Printf("%s %s\n", cmd.Name(), strconv.Itoa(offset))
//		bucketListRun(cmd, args)
//	},
//}

func configShowRun(cmd *cobra.Command, args []string) {
	configFileUsed := viper.ConfigFileUsed()
	fmt.Fprintln(os.Stderr, "Using config file:", configFileUsed)
	err := printFileContent(configFileUsed)
	if err != nil {
		Error(cmd, args, err)
	}

	configShowRunOutput(cmd, args)
}

func configShowRunOutput(cmd *cobra.Command, args []string) {
	//	templateContent := `
	//total {{.Count}}
	//{{- if eq (len .List) 0}}
	//Status: {{.Code}}
	//{{- else}}
	//{{- range .List}}
	//{{.ObjectAmount}} {{.StorageNetwork}} {{.BucketPrinciple}} {{.UsedSpace}} {{.CreatedDate}} {{.BucketName}}
	//{{- end}}
	//{{- end}}
	//`
	//
	//	t, err := template.New("configShowTemplate").Parse(templateContent)
	//	if err != nil {
	//		Error(cmd, args, err)
	//	}
	//
	//	err = t.Execute(os.Stdout, bucketListOutput)
	//	if err != nil {
	//		Error(cmd, args, err)
	//	}
}

//type BucketListOutput struct {
//	RequestId string         `json:"requestId,omitempty"`
//	Code      int32          `json:"code,omitempty"`
//	Msg       string         `json:"msg,omitempty"`
//	Status    string         `json:"status,omitempty"`
//	Count     int            `json:"count,omitempty"`
//	PageIndex int            `json:"pageIndex,omitempty"`
//	PageSize  int            `json:"pageSize,omitempty"`
//	List      []BucketOutput `json:"list,omitempty"`
//}
//
//type BucketOutput struct {
//	Id                  int       `json:"id" comment:"桶ID"`
//	BucketName          string    `json:"bucketName" comment:"桶名称（3-63字长度限制）"`
//	StorageNetworkCode  int       `json:"storageNetworkCode" comment:"存储网络编码（10001-IPFS）"`
//	BucketPrincipleCode int       `json:"bucketPrincipleCode" comment:"桶策略编码（10001-公开，10000-私有）"`
//	UsedSpace           int64     `json:"usedSpace" comment:"已使用空间（字节）"`
//	ObjectAmount        int       `json:"objectAmount" comment:"对象数量"`
//	CreatedAt           time.Time `json:"createdAt" comment:"创建时间"`
//	StorageNetwork      string    `json:"storageNetwork" comment:"存储网络（10001-IPFS）"`
//	BucketPrinciple     string    `json:"bucketPrinciple" comment:"桶策略（10001-公开，10000-私有）"`
//	CreatedDate         string    `json:"createdDate" comment:"创建日期"`
//}

// endregion Config show
