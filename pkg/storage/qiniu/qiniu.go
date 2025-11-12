package qiniu

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"io"
	"path/filepath"
	"personal_blog/global"
	"personal_blog/pkg/storage"
	"strings"
	"time"

	qstorage "github.com/qiniu/go-sdk/v7/storage"
)

// Driver 七牛云存储驱动
// 命名为 Driver 以避免包名重复导致的命名冗余（revive stutter）。
type Driver struct{}

// New 创建并返回一个七牛云存储驱动实例。
func New() *Driver { return &Driver{} }

// Name 返回驱动名称，用于注册与选择。
func (d *Driver) Name() string { return "qiniu" }

// 在包初始化时注册七牛驱动单例
func init() {
	storage.RegisterDriver("qiniu", New())
}

// Delete 删除对象；若返回 612（资源不存在）视为成功
func (d *Driver) Delete(_ context.Context, key string) error {
	cfg := sdkConfig()
	mac, bucket, err := credentialsFromConfig()
	if err != nil {
		return err
	}
	bm := qstorage.NewBucketManager(mac, cfg)
	if err := bm.Delete(bucket, key); err != nil {
		// 七牛错误码 612 表示资源不存在，视为成功
		if strings.Contains(err.Error(), "612") {
			return nil
		}
		return err
	}
	return nil
}

// Upload 使用 Resumable V2 接口进行分片上传，支持 io.Reader
// 返回 StorageObject，其中 URL 使用配置的自定义域名
func (d *Driver) Upload(
	ctx context.Context,
	r io.Reader,
	filename string,
) (storage.StorageObject, error) {
	cfg := sdkConfig()
	mac, bucket, err := credentialsFromConfig()
	if err != nil {
		return storage.StorageObject{}, err
	}

	// 生成对象键：<prefix>/<yyyyMMdd>/<rand16><ext>
	key := generateObjectKey(filename)
	if kp := strings.Trim(global.Config.Storage.Qiniu.KeyPrefix, "/"); kp != "" {
		key = filepath.ToSlash(filepath.Join(kp, key))
	}

	// 上传凭证
	putPolicy := qstorage.PutPolicy{Scope: bucket}
	upToken := putPolicy.UploadToken(mac)

	// 统计上传字节数（兼容返回 Size）
	counter := &byteCounter{}
	reader := io.TeeReader(r, counter)

	ru := qstorage.NewResumeUploaderV2(cfg)

	var putRet qstorage.PutRet
	uploadErr := storage.DoWithBackoff( // 通过重复/抖动，从而自动恢复一些临时错误
		ctx,
		3,
		300*time.Millisecond,
		200*time.Millisecond,
		func() error {
			// 无需提前知道大小，走 PutWithoutSize 以支持流式上传
			// 无需加载整个文件到内存，所以支持超大文件（理论无限制）
			return ru.PutWithoutSize(ctx, &putRet, upToken, key, reader, &qstorage.RputV2Extra{})
		},
	)
	if uploadErr != nil {
		return storage.StorageObject{}, uploadErr
	}

	// 拼接访问 URL（优先使用配置 domain）
	domain := strings.TrimSuffix(global.Config.Storage.Qiniu.Domain, "/")
	url := domain
	if url == "" {
		// 后备：若未配置 domain，使用 http://<system.host>:<system.port><static.prefix>
		base := strings.TrimSuffix(global.Config.System.Host, "/")
		port := global.Config.System.Port
		prefix := strings.TrimSuffix(global.Config.Static.Prefix, "/")
		url = base
		if port > 0 {
			url = url + ":" + fmt.Sprintf("%d", port)
		}
		url = url + prefix
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	fullURL := strings.TrimSuffix(url, "/") + "/" + key

	return storage.StorageObject{
		Key:  key,
		URL:  fullURL,
		Size: int64(counter.n),
		Type: "",
		Name: filename,
	}, nil
}

// sdkConfig 构建七牛 SDK 配置，挂载共享 HTTP 连接的 Transport
func sdkConfig() *qstorage.Config {
	useHTTPS := strings.HasPrefix(global.Config.Storage.Qiniu.Domain, "https://")
	// https://cdn.example.co
	return &qstorage.Config{
		UseHTTPS: useHTTPS,
	}
}

// credentialsFromConfig 读取 AK/SK 与 bucket
// 优先从配置指定的环境变量名读取；若环境变量未设置，则回退直接使用配置中的值。
// 这样可以复用现有字段：AccessKeyEnv / SecretKeyEnv
func credentialsFromConfig() (*auth.Credentials, string, error) {
	q := global.Config.Storage.Qiniu

	ak := strings.TrimSpace(q.AccessKey)
	sk := strings.TrimSpace(q.SecretKey)
	if ak == "" || sk == "" {
		return nil, "", errors.New("qiniu credentials not configured")
	}
	if strings.TrimSpace(q.Bucket) == "" {
		return nil, "", errors.New("qiniu bucket not configured")
	}
	mac := auth.New(ak, sk)
	return mac, q.Bucket, nil
}

// generateObjectKey 生成符合规范的对象键
func generateObjectKey(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".bin"
	}
	b := make([]byte, 16)
	// 避免未检查错误触发 errcheck；失败时退回时间戳作为随机因子
	if _, err := rand.Read(b); err != nil {
		// 使用当前时间与文件名混合生成退路值
		fallback := time.Now().UnixNano()
		b = []byte(fmt.Sprintf("%x", fallback))
	}
	randHex := hex.EncodeToString(b)
	dateDir := time.Now().Format("20060102")
	return filepath.ToSlash(filepath.Join(dateDir, randHex+ext))
}

// byteCounter 用于统计上传字节数（配合 TeeReader）
type byteCounter struct{ n int }

func (b *byteCounter) Write(p []byte) (int, error) {
	b.n += len(p)
	return len(p), nil
}
