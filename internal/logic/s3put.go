package logic

import (
	"context"
	"errors"
	"io"
	"os"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"

	"github.com/steemfans/multiple-upload-to-s3/internal/model/entity"
	"github.com/steemfans/multiple-upload-to-s3/internal/service/aws"
	"github.com/steemfans/multiple-upload-to-s3/internal/service/sqlite"
	"github.com/steemfans/multiple-upload-to-s3/utility"
)

// file exist -- true;
// file not exist -- false.
func CheckDb(path, dbName string) (result bool) {
	if _, err := os.Stat(path + "/" + dbName); err == nil {
		return true
	}
	return false
}

// taskId is 0 means new task;
// taskId is not 0 and lastPartNum is 0 means no upload part;
// taskId and lastPartNum are not 0 means there is an unfinished task.
func CheckTask(ctx context.Context, bucketName, fileKey, src string) (taskId, lastPartNum int, uploadId string, err error) {
	db, err := gdb.Instance()
	if err != nil {
		glog.Warning(ctx, "Get db instance failed:", err)
		return
	}
	var task *entity.Tasks
	err = db.Model("tasks").Where("bucket_name = ? and file_key = ? and src = ?", bucketName, fileKey, src).Scan(&task)
	if err != nil {
		glog.Warning(ctx, "Get task failed:", err)
		return
	}
	if task == nil {
		err = errors.New("task_empty")
		glog.Info(ctx, "task empty.")
		return
	}
	taskId = task.Id
	uploadId = task.UploadId
	var part *entity.Parts
	err = db.Model("parts").Where("task_id = ?", taskId).Order("part_num", "desc").Limit(1).Scan(&part)
	if err != nil {
		glog.Warning(ctx, "Get part failed:", err)
		return
	}
	if part == nil {
		glog.Info(ctx, "part empty.")
		return
	}
	lastPartNum = part.PartNum
	return
}

func CreateTask(ctx context.Context, uploadId, bucketName, fileKey, src string) (id int, err error) {
	db, err := gdb.Instance()
	if err != nil {
		glog.Warning(ctx, "Get db instance failed:", err)
		return
	}

	result, err := db.Model("tasks").Insert(g.Map{
		"upload_id":   uploadId,
		"bucket_name": bucketName,
		"file_key":    fileKey,
		"src":         src,
	})
	if err != nil {
		glog.Warning(ctx, "Get task failed:", err)
		return
	}
	tmpId, err := result.LastInsertId()
	return int(tmpId), err
}

func Upload(ctx context.Context, taskId int, partNum int, fileContent []byte) (err error) {
	db, err := gdb.Instance()
	if err != nil {
		glog.Warning(ctx, "Get db instance failed:", err)
		return
	}
	var task *entity.Tasks
	err = db.Model("tasks").Where("id = ?", taskId).Scan(&task)
	if err != nil {
		glog.Warning(ctx, "Get task failed:", err)
		return
	}
	if task == nil {
		err = errors.New("task_empty")
		glog.Info(ctx, "task empty.")
		return
	}

	var out = s3.CreateMultipartUploadOutput{
		Bucket:   awsSDK.String(task.BucketName),
		Key:      awsSDK.String(task.FileKey),
		UploadId: awsSDK.String(task.UploadId),
	}
	result, err := aws.SS3.UploadPart(ctx, &out, fileContent, int32(partNum))
	if err != nil {
		glog.Fatal(ctx, "UploadPart failed.", err)
	}
	_, err = db.Model("parts").Insert(g.Map{
		"task_id":     taskId,
		"part_num":    partNum,
		"content_md5": utility.Md5(fileContent),
		"etag":        result.ETag,
		"status":      1,
	})
	return
}

