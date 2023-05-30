package cmd

import (
	chainstoragesdk "github.com/paradeum-team/chainstorage-sdk/sdk"
	"github.com/paradeum-team/chainstorage-sdk/sdk/model"
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
