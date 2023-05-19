package cmd

import (
	"chainstoragesdk"
	sdkcode "chainstoragesdk/code"
	"chainstoragesdk/model"
	"context"
	"fmt"
	"github.com/Code-Hex/pget"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ulule/deepcopier"
	"net/http"
	"os"
	"text/template"
	"time"
)

func init() {
	//objectListCmd.Flags().StringP("Bucket", "b", "", "桶名称")
	//objectListCmd.Flags().StringP("Object", "r", "", "对象名称")
	//objectListCmd.Flags().StringP("Cid", "c", "", "Cid")
	//objectListCmd.Flags().IntP("Offset", "o", 10, "查询偏移量")
	//
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
	//
	//objectDownloadCmd.Flags().StringP("Bucket", "b", "", "桶名称")
	//objectDownloadCmd.Flags().StringP("Object", "o", "", "对象名称")
	//objectDownloadCmd.Flags().StringP("Cid", "c", "", "Cid")
	//objectDownloadCmd.Flags().BoolP("Target", "t", false, "输出路径")
}

// region Object List

//var objectListCmd = &cobra.Command{
//	Use:     "lso",
//	Short:   "lso",
//	Long:    "List object",
//	Example: "gcscmd ls cs://BUCKET [--name=<name>] [--cid=<cid>] [--Offset=<Offset>]",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		//cmd.Help()
//		//fmt.Printf("%s %s\n", cmd.Name(), strconv.Itoa(offset))
//		objectListRun(cmd, args)
//	},
//}

func objectListRun(cmd *cobra.Command, args []string) {
	itemName := ""
	pageSize := 10
	pageIndex := 1

	// 桶名称
	bucketName := GetBucketName(args)
	if err := checkBucketName(bucketName); err != nil {
		Error(cmd, args, err)
	}

	// 对象名称
	objectName, err := cmd.Flags().GetString("name")
	if err != nil {
		Error(cmd, args, err)
	}

	if len(objectName) > 0 {
		if err := checkObjectName(objectName); err != nil {
			Error(cmd, args, err)
		}
	}

	// 对象CID
	objectCid, err := cmd.Flags().GetString("cid")
	if err != nil {
		Error(cmd, args, err)
	}

	// 设置参数
	if len(objectName) > 0 {
		itemName = objectName
	} else if len(objectCid) > 0 {
		itemName = objectCid
	}

	// 查询偏移量
	offset := viper.GetInt("cmd.list_offset")
	if offset > 0 || offset < 1000 {
		pageSize = offset
	}

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

	// 列出桶对象
	response, err := sdk.Object.GetObjectList(bucketId, itemName, pageSize, pageIndex)
	if err != nil {
		Error(cmd, args, err)
	}

	objectListRunOutput(cmd, args, response)
}

