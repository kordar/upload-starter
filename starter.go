package upload_starter

import (
	"github.com/kordar/goframework-upload"
	logger "github.com/kordar/gologger"
	"github.com/kordar/goupload_cos"
	"github.com/kordar/goupload_local"
	"github.com/spf13/cast"
	"io/fs"
	"path"
	"strings"
)

type UploadModule struct {
	name string
	args map[string]interface{}
	load func(moduleName string, itemId string, item map[string]string)
}

func NewUploadModule(name string, load func(moduleName string, itemId string, item map[string]string), args map[string]interface{}) *UploadModule {
	return &UploadModule{name, args, load}
}

func (m UploadModule) Name() string {
	return m.name
}

func (m UploadModule) _load(id string, cfg map[string]string) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	driver := cfg["driver"]
	if driver == "cos" {
		if cfg["bucket"] == "" || cfg["region"] == "" || cfg["secret_key"] == "" || cfg["secret_id"] == "" {
			logger.Fatalf("[%s] invalid client parameters for instance cos, bucket,region,secret_id,secret_key are required.", m.Name())
		}
		client := goupload_cos.NewCOSClient(cfg["bucket"], cfg["region"], cfg["secret_id"], cfg["secret_key"])
		if err := goframework_upload.AddUploaderInstance(id, client); err != nil {
			logger.Error("[%s] create cos fail, err=%v", err)
		}
	}

	if driver == "local" {
		if cfg["bucket"] == "" || cfg["root"] == "" {
			logger.Fatalf("[%s] invalid client parameters for instance local, bucket,root are required.", m.Name())
		}

		// 默认过滤隐藏文件
		filterDirItem := func(s string, entry fs.DirEntry) bool {
			base := path.Base(s)
			return strings.HasPrefix(base, ".")
		}

		if m.args != nil && m.args["filter"] != nil {
			filterDirItem = m.args["filter"].(goupload_local.FilterDirItem)
		}

		client := goupload_local.NewLocalUploader(cfg["root"], cfg["bucket"], filterDirItem)
		if err := goframework_upload.AddUploaderInstance(id, client); err != nil {
			logger.Error("[%s] create local fail, err=%v", err)
		}

	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m UploadModule) Load(value interface{}) {
	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, cast.ToStringMapString(value))
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMapString(item))
	}

}

func (m UploadModule) Close() {
}
