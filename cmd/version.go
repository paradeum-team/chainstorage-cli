package cmd

import (
	"fmt"
	chainstoragesdk "github.com/paradeum-team/chainstorage-sdk"
	"os"
	"runtime"
)

// CurrentVersionNumber is the current application's version literal
const CurrentVersionNumber = "0.0.4"

func GetApiVersion() string {
	sdk, err := chainstoragesdk.New(&appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "get api version fail, error:%+v\n", err)
		os.Exit(1)
	}

	response, err := sdk.GetApiVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "get api version fail, error:%+v\n", err)
		os.Exit(1)
	}

	return response.Data.Version
}

type VersionInfo struct {
	Version    string
	ApiVersion string
	System     string
	Golang     string
}

func GetVersionInfo() *VersionInfo {
	return &VersionInfo{
		ApiVersion: GetApiVersion(),
		Version:    CurrentVersionNumber,
		System:     runtime.GOARCH + "/" + runtime.GOOS,
		Golang:     runtime.Version(),
	}
}
