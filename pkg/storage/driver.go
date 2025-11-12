package storage

import (
	"context"
	"io"
)

// StorageObject 统一描述一次存储操作的结果与元信息。
// 约定与语义：
// - Key：驱动内部的“对象名/相对路径”，不包含前缀，例如 local 不含 "/images/"；qiniu 为桶内对象键。
// - URL：可直接用于前端访问的路径或完整地址；local 为 "/images/" + Key；qiniu 为 "https://domain/" + Key。
// - Size/Type/Name：可选元信息，供业务使用；驱动可按需填充（不强制）。
//
//revive:disable:exported // 保留 StorageObject 命名以明确语义与外部可读性
type StorageObject struct {
	// Key 是对象的存储键（相对路径或对象名），不含任何静态前缀
	Key string
	// URL 是可访问地址：local 使用相对路径前缀 "/images/"，qiniu 使用完整域名
	URL string
	// Size 为对象大小（字节），如未知可填 0
	Size int64
	// Type 为 MIME 类型或业务自定义类型，如未知可留空
	Type string
	// Name 为原始文件名或业务命名，便于记录与排查
	Name string
}

//revive:enable:exported

// Driver 为统一的存储驱动接口，需实现并发安全。
// 调用约定：
// - Delete(ctx, key)：key 为存储键（不含前缀）。，资源不存在必须视为成功（幂等）
// - Upload(ctx, r, filename)：r 为流式数据源，filename 用作对象名/建议名；实现应尽量避免整文件缓冲。
// - Name()：返回驱动名（如 "local"、"qiniu"），用于日志与分支控制。
type Driver interface {
	Name() string
	Delete(ctx context.Context, key string) error
	Upload(ctx context.Context, r io.Reader, filename string) (StorageObject, error)
}