func CompleteTask(ctx context.Context, taskId int) (result *s3.CompleteMultipartUploadOutput, err error) {
	var completedParts []s3types.CompletedPart
	var parts []entity.Parts
	var task entity.Tasks

	db, err := gdb.Instance()
	if err != nil {
		glog.Warning(ctx, "Get db instance failed:", err)
		return
	}

	err = db.Model("tasks").Where("id = ?", taskId).Scan(&task)
	if err != nil {
		glog.Warning(ctx, "Get task failed:", err)
		return
	}

	err = db.Model("parts").Where("task_id = ?", taskId).Order("part_num", "asc").Scan(&parts)
	if err != nil {
		glog.Warning(ctx, "Get part failed:", err)
		return
	}

	for i, p := range parts {
		completedParts[i] = s3types.CompletedPart{
			ETag:       awsSDK.String(p.Etag),
			PartNumber: int32(p.PartNum),
		}
	}

	var out = s3.CreateMultipartUploadOutput{
		Bucket:   awsSDK.String(task.BucketName),
		Key:      awsSDK.String(task.FileKey),
		UploadId: awsSDK.String(task.UploadId),
	}

	return aws.SS3.CompleteMultipartUpload(ctx, &out, completedParts)
}

func S3Put(ctx context.Context, parser *gcmd.Parser) (err error) {
	bucketName := parser.GetOpt("bucketname").String()
	objectName := parser.GetOpt("objectname").String()
	src := parser.GetOpt("src").String()
	endpoint := os.Getenv("ENDPOINT")
	currentWorkPath := gfile.Pwd()
	// BLOCK_SIZE unit is MB, blockSize unit is Byte.
	blockSize := genv.Get("BLOCK_SIZE", 50).Int64() * 1024 * 1024
	dbName := sqlite.GetDbName(bucketName, objectName)

	glog.Info(ctx, src, endpoint)
	// check if db exist
	if !CheckDb(currentWorkPath, dbName) {
		// create db
		err = sqlite.GenerateNewDb(currentWorkPath, dbName)
		if err != nil {
			glog.Warning(ctx, "Create db failed.", err)
			return
		}
	}

	// init db config
	err = sqlite.InitUploadTaskDbConfig(currentWorkPath + "/" + dbName)
	if err != nil {
		glog.Warning(ctx, "Init db config failed.", err)
		return
	}

	var out *s3.CreateMultipartUploadOutput
	var uploadId string
	var taskId, lastPartNum int
	// check if new task
	taskId, lastPartNum, _, err = CheckTask(ctx, bucketName, objectName, src)
	if taskId == 0 {
		// create new task
		out, err = aws.SS3.CreateMultipartUpload(ctx, bucketName, objectName, src)
		if err != nil {
			glog.Warning(ctx, "Create Multipart Upload failed.", err)
			return
		}
		uploadId = awsSDK.ToString(out.UploadId)
		taskId, err = CreateTask(ctx, uploadId, bucketName, objectName, src)
		if err != nil {
			glog.Warning(ctx, "Create task db data failed.", err)
			return
		}
	}

	file, err := os.Open(src)
	if err != nil {
		glog.Warning(ctx, "Error to open file.", err)
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buf := make([]byte, blockSize)

	var offset int64

	currentPartNum := 1
	offset = 0

	if lastPartNum != 0 {
		glog.Info(ctx, "Start last upload, please wait.")
		currentPartNum = lastPartNum + 1
		offset = blockSize * int64(lastPartNum)
		if offset >= fileSize {
			CompleteTask(ctx, taskId)
			return
		}
	}

	file.Seek(offset, 0)

	for offset < fileSize {
		_, err = file.Read(buf)
		if err == io.EOF {
			CompleteTask(ctx, taskId)
			return
		}
		glog.Debug(ctx, "process:", offset, currentPartNum)
		err = Upload(ctx, taskId, currentPartNum, buf)
		if err != nil {
			glog.Fatal(ctx, "Upload failed", err)
		}
		currentPartNum += 1
		offset += blockSize
		if offset >= fileSize {
			CompleteTask(ctx, taskId)
			return
		}
		_, err = file.Seek(offset, 0)
		if err != nil {
			glog.Fatal(ctx, "File seek failed.", err)
		}
	}

	return
}
