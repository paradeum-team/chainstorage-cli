package cmd

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	chainstoragesdk "github.com/paradeum-team/chainstorage-sdk"
	sdkcode "github.com/paradeum-team/chainstorage-sdk/code"
	"github.com/paradeum-team/chainstorage-sdk/consts"
	"github.com/paradeum-team/chainstorage-sdk/model"
	"github.com/paradeum-team/chainstorage-sdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ulule/deepcopier"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func init() {
	//carUploadCmd.Flags().StringP("Bucket", "b", "", "桶名称")
	//carUploadCmd.Flags().StringP("Object", "o", "", "上传对象路径")
	//
	//carImportCmd.Flags().StringP("Bucket", "b", "", "桶名称")
	//carImportCmd.Flags().StringP("Carfile", "c", "", "car文件标识")

	//objectRenameCmd.Flags().StringP("Bucket", "b", "", "桶名称")
	//objectRenameCmd.Flags().StringP("Object", "o", "", "对象名称")
	//objectRenameCmd.Flags().StringP("Cid", "c", "", "Cid")
	//objectRenameCmd.Flags().StringP("Rename", "r", "", "重命名")
	//objectRenameCmd.Flags().BoolP("Force", "f", false, "有冲突的时候强制覆盖")
	//
	//objectRemoveCmd.Flags().StringP("Bucket", "b", "", "桶名称")
	//objectRemoveCmd.Flags().StringP("Object", "o", "", "对象名称")
	//objectRemoveCmd.Flags().StringP("Cid", "c", "", "Cid")
	//objectRemoveCmd.Flags().BoolP("Force", "f", false, "有冲突的时候强制覆盖")
}

// region CAR Upload

//var carUploadCmd = &cobra.Command{
//	Use:     "put",
//	Short:   "put",
//	Long:    "upload object",
//	Example: "gcscmd put FILE[/DIR...] cs://BUCKET",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		carUploadRun(cmd, args)
//	},
//}

func carUploadRun(cmd *cobra.Command, args []string) {
	// 桶名称
	bucketName := GetBucketName(args)
	if err := checkBucketName(bucketName); err != nil {
		Error(cmd, args, err)
	}

	// 上传对象路径
	dataPath := GetDataPath(args)

	//// 上传 carfile
	//carFile, err := cmd.Flags().GetString("carFile")
	//if err != nil {
	//	Error(cmd, args, err)
	//}

	sdk, err := chainstoragesdk.New()
	if err != nil {
		Error(cmd, args, err)
	}

	// 确认桶数据有效性
	respBucket, err := sdk.Bucket.GetBucketByName(bucketName)
	if err != nil {
		Error(cmd, args, err)
	}

	code := int(respBucket.Code)
	if code != http.StatusOK {
		Error(cmd, args, errors.New(respBucket.Msg))
	}

	// 桶ID
	bucketId := respBucket.Data.Id

	// 对象上传
	response, err := UploadData(sdk, bucketId, dataPath, false)
	if err != nil {
		Error(cmd, args, err)
	}

	carUploadRunOutput(cmd, args, response)
}

