package cmd

import (
	"bufio"
	"fmt"
	sdkcode "github.com/paradeum-team/chainstorage-sdk/sdk/code"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Error(cmd *cobra.Command, args []string, err error) {
	log.Errorf("execute %s args:%v error:%v\n", cmd.Name(), args, err)
	fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}

func GetBucketName(args []string) string {
	bucketName := ""
	if len(args) == 0 {
		return bucketName
	}

	//bucketPrefix := viper.GetString("cli.bucketPrefix")
	bucketPrefix := cliConfig.BucketPrefix

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

	//bucketPrefix := viper.GetString("cli.bucketPrefix")
	bucketPrefix := cliConfig.BucketPrefix

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

func GetTimestampString() string {
	timestampString := time.Now().Format("2006-01-02 15:04:05.000000000") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	return timestampString
}

func isFolderNotEmpty(path string) (bool, error) {
	// Check if the path is a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !fileInfo.IsDir() {
		return false, nil
	}

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// Read the directory entries
	_, err = dir.Readdirnames(1)
	if err == nil {
		// Directory is not empty
		return true, nil
	} else if err == io.EOF {
		// Directory is empty
		return false, nil
	} else {
		// An error occurred while reading the directory
		return false, err
	}
}

func getFolderSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			size += fileInfo.Size()
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

//func main() {
//	path := "/path/to/folder"
//
//	size, err := folderSize(path)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("Folder size: %d bytes\n", size)
//}

func printFileContent(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(os.Stdout, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
