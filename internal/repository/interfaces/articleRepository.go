package interfaces

import (
	"context"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	// IncOrCreateCategory 在分类不存在时创建并计数为1，存在时将计数+1
	IncOrCreateCategory(ctx context.Context, tx *gorm.DB, category string) error
	// DecOrDeleteCategory 将分类计数-1；若减少前计数为1则删除该分类
	DecOrDeleteCategory(ctx context.Context, tx *gorm.DB, category string) error
	// AddOrIncTag 在标签不存在时创建并计数为1，存在时将计数+1
	AddOrIncTag(ctx context.Context, tx *gorm.DB, tags []string) error
	// DecOrDeleteTag 将标签计数-1；若减少前计数为1则删除该标签
	DecOrDeleteTag(ctx context.Context, tx *gorm.DB, tags []string) error
}