func carUploadRunOutput(cmd *cobra.Command, args []string, resp model.ObjectCreateResponse) {
	respCode := int(resp.Code)

	if respCode != http.StatusOK {
		Error(cmd, args, errors.New(resp.Msg))
	}

	carUploadOutput := CarUploadOutput{
		RequestId: resp.RequestId,
		Code:      resp.Code,
		Msg:       resp.Msg,
		Status:    resp.Status,
	}

	err := deepcopier.Copy(&resp.Data).To(&carUploadOutput.Data)
	if err != nil {
		Error(cmd, args, err)
	}

	//对象上传
	//通过命令向固定桶内上传对象，包括文件、目录
	//
	//模版
	//
	//gcscmd put FILE[/DIR...] cs://BUCKET
	//BUCKET
	//
	//桶名称
	//
	//命令行例子
	//
	//上传文件
	//
	//当前目录
	//
	//gcscmd put ./aaa.mp4 cs://bbb
	//绝对路径
	//
	//gcscmd put /home/pz/aaa.mp4 cs://bbb
	//相对路径
	//
	//gcscmd put ../pz/aaa.mp4 cs://bbb
	//上传目录
	//
	//gcscmd put ./aaaa cs://bbb
	//上传 carfile
	//
	//gcscmd put ./aaa.car cs://bbb --carfile
	//响应
	//
	//过程
	//
	//################                                                                15%
	//Tarkov.mp4
	//完成
	//
	//CID:    QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
	//Name:Tarkov.mp4
	//报错
	//
	//Error: This file is a car file, add --carfile to confirm uploading car

	templateContent := `
CID: {{.Data.ObjectCid}}
Name: {{.Data.ObjectName}}
`

	t, err := template.New("carUploadTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, carUploadOutput)
	if err != nil {
		Error(cmd, args, err)
	}
}

type CarUploadOutput struct {
	RequestId string       `json:"requestId,omitempty"`
	Code      int32        `json:"code,omitempty"`
	Msg       string       `json:"msg,omitempty"`
	Status    string       `json:"status,omitempty"`
	Data      ObjectOutput `json:"objectOutput,omitempty"`
}

type CarUploadResponse struct {
	RequestId string      `json:"requestId,omitempty"`
	Code      int32       `json:"code,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	Status    string      `json:"status,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// 上传数据
func UploadData(sdk *chainstoragesdk.CssClient, bucketId int, dataPath string, isCarFile bool) (model.ObjectCreateResponse, error) {
	//response := model.CarResponse{}
	response := model.ObjectCreateResponse{}

	// 数据路径为空
	if len(dataPath) == 0 {
		return response, sdkcode.ErrCarUploadFileInvalidDataPath
	}

	// 数据路径无效
	fileInfo, err := os.Stat(dataPath)
	if os.IsNotExist(err) {
		return response, sdkcode.ErrCarUploadFileInvalidDataPath
	} else if err != nil {
		return response, err
	}

	fileDestination := dataPath
	if !isCarFile {
		// add constant
		//carVersion := 1
		fileDestination = sdk.Car.GenerateTempFileName(utils.CurrentDate()+"_", ".tmp")
		//fileDestination := GenerateTempFileName("", ".tmp")
		//fmt.Printf("UploadData carVersion:%d, fileDestination:%s, dataPath:%s\n", carVersion, fileDestination, dataPath)

		// 创建Car文件
		err = sdk.Car.CreateCarFile(dataPath, fileDestination)
		if err != nil {
			fmt.Printf("Error:%+v\n", err)
			return response, sdkcode.ErrCarUploadFileCreateCarFileFail
		}

		// todo: 清除CAR文件，添加utils?
		defer func(fileDestination string) {
			if !viper.GetBool("cmd.clean_tmp_data") {
				return
			}

			err := os.Remove(fileDestination)
			if err != nil {
				fmt.Printf("Error:%+v\n", err)
				//logger.Errorf("file.Delete %s err: %v", fileDestination, err)
			}
		}(fileDestination)
	}

	// 解析CAR文件，获取DAG信息，获取文件或目录的CID
	rootLink := model.RootLink{}
	err = sdk.Car.ParseCarFile(fileDestination, &rootLink)
	if err != nil {
		fmt.Printf("Error:%+v\n", err)
		return response, sdkcode.ErrCarUploadFileParseCarFileFail
	}

	rootCid := rootLink.RootCid.String()
	objectCid := rootLink.Cid.String()
	objectSize := int64(rootLink.Size)
	objectName := rootLink.Name

	if isCarFile {
		objectCid = rootCid
		objectSize = fileInfo.Size()

		filename := filepath.Base(dataPath)
		filename = strings.TrimSuffix(filename, ".car")
		objectName = filename
	}

	// 设置请求参数
	carFileUploadReq := model.CarFileUploadReq{}
	carFileUploadReq.BucketId = bucketId
	carFileUploadReq.ObjectCid = objectCid
	carFileUploadReq.ObjectSize = objectSize
	carFileUploadReq.ObjectName = objectName
	carFileUploadReq.FileDestination = fileDestination
	carFileUploadReq.CarFileCid = rootCid

	// 上传为目录的情况
	if fileInfo.IsDir() || isCarFile {
		carFileUploadReq.ObjectTypeCode = consts.ObjectTypeCodeDir
	}

	// 计算文件sha256
	sha256, err := utils.GetFileSha256ByPath(fileDestination)
	if err != nil {
		fmt.Printf("Error:%+v\n", err)
		return response, sdkcode.ErrCarUploadFileComputeCarFileHashFail
	}
	carFileUploadReq.RawSha256 = sha256

	// 使用Root CID秒传检查
	objectExistResponse, err := sdk.Object.IsExistObjectByCid(objectCid)
	if err != nil {
		fmt.Printf("Error:%+v\n", err)
		return response, sdkcode.ErrCarUploadFileReferenceObjcetFail
	}

	// CID存在，执行秒传操作
	objectExistCheck := objectExistResponse.Data
	if objectExistCheck.IsExist {
		response, err := sdk.Car.ReferenceObject(&carFileUploadReq)
		if err != nil {
			fmt.Printf("Error:%+v\n", err)
			return response, sdkcode.ErrCarUploadFileReferenceObjcetFail
		}

		return response, err
	}

	// CAR文件大小，超过分片阈值
	carFileSize := fileInfo.Size()
	carFileShardingThreshold := sdk.Config.CarFileShardingThreshold

	// 生成CAR分片文件上传
	if carFileSize > int64(carFileShardingThreshold) {
		response, err = UploadBigCarFile(sdk, &carFileUploadReq)
		if err != nil {
			return response, sdkcode.ErrCarUploadFileFail
		}

		return response, nil
	}

	// 普通上传
	file, err := os.Open(fileDestination)
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return response, sdkcode.ErrCarUploadFileFail
	}
	size := fi.Size()

	bar := pb.Start64(size).SetWriter(os.Stdout).Set(pb.Bytes, true)
	bar.SetRefreshRate(100 * time.Millisecond)
	defer bar.Finish()

	extReader := bar.NewProxyReader(file)

	response, err = sdk.Car.UploadCarFileExt(&carFileUploadReq, extReader)
	if err != nil {
		return response, sdkcode.ErrCarUploadFileFail
	}

	return response, err
}

// 上传大CAR文件
func UploadBigCarFile(sdk *chainstoragesdk.CssClient, req *model.CarFileUploadReq) (model.ObjectCreateResponse, error) {
	response := model.ObjectCreateResponse{}

	// 生成CAR分片文件
	shardingCarFileUploadReqs := []model.CarFileUploadReq{}
	err := sdk.Car.GenerateShardingCarFiles(req, &shardingCarFileUploadReqs)
	if err != nil {
		return response, err
	}
	//todo:delete CAR分片文件
	defer func(shardingCarFileUploadReqs []model.CarFileUploadReq) {
		if !viper.GetBool("cmd.clean_tmp_data") {
			return
		}

		for i := range shardingCarFileUploadReqs {
			fileDestination := shardingCarFileUploadReqs[i].FileDestination
			err := os.Remove(fileDestination)
			if err != nil {
				fmt.Printf("Error:%+v\n", err)
				//logger.Errorf("file.Delete %s err: %v", fileDestination, err)
			}
		}
	}(shardingCarFileUploadReqs)

	totalSize := int64(0)
	for i, _ := range shardingCarFileUploadReqs {
		totalSize += shardingCarFileUploadReqs[i].ObjectSize
	}

	bar := pb.Start64(totalSize).SetWriter(os.Stdout).Set(pb.Bytes, true)
	bar.SetRefreshRate(100 * time.Millisecond)
	defer bar.Finish()

	// 上传CAR文件分片
	//uploadingReqs := []model.CarFileUploadReq{}
	//deepcopier.Copy(&shardingCarFileUploadReqs).To(&uploadingReqs)
	// todo: 添加配置，重试3次，每次间隔3秒
	maxRetries := 3
	retryDelay := time.Duration(3) * time.Second

	uploadRespList := []model.ShardingCarFileUploadResponse{}
	for i, _ := range shardingCarFileUploadReqs {
		for j := 0; j < maxRetries; j++ {
			uploadingReq := model.CarFileUploadReq{}
			deepcopier.Copy(&shardingCarFileUploadReqs[i]).To(&uploadingReq)

			file, err := os.Open(uploadingReq.FileDestination)
			defer file.Close()
			//fi, err := file.Stat()
			//size := fi.Size()
			extReader := bar.NewProxyReader(file)

			uploadResp, err := sdk.Car.UploadShardingCarFileExt(&uploadingReq, extReader)
			if err == nil && uploadResp.Code == http.StatusOK {
				uploadRespList = append(uploadRespList, uploadResp)
				break
			}
			// todo: log err?

			if j == maxRetries-1 {
				// 尝试maxRetries次失败
				if err != nil {
					return response, err
				} else if uploadResp.Code != http.StatusOK {
					return response, errors.New(response.Msg)
				}
			}

			time.Sleep(retryDelay)
		}
	}

	// 确认分片上传成功
	response, err = sdk.Car.ConfirmShardingCarFiles(req)
	if err != nil {
		return response, err
	}

	return response, nil
}

// endregion CAR Upload

// region CAR Import

//var carImportCmd = &cobra.Command{
//	Use:     "import",
//	Short:   "import",
//	Long:    "import car file",
//	Example: "gcscmd import  ./aaa.car cs://BUCKET",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		//cmd.Help()
//		//fmt.Printf("%s %s\n", cmd.Name(), strconv.Itoa(offset))
//		carImportRun(cmd, args)
//	},
//}

func carImportRun(cmd *cobra.Command, args []string) {
	// 桶名称
	bucketName := GetBucketName(args)
	if err := checkBucketName(bucketName); err != nil {
		Error(cmd, args, err)
	}

	// 上传对象路径
	dataPath := GetDataPath(args)

	// CAR文件类型检查
	if !strings.HasSuffix(strings.ToLower(dataPath), ".car") {
		err := sdkcode.ErrCarUploadFileInvalidDataPath
		Error(cmd, args, err)
	}

	//// CAR文件标识
	//carFile, err := cmd.Flags().GetString("carFile")
	//if err != nil {
	//	Error(cmd, args, err)
	//}

	sdk, err := chainstoragesdk.New()
	if err != nil {
		Error(cmd, args, err)
	}

	// 确认桶数据有效性
	respBucket, err := sdk.Bucket.GetBucketByName(bucketName)
	if err != nil {
		Error(cmd, args, err)
	}

	code := int(respBucket.Code)
	if code != http.StatusOK {
		Error(cmd, args, errors.New(respBucket.Msg))
	}

	// 桶ID
	bucketId := respBucket.Data.Id

	// 对象上传
	response, err := UploadData(sdk, bucketId, dataPath, true)
	if err != nil {
		Error(cmd, args, err)
	}

	carImportRunOutput(cmd, args, response)
}

func carImportRunOutput(cmd *cobra.Command, args []string, resp model.ObjectCreateResponse) {
	code := resp.Code
	if code != http.StatusOK {
		Error(cmd, args, errors.New(resp.Msg))
	}

	carImportOutput := CarImportOutput{
		RequestId: resp.RequestId,
		Code:      resp.Code,
		Msg:       resp.Msg,
		Status:    resp.Status,
	}

	err := deepcopier.Copy(&resp.Data).To(&carImportOutput.Data)
	if err != nil {
		Error(cmd, args, err)
	}

	//	导入 car 文件
	//	通过命令向固定桶内导入 car 文件对象
	//
	//	模版
	//
	//	gcscmd import  ./aaa.car cs://BUCKET
	//	BUCKET
	//
	//	桶名称
	//
	//	carfile
	//
	//	car文件标识
	//
	//	命令行例子
	//
	//	当前目录
	//
	//	gcscmd import ./aaa.car cs://bbb
	//	绝对路径
	//
	//	gcscmd import /home/pz/aaa.car cs://bbb
	//	相对路径
	//
	//	gcscmd import ../pz/aaa.car cs://bbb
	//	响应
	//
	//	过程
	//
	//	################                                                                15%
	//		QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo        Tarkov.mp4
	//	完成
	//
	//CID:    QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
	//Name:Tarkov.mp4
	//	报错
	//
	//Error: This is not a carfile

	templateContent := `
CID: {{.ObjectCid}}
Name: {{.ObjectName}}
`

	t, err := template.New("carImportTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, carImportOutput)
	if err != nil {
		Error(cmd, args, err)
	}
}

type CarImportOutput struct {
	RequestId string       `json:"requestId,omitempty"`
	Code      int32        `json:"code,omitempty"`
	Msg       string       `json:"msg,omitempty"`
	Status    string       `json:"status,omitempty"`
	Data      ObjectOutput `json:"objectOutput,omitempty"`
}

//type ObjectOutput struct {
//	Id             int       `json:"id" comment:"对象ID"`
//	BucketId       int       `json:"bucketId" comment:"桶主键"`
//	ObjectName     string    `json:"objectName" comment:"对象名称（255字限制）"`
//	ObjectTypeCode int       `json:"objectTypeCode" comment:"对象类型编码"`
//	ObjectSize     int64     `json:"objectSize" comment:"对象大小（字节）"`
//	IsMarked       int       `json:"isMarked" comment:"星标（1-已标记，0-未标记）"`
//	ObjectCid      string    `json:"objectCid" comment:"对象CID"`
//	CreatedAt      time.Time `json:"createdAt" comment:"创建时间"`
//	UpdatedAt      time.Time `json:"updatedAt" comment:"最后更新时间"`
//	CreatedDate    string    `json:"createdDate" comment:"创建日期"`
//}

// endregion CAR Import

//func makeBar(req *model.CarFileUploadReq) *pb.ProgressBar {
//	objectSize := int(req.ObjectSize)
//	bar := pb.New(objectSize).
//
//	bar := pb.New(int(sourceSize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
//	bar.ShowSpeed = true
//	bar.
//	// show percents (by default already true)
//	bar.ShowPercent = true
//
//	// show bar (by default already true)
//	bar.ShowBar = true
//
//	bar.ShowCounters = true
//
//	bar.ShowTimeLeft = true
//}