func objectListRunOutput(cmd *cobra.Command, args []string, resp model.ObjectPageResponse) {
	code := int(resp.Code)
	if code != http.StatusOK {
		Error(cmd, args, errors.New(resp.Msg))

		//err := errors.Errorf("code:%d, message:&s\n", resp.Code, resp.Msg)
		//if code == sdkcode.ErrInvalidBucketId.Code() {
		//	err = errors.Errorf("bucket can't be found")
		//}

		//Error(cmd, args, err)
	}

	respData := resp.Data
	objectListOutput := ObjectListOutput{
		RequestId: resp.RequestId,
		Code:      resp.Code,
		Msg:       resp.Msg,
		Status:    resp.Status,
		Count:     respData.Count,
		PageIndex: respData.PageIndex,
		PageSize:  respData.PageSize,
		List:      []ObjectOutput{},
	}

	if len(respData.List) > 0 {
		for i, _ := range respData.List {
			objectOutput := ObjectOutput{}
			deepcopier.Copy(respData.List[i]).To(&objectOutput)

			// 创建时间
			objectOutput.CreatedDate = objectOutput.CreatedAt.Format("2006-01-02")
			objectListOutput.List = append(objectListOutput.List, objectOutput)
		}
	}

	templateContent := `
total {{.Count}}
{{- if eq (len .List) 0}}
Status: {{.Code}}
{{else}}
{{- range .List}}
{{.ObjectCid}} {{.ObjectSize}} {{.CreatedDate}} {{.ObjectName}}
{{ end}}
{{- end}}`

	t, err := template.New("objectListTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, objectListOutput)
	if err != nil {
		Error(cmd, args, err)
	}
}

type ObjectListOutput struct {
	RequestId string         `json:"requestId,omitempty"`
	Code      int32          `json:"code,omitempty"`
	Msg       string         `json:"msg,omitempty"`
	Status    string         `json:"status,omitempty"`
	Count     int            `json:"count,omitempty"`
	PageIndex int            `json:"pageIndex,omitempty"`
	PageSize  int            `json:"pageSize,omitempty"`
	List      []ObjectOutput `json:"list,omitempty"`
}

type ObjectOutput struct {
	Id             int       `json:"id" comment:"对象ID"`
	BucketId       int       `json:"bucketId" comment:"桶主键"`
	ObjectName     string    `json:"objectName" comment:"对象名称（255字限制）"`
	ObjectTypeCode int       `json:"objectTypeCode" comment:"对象类型编码"`
	ObjectSize     int64     `json:"objectSize" comment:"对象大小（字节）"`
	IsMarked       int       `json:"isMarked" comment:"星标（1-已标记，0-未标记）"`
	ObjectCid      string    `json:"objectCid" comment:"对象CID"`
	CreatedAt      time.Time `json:"createdAt" comment:"创建时间"`
	UpdatedAt      time.Time `json:"updatedAt" comment:"最后更新时间"`
	CreatedDate    string    `json:"createdDate" comment:"创建日期"`
}

// endregion Object List

// region Object Rename

//var objectRenameCmd = &cobra.Command{
//	Use:     "rn",
//	Short:   "rn",
//	Long:    "rename object",
//	Example: "gcscmd rn cs://BUCKET] [--name=<name>] [--cid=<cid>] [--rename=<rename>] [--force]",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		objectRenameRun(cmd, args)
//	},
//}

func objectRenameRun(cmd *cobra.Command, args []string) {
	// 桶名称
	bucketName := GetBucketName(args)
	if err := checkBucketName(bucketName); err != nil {
		Error(cmd, args, err)
	}

	// 对象名称
	objectName, err := cmd.Flags().GetString("name")
	if err != nil {
		Error(cmd, args, err)
	}

	if err := checkObjectName(objectName); err != nil {
		Error(cmd, args, err)
	}

	//// 对象CID
	//objectCid, err := cmd.Flags().GetString("cid")
	//if err != nil {
	//	Error(cmd, args, err)
	//}

	// 重命名
	rename, err := cmd.Flags().GetString("rename")
	if err != nil {
		Error(cmd, args, err)
	}

	if err := checkObjectName(rename); err != nil {
		Error(cmd, args, err)
	}

	// todo: return succeed?
	if rename == objectName {
		Error(cmd, args, errors.New("the new name of object can't be equal to the raw name of object"))
	}

	// 强制覆盖
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		Error(cmd, args, err)
	}

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

	// 确认对象数据有效性
	respObject, err := sdk.Object.GetObjectByName(bucketId, objectName)
	if err != nil {
		Error(cmd, args, err)
	}

	code = int(respObject.Code)
	if code != http.StatusOK {
		Error(cmd, args, errors.New(respObject.Msg))
	}

	// 对象ID
	objectId := respObject.Data.Id

	// 重命名对象
	response, err := sdk.Object.RenameObject(objectId, rename, force)
	if err != nil {
		Error(cmd, args, err)
	}

	objectRenameRunOutput(cmd, args, response)
}

