package system

import (
	"context"
	"errors"
	"personal_blog/internal/model/entity"
	"personal_blog/internal/repository/interfaces"

	"gorm.io/gorm"
)

// ArticleGormRepository 文章仓储实现
type ArticleGormRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储实例
func NewArticleRepository(db *gorm.DB) interfaces.ArticleRepository {
	return &ArticleGormRepository{db: db}
}

// IncOrCreateCategory 在分类不存在时创建并计数为1，存在时将计数+1
func (r *ArticleGormRepository) IncOrCreateCategory(
	ctx context.Context,
	tx *gorm.DB,
	category string) error {
	var c entity.ArticleCategory
	if errors.Is(
		tx.WithContext(ctx).Where("category = ?", category).First(&c).Error,
		gorm.ErrRecordNotFound,
	) {
		return tx.Create(&entity.ArticleCategory{Category: category, Number: 1}).Error
	}
	return tx.Model(&c).Update("number", gorm.Expr("number + ?", 1)).Error
}

// DecOrDeleteCategory 将分类计数-1；若减少前计数为1则删除该分类
func (r *ArticleGormRepository) DecOrDeleteCategory(
	ctx context.Context,
	tx *gorm.DB,
	category string,
) error {
	var c entity.ArticleCategory
	if err := tx.WithContext(ctx).Where("category = ?", category).
		First(&c).
		Update("number", gorm.Expr("number - ?", 1)).
		Error; err != nil {
		return err
	}
	if c.Number == 1 {
		if err := tx.WithContext(ctx).Delete(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

// AddOrIncTag 在标签不存在时创建并计数为1，存在时将计数+1
func (r *ArticleGormRepository) AddOrIncTag(
    ctx context.Context,
    tx *gorm.DB,
    tags []string,
) error {
    for _, tag := range tags {
        if tag == "" {
            continue
        }
        var t entity.ArticleTag
        if errors.Is(
            tx.WithContext(ctx).Where("tag = ?", tag).First(&t).Error,
            gorm.ErrRecordNotFound,
        ) {
            if err := tx.WithContext(ctx).Create(&entity.ArticleTag{Tag: tag, Number: 1}).Error; err != nil {
                return err
            }
            continue
        }
        if err := tx.WithContext(ctx).
            Model(&entity.ArticleTag{}).
            Where("tag = ?", tag).
            Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
            return err
        }
    }
    return nil
}

// DecOrDeleteTag 将标签计数-1；若减少前计数为1则删除该标签
func (r *ArticleGormRepository) DecOrDeleteTag(
    ctx context.Context,
    tx *gorm.DB,
    tags []string,
) error {
    for _, tag := range tags {
        if tag == "" {
            continue
        }
        var t entity.ArticleTag
        if err := tx.WithContext(ctx).Where("tag = ?", tag).First(&t).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                continue
            }
            return err
        }
        if t.Number <= 1 {
            if err := tx.WithContext(ctx).Delete(&t).Error; err != nil {
                return err
            }
            continue
        }
        if err := tx.WithContext(ctx).
            Model(&entity.ArticleTag{}).
            Where("tag = ?", tag).
            Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
            return err
        }
    }
    return nil
}

// Transaction 事物统一处理，用以保证原子性
func (r *ArticleGormRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
