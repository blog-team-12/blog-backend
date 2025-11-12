package local

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"personal_blog/global"
	"personal_blog/pkg/storage"
	"strings"
	"time"
)

// Driver 本地文件存储驱动
// 命名为 Driver 以避免包名重复导致的命名冗余（revive stutter）。
type Driver struct{}

// New 创建并返回一个本地存储驱动实例。
func New() *Driver { return &Driver{} }

// Name 返回驱动名称，用于注册与选择。
func (d *Driver) Name() string { return "local" }

// 在包初始化时注册本地驱动单例
func init() {
	storage.RegisterDriver("local", New())
}

// Delete 删除本地文件；资源不存在视为成功
func (d *Driver) Delete(_ context.Context, key string) error {
	// 兼容现有静态目录：使用配置中的静态根路径进行拼接
	root := strings.TrimSpace(global.Config.Static.Path)
	realPath := filepath.Join(root, key)
	if err := os.Remove(realPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

// Upload 暂不改造上传流程，保留骨架，后续可切换到统一驱动上传
func (d *Driver) Upload(
	_ context.Context,
	r io.Reader,
	filename string,
) (storage.StorageObject, error) {
	key, err := generateLocalKey(filename)
	if err != nil {
		return storage.StorageObject{}, err
	}

	// 目标路径
	root := strings.TrimSpace(global.Config.Static.Path)
	realPath := filepath.Join(root, key)
	if err := os.MkdirAll(filepath.Dir(realPath), 0755); err != nil {
		return storage.StorageObject{}, err
	}

	// 写入文件（流式复制）
	f, err := os.Create(realPath)
	if err != nil {
		return storage.StorageObject{}, err
	}
	defer func() { _ = f.Close() }()
	n, err := io.Copy(f, r)
	if err != nil {
		return storage.StorageObject{}, err
	}

	baseURL := composeBaseURL()
	fullURL := strings.TrimSuffix(baseURL, "/") + "/" + key

	obj := storage.StorageObject{
		Key:  key,
		URL:  fullURL,
		Size: n,
		Type: "",
		Name: filename,
	}
	return obj, nil
}

// generateLocalKey 生成唯一键：<yyyyMMdd>/<rand16><ext>，并按需添加业务前缀
func generateLocalKey(filename string) (string, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".bin"
	}
	ext = strings.ToLower(ext)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	randHex := hex.EncodeToString(b)
	dateDir := time.Now().Format("20060102")
	key := filepath.ToSlash(filepath.Join(dateDir, randHex+ext))
	if kp := strings.Trim(global.Config.Storage.Local.KeyPrefix, "/"); kp != "" {
		key = filepath.ToSlash(filepath.Join(kp, key))
	}
	return key, nil
}

// composeBaseURL 计算本地驱动返回的基础 URL
func composeBaseURL() string {
	baseURL := strings.TrimSuffix(global.Config.Storage.Local.BaseURL, "/")
	if baseURL == "" {
		host := strings.TrimSuffix(global.Config.System.Host, "/")
		port := global.Config.System.Port
		prefix := strings.TrimSuffix(global.Config.Static.Prefix, "/")
		baseURL = host
		if port > 0 {
			baseURL = baseURL + ":" + fmt.Sprintf("%d", port)
		}
		baseURL = baseURL + prefix
	}
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	return baseURL
}