func objectRenameRunOutput(cmd *cobra.Command, args []string, resp model.ObjectRenameResponse) {
	respCode := int(resp.Code)

	if respCode == sdkcode.ErrObjectNameConflictInBucket.Code() {
		err := errors.New("Error: conflicting rename filename, add --force to confirm overwrite\n")
		Error(cmd, args, err)
	} else if respCode != http.StatusOK {
		Error(cmd, args, errors.New(resp.Msg))
	}

	objectRenameOutput := ObjectRenameOutput{
		RequestId: resp.RequestId,
		Code:      resp.Code,
		Msg:       resp.Msg,
		Status:    resp.Status,
	}

	templateContent := `
Succeed
Status: {{.Code}}
`

	t, err := template.New("objectRenameTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, objectRenameOutput)
	if err != nil {
		Error(cmd, args, err)
	}
}

type ObjectRenameOutput struct {
	RequestId string       `json:"requestId,omitempty"`
	Code      int32        `json:"code,omitempty"`
	Msg       string       `json:"msg,omitempty"`
	Status    string       `json:"status,omitempty"`
	Data      ObjectOutput `json:"objectOutput,omitempty"`
}

// endregion Object Rename

// region Object Remove

//var objectRemoveCmd = &cobra.Command{
//	Use:     "rmo",
//	Short:   "rmo",
//	Long:    "remove object",
//	Example: "gcscmd rmo cs://BUCKET] [--name=<name>] [--cid=<cid>] [--remove=<remove>] [--force]",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		objectRemoveRun(cmd, args)
//	},
//}

func objectRemoveRun(cmd *cobra.Command, args []string) {
	// 桶名称
	bucketName := GetBucketName(args)
	if err := checkBucketName(bucketName); err != nil {
		Error(cmd, args, err)
	}

	// 对象名称
	objectName, err := cmd.Flags().GetString("name")
	if err != nil {
		Error(cmd, args, err)
	}

	if err := checkObjectName(objectName); err != nil {
		Error(cmd, args, err)
	}

	//// 对象CID
	//objectCid, err := cmd.Flags().GetString("cid")
	//if err != nil {
	//	Error(cmd, args, err)
	//}

	//// 强制覆盖
	//force, err := cmd.Flags().GetBool("force")
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

	// 确认对象数据有效性
	respObject, err := sdk.Object.GetObjectByName(bucketId, objectName)
	if err != nil {
		Error(cmd, args, err)
	}

	code = int(respObject.Code)
	if code != http.StatusOK {
		Error(cmd, args, errors.New(respObject.Msg))
	}

	// 对象ID
	objectId := respObject.Data.Id
	objectIdList := []int{objectId}

	//// todo: 批量删除? cid or objectName
	//itemName := ""
	//pageSize := 1000
	//pageIndex := 1
	//
	//// 设置参数
	//if len(objectName) > 0 {
	//	itemName = objectName
	//} else if len(objectCid) > 0 {
	//	itemName = objectCid
	//}
	//
	//pageRespObject, err := sdk.Object.GetObjectList(bucketId, itemName, pageSize, pageIndex)
	//if err != nil {
	//	Error(cmd, args, err)
	//}
	//
	//code = int(pageRespObject.Code)
	//if code != http.StatusOK {
	//	Error(cmd, args, errors.New(pageRespObject.Msg))
	//}
	//
	//objectIdList := []int{}
	//for i := range pageRespObject.Data.List {
	//	objectIdList = append(objectIdList, pageRespObject.Data.List[i].Id)
	//}

	// 重命名对象
	response, err := sdk.Object.RemoveObject(objectIdList)
	if err != nil {
		Error(cmd, args, err)
	}

	objectRemoveRunOutput(cmd, args, response)
}

func objectRemoveRunOutput(cmd *cobra.Command, args []string, resp model.ObjectRemoveResponse) {
	respCode := int(resp.Code)

	if respCode != http.StatusOK {
		Error(cmd, args, errors.New(resp.Msg))
	}

	objectRemoveOutput := ObjectRemoveOutput{
		RequestId: resp.RequestId,
		Code:      resp.Code,
		Msg:       resp.Msg,
		Status:    resp.Status,
	}

	//	删除对象
	//	通过命令删除固定桶内对象
	//
	//	模版
	//
	//	gcscmd rm cs://[BUCKET] [--name=<name>] [--cid=<cid>] [--force]
	//	BUCKET
	//
	//	桶名称
	//
	//	cid
	//
	//	添加对应的 CID
	//
	//	name
	//
	//	对象名
	//
	//	force
	//
	//	无添加筛选条件或命中多的对象时需要添加
	//
	//	命令行例子
	//
	//	清空桶
	//
	//	gcscmd rm cs://bbb --force
	//	使用对象名删除单文件
	//
	//	gcscmd rm  cs://bbb --name Tarkov.mp4
	//	使用模糊查询删除对象
	//
	//	gcscmd rm  cs://bbb --name .mp4 --force
	//	使用对象名删除单目录
	//
	//	gcscmd rm  cs://bbb --name aaa
	//	使用CID删除单对象
	//
	//	gcscmd rm  cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo
	//	使用 CID 删除多个对象(命中多个对象时加)
	//
	//	gcscmd rm  cs://bbb --cid QmWgnG7pPjG31w328hZyALQ2BgW5aQrZyKpT47jVpn8CNo --force
	//	响应
	//
	//	成功
	//
	//Status: 200
	//	多对象没有添加 force
	//
	//Error: multiple object  are matching this query, add --force to confirm the bulk removal

	templateContent := `
Succeed
Status: {{.Code}}
`

	t, err := template.New("objectRemoveTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, objectRemoveOutput)
	if err != nil {
		Error(cmd, args, err)
	}
}

type ObjectRemoveOutput struct {
	RequestId string       `json:"requestId,omitempty"`
	Code      int32        `json:"code,omitempty"`
	Msg       string       `json:"msg,omitempty"`
	Status    string       `json:"status,omitempty"`
	Data      ObjectOutput `json:"objectOutput,omitempty"`
}

// endregion Object Remove

// region Object Download

//var objectDownloadCmd = &cobra.Command{
//	Use:     "get",
//	Short:   "get",
//	Long:    "download object",
//	Example: "gcscmd get cs://BUCKET [--name=<name>] [--cid=<cid>]",
//
//	Run: func(cmd *cobra.Command, args []string) {
//		objectDownloadRun(cmd, args)
//	},
//}

func objectDownloadRun(cmd *cobra.Command, args []string) {
	// 桶名称
	bucketName := GetBucketName(args)
	if err := checkBucketName(bucketName); err != nil {
		Error(cmd, args, err)
	}

	// 对象名称
	objectName, err := cmd.Flags().GetString("name")
	if err != nil {
		Error(cmd, args, err)
	}

	if err := checkObjectName(objectName); err != nil {
		Error(cmd, args, err)
	}

	//// 对象CID
	//objectCid, err := cmd.Flags().GetString("cid")
	//if err != nil {
	//	Error(cmd, args, err)
	//}

	//// 输出路径
	//target, err := cmd.Flags().GetString("Target")
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

	// 确认对象数据有效性
	respObject, err := sdk.Object.GetObjectByName(bucketId, objectName)
	if err != nil {
		Error(cmd, args, err)
	}

	code = int(respObject.Code)
	if code != http.StatusOK {
		Error(cmd, args, errors.New(respObject.Msg))
	}

	// 对象ID
	//objectId := respObject.Data.Id
	//objectName := respObject.Data.ObjectName
	objectCid := respObject.Data.ObjectCid
	//downloadEndpoint := "https://test-ipfs-gateway.netwarps.com/ipfs/"
	//downloadUrl := fmt.Sprintf("%s%s", downloadEndpoint, objectCid)
	ipfsGateway := viper.GetString("cmd.ipfs_gateway")
	downloadUrl := fmt.Sprintf("https://%s%s", ipfsGateway, objectCid)

	// todo: remove it
	//downloadUrl = "https://test-ipfs-gateway.netwarps.com/ipfs/bafybeiguyrqm6z76mrhntk64fiwwdpjqv64ny3ugw64owznlbeotknvypa"
	cli := pget.New()
	cli.URLs = []string{downloadUrl}
	cli.Output = objectName

	version := ""
	downloadArgs := []string{"-t", "10"}

	if err := cli.Run(context.Background(), version, downloadArgs); err != nil {
		//if cli.Trace {
		//	fmt.Fprintf(os.Stderr, "Error:\n%+v\n", err)
		//} else {
		//	fmt.Fprintf(os.Stderr, "Error:\n  %v\n", err)
		//}
		Error(cmd, args, err)
	}

	objectDownloadRunOutput(cmd, args, respObject)
}

func objectDownloadRunOutput(cmd *cobra.Command, args []string, resp model.ObjectCreateResponse) {
	respCode := int(resp.Code)

	if respCode != http.StatusOK {
		Error(cmd, args, errors.New(resp.Msg))
	}

	objectDownloadOutput := ObjectDownloadOutput{
		RequestId: resp.RequestId,
		Code:      resp.Code,
		Msg:       resp.Msg,
		Status:    resp.Status,
	}

	err := deepcopier.Copy(&resp.Data).To(&objectDownloadOutput.Data)
	if err != nil {
		Error(cmd, args, err)
	}

	templateContent := `
CID: {{.Data.ObjectCid}}
Name: {{.Data.ObjectName}}
`

	t, err := template.New("objectDownloadTemplate").Parse(templateContent)
	if err != nil {
		Error(cmd, args, err)
	}

	err = t.Execute(os.Stdout, objectDownloadOutput)
	if err != nil {
		Error(cmd, args, err)
	}
}

type ObjectDownloadOutput struct {
	RequestId string       `json:"requestId,omitempty"`
	Code      int32        `json:"code,omitempty"`
	Msg       string       `json:"msg,omitempty"`
	Status    string       `json:"status,omitempty"`
	Data      ObjectOutput `json:"objectOutput,omitempty"`
}

// endregion Object Download
