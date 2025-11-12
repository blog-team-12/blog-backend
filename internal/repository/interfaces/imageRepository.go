package interfaces

import (
    "context"
    "personal_blog/internal/model/entity"
)

// ImageRepository 图片仓储接口
type ImageRepository interface {
    // Create 创建图片记录
    Create(ctx context.Context, img *entity.Image) error
    // GetByID 根据ID查询图片
    GetByID(ctx context.Context, id uint) (*entity.Image, error)
    // ListByUser 根据用户ID分页查询图片
    ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*entity.Image, int64, error)
    // DeleteByID 根据ID删除图片（软删除）
    DeleteByID(ctx context.Context, id uint) error
}