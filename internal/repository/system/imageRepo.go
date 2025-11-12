package system

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "personal_blog/internal/model/entity"
    "personal_blog/internal/repository/interfaces"
)

// ImageGormRepository 图片仓储GORM实现
type ImageGormRepository struct {
    db *gorm.DB
}

// NewImageRepository 创建图片仓储实例，返回接口类型
func NewImageRepository(db *gorm.DB) interfaces.ImageRepository {
    return &ImageGormRepository{db: db}
}

// Create 创建图片记录
func (r *ImageGormRepository) Create(ctx context.Context, img *entity.Image) error {
    return r.db.WithContext(ctx).Create(img).Error
}

// GetByID 根据ID查询图片
func (r *ImageGormRepository) GetByID(ctx context.Context, id uint) (*entity.Image, error) {
    var img entity.Image
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&img).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &img, nil
}

// ListByUser 根据用户ID分页查询图片
func (r *ImageGormRepository) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*entity.Image, int64, error) {
    var imgs []*entity.Image
    var total int64
    q := r.db.WithContext(ctx).Model(&entity.Image{}).Where("user_id = ?", userID)
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    offset := (page - 1) * pageSize
    if err := q.Offset(offset).Limit(pageSize).Order("id DESC").Find(&imgs).Error; err != nil {
        return nil, 0, err
    }
    return imgs, total, nil
}

// DeleteByID 根据ID删除图片（软删除）
func (r *ImageGormRepository) DeleteByID(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&entity.Image{}, id).Error
}