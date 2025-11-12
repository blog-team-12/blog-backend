package system

import (
    "context"
    "errors"
    "io"
    "mime/multipart"
    "path/filepath"
    "personal_blog/global"
    "personal_blog/internal/model/consts"
    "personal_blog/internal/model/entity"
    "personal_blog/internal/repository"
    repoSystem "personal_blog/internal/repository/system"
    "personal_blog/pkg/imageUtils"
    "personal_blog/pkg/storage"
    "strings"
)

// ImageService 图片业务服务
type ImageService struct {
	repos repoSystem.Supplier
}

// NewImageService 创建图片服务
func NewImageService(group *repository.Group) *ImageService {
	return &ImageService{repos: group.SystemRepositorySupplier}
}

// UploadImage 处理上传：校验、调用存储驱动、落库
func (s *ImageService) UploadImage(
    ctx context.Context,
    userID uint,
    filename string,
    reader io.Reader,
    category consts.Category,
) (*entity.Image, error) {
	// 校验扩展名
	if !isAllowedExt(filename) {
		return nil, errors.New("file type not allowed")
	}
	// 校验大小（MB -> bytes），若配置为0则不限制
	if global.Config.Static.MaxSize > 0 {
		limit := int64(global.Config.Static.MaxSize) * 1024 * 1024
		reader = io.LimitReader(reader, limit)
	}

    // 使用工具包上传并绑定字段到模型（不入库）
    img := &entity.Image{ // 先设置基础字段
        Name:     filename,
        Category: category,
        UserID:   &userID,
    }
    if _, err := imageUtils.UploadAndBind(ctx, img, reader, filename); err != nil {
        return nil, err
    }
    if err := s.repos.GetImageRepository().Create(ctx, img); err != nil {
        return nil, err
    }
    return img, nil
}

// UploadImageWithDriver 指定驱动名上传并绑定数据库记录
func (s *ImageService) UploadImageWithDriver(
    ctx context.Context,
    userID uint,
    filename string,
    reader io.Reader,
    category consts.Category,
    driverName string,
) (*entity.Image, error) {
	// 拓展名限制
	if !isAllowedExt(filename) {
		return nil, errors.New("file type not allowed")
	}
	// 限制流量读取
	if global.Config.Static.MaxSize > 0 {
		limit := int64(global.Config.Static.MaxSize) * 1024 * 1024
		reader = io.LimitReader(reader, limit)
	}
    // 使用工具包按指定驱动上传并绑定字段到模型（不入库）
    img := &entity.Image{ // 先设置基础字段
        Name:     filename,
        Category: category,
        UserID:   &userID,
    }
    if _, err := imageUtils.UploadWithDriver(ctx, strings.ToLower(strings.TrimSpace(driverName)), img, reader, filename); err != nil {
        return nil, err
    }
    if err := s.repos.GetImageRepository().Create(ctx, img); err != nil {
        return nil, err
    }
    return img, nil
}

// ListUserImages 分页列出用户图片
func (s *ImageService) ListUserImages(
	ctx context.Context,
	userID uint,
	page, pageSize int,
) ([]*entity.Image, int64, error) {
	return s.repos.GetImageRepository().ListByUser(ctx, userID, page, pageSize)
}

// DeleteImage 删除图片：先删存储，再删记录
func (s *ImageService) DeleteImage(
	ctx context.Context,
	userID uint,
	id uint,
) error {
	img, err := s.repos.GetImageRepository().GetByID(ctx, id)
	if err != nil {
		return err
	}
	if img == nil {
		return errors.New("无图片资源")
	}
    if img.UserID == nil || *img.UserID != userID {
        return errors.New("无权限删除 其他用户的图片")
    }

	drv := storage.FromName(strings.ToLower(storage.Current().Name()))
	if drv == nil {
		return errors.New("驱动未初始化")
	}
	if err := drv.Delete(ctx, img.Key); err != nil {
		return err
	}
	return s.repos.GetImageRepository().DeleteByID(ctx, id)
}

// isAllowedExt 检查文件扩展名是否允许
func isAllowedExt(
    filename string,
) bool {
    ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
    if ext == "" {
        return false
    }
    allowed := global.Config.Static.AllowedTypes
    if len(allowed) == 0 {
        // 未配置则默认允许常见类型
        allowed = []string{"png", "jpg", "jpeg", "gif", "webp"}
		global.Log.Warn("未配置允许的文件类型，已使用默认配置")
    } else {
        // 兼容配置中带点的写法，例如 ".png"，统一去掉前导点并小写化
        for i, a := range allowed {
            allowed[i] = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(a), "."))
        }
    }
    for _, a := range allowed {
        if a == ext {
            return true
        }
    }
    return false
}

// FromMultipart 便捷方法：从 *multipart.FileHeader 上传
func (s *ImageService) FromMultipart(
	ctx context.Context,
	userID uint,
	fh *multipart.FileHeader,
	category consts.Category,
) (*entity.Image, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return s.UploadImage(ctx, userID, fh.Filename, f, category)
}
