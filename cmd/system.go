package cmd

import (
	chainstoragesdk "github.com/paradeum-team/chainstorage-sdk"
	"github.com/paradeum-team/chainstorage-sdk/model"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

func ipfsVersionRun(cmd *cobra.Command, args []string) {
	sdk, err := chainstoragesdk.New(&appConfig)
	if err != nil {
		Error(cmd, args, err)
	}

	response, err := sdk.GetIpfsVersion()
	if err != nil {
		Error(cmd, args, err)
	}

	ipfsVersionRunOutput(cmd, args, response)
}

func ipfsVersionRunOutput(cmd *cobra.Command, args []string, resp model.VersionResponse) {
	//code := int(resp.Code)
	//if code != http.StatusOK {
	//	Error(cmd, args, errors.New(resp.Msg))
	//}

	respData := resp.Data

	templateContent := `
IPFS Version: {{.Version}}
`

	t, err := template.New("ipfsVersionTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, respData)
	if err != nil {
		Error(cmd, args, err)
	}
}

func versionRun(cmd *cobra.Command, args []string) {
	versionInfo := GetVersionInfo()
	versionRunOutput(cmd, args, versionInfo)
}

func versionRunOutput(cmd *cobra.Command, args []string, versionInfo *VersionInfo) {
	templateContent := `
Client Version: version.Info{Version:"v{{.Version}}"}
Server Version: version.Info{Version:"v{{.ApiVersion}}"}
`

	t, err := template.New("versionTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, versionInfo)
	if err != nil {
		Error(cmd, args, err)
	}
}
