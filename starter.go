package upload_starter

import (
	"github.com/kordar/goupload"
	"github.com/kordar/goupload_cos"
	"github.com/spf13/cast"
)

var (
	bucketManager *goupload.BucketManager
)

func GetUploadMgr() *goupload.BucketManager {
	return bucketManager
}

func GetUploaderByBucket(name string) goupload.IUpload {
	return bucketManager.GetHandler(name)
}

type UploadModule struct {
}

func (m UploadModule) Name() string {
	return "upload"
}

func (m UploadModule) Close() {
}

func (m UploadModule) Load(value interface{}) {
	bucketManager = goupload.NewBucketManager()
	items := cast.ToStringMap(value)
	for bucket, item := range items {
		cfg := cast.ToStringMap(item)
		switch cfg["driver"] {
		case "cos":
			m.createCos(bucket, cfg)
			break
		default:

		}
	}
}

func (m UploadModule) createCos(bucket string, cfg map[string]interface{}) {
	region := cast.ToString(cfg["region"])
	secretId := cast.ToString(cfg["secretId"])
	secretKey := cast.ToString(cfg["secretKey"])
	client := goupload_cos.NewCOSClient(bucket, region, secretId, secretKey)
	bucketManager.SetUploadHandlers(client)
	auto := cast.ToBool(cfg["autoInit"])
	if auto {
		client.CreateBucket()
	}
}
