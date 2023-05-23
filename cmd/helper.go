package cmd

import (
	"fmt"
	sdkcode "github.com/paradeum-team/chainstorage-sdk/sdk/code"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strings"
)

func Error(cmd *cobra.Command, args []string, err error) {
	fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}

func GetBucketName(args []string) string {
	bucketName := ""
	if len(args) == 0 {
		return bucketName
	}

	bucketPrefix := viper.GetString("cmd.bucket_prefix")

	for i := range args {
		arg := args[i]
		if strings.HasPrefix(arg, bucketPrefix) {
			bucketName = strings.TrimPrefix(arg, bucketPrefix)
			break
		}
	}

	return bucketName
}

func GetDataPath(args []string) string {
	dataPath := ""
	if len(args) == 0 {
		return dataPath
	}

	bucketPrefix := viper.GetString("cmd.bucket_prefix")

	for i := range args {
		arg := args[i]
		if strings.HasPrefix(arg, bucketPrefix) {
			continue
		}

		if _, err := os.Stat(arg); !os.IsNotExist(err) {
			return arg
		}
	}

	return dataPath
}

// 检查桶名称
func checkBucketName(bucketName string) error {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return sdkcode.ErrInvalidBucketName
	}

	// 桶名称异常，名称范围必须在 3-63 个字符之间并且只能包含小写字符、数字和破折号，请重新尝试
	isMatch := regexp.MustCompile(`^[a-z0-9-]*$`).MatchString(bucketName)
	if !isMatch {
		return sdkcode.ErrInvalidBucketName
	}

	return nil
}

// 检查对象名称
func checkObjectName(objectName string) error {
	if len(objectName) == 0 || len(objectName) > 255 {
		return sdkcode.ErrInvalidObjectName
	}

	isMatch := regexp.MustCompile("[<>:\"/\\|?*\u0000-\u001F]").MatchString(objectName)
	if isMatch {
		return sdkcode.ErrInvalidObjectName
	}

	isMatch = regexp.MustCompile(`^(con|prn|aux|nul|com\d|lpt\d)$`).MatchString(objectName)
	if isMatch {
		return sdkcode.ErrInvalidObjectName
	}

	if objectName == "." || objectName == ".." {
		return sdkcode.ErrInvalidObjectName
	}

	return nil
}
