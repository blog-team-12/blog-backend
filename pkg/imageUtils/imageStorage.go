package imageUtils

import (
	"context"
	"errors"
	"io"
	"personal_blog/global"
	"personal_blog/internal/model/consts"
	"personal_blog/internal/model/entity"
	"personal_blog/pkg/storage"
	"strings"
)

// ImageURL 基于模型生成可访问URL
// - 优先使用模型中的 URL
// - 否则使用静态前缀 Prefix 与 Key 组装
func ImageURL(img *entity.Image) string {
	if img == nil {
		return ""
	}
	if u := strings.TrimSpace(img.URL); u != "" {
		return u
	}
	prefix := strings.TrimSuffix(strings.TrimSpace(global.Config.Static.Prefix), "/")
	key := strings.TrimPrefix(strings.TrimSpace(img.Key), "/")
	if prefix == "" {
		return "/" + key
	}
	return prefix + "/" + key
}

// ApplyStorageObject 将上传返回的 StorageObject 写入模型字段
// - Key/URL/Name 直接回写
// - Storage 根据驱动名映射到枚举
func ApplyStorageObject(img *entity.Image, obj storage.StorageObject, driverName string) {
	if img == nil {
		return
	}
	img.Key = obj.Key
	img.URL = obj.URL
	if obj.Name != "" {
		img.Name = obj.Name
	}
	switch strings.ToLower(strings.TrimSpace(driverName)) {
	case "qiniu":
		img.Storage = consts.Qiniu
	default:
		img.Storage = consts.Local
	}
}

// UploadViaCurrent 使用当前驱动执行上传并返回存储对象
func UploadViaCurrent(ctx context.Context, r io.Reader, filename string) (storage.StorageObject, error) {
	drv := storage.Current()
	if drv == nil {
		return storage.StorageObject{}, ErrNoDriver
	}
	return drv.Upload(ctx, r, filename)
}

// ErrNoDriver 表示当前未初始化存储驱动
var ErrNoDriver = errors.New("no storage driver initialized")

// UploadAndBind 使用当前驱动上传并将结果绑定到模型
func UploadAndBind(ctx context.Context, img *entity.Image, r io.Reader, filename string) (storage.StorageObject, error) {
	obj, err := UploadViaCurrent(ctx, r, filename)
	if err != nil {
		return storage.StorageObject{}, err
	}
	drv := storage.Current()
	name := ""
	if drv != nil {
		name = drv.Name()
	}
	ApplyStorageObject(img, obj, name)
	return obj, nil
}

// UploadWithDriver 指定驱动名执行上传并绑定模型
func UploadWithDriver(
	ctx context.Context,
	driverName string,
	img *entity.Image,
	r io.Reader,
	filename string,
) (storage.StorageObject, error) {
	drv := storage.FromName(driverName)
	if drv == nil {
		return storage.StorageObject{}, ErrNoDriver
	}
	obj, err := drv.Upload(ctx, r, filename)
	if err != nil {
		return storage.StorageObject{}, err
	}
	ApplyStorageObject(img, obj, drv.Name())
	return obj, nil
}
