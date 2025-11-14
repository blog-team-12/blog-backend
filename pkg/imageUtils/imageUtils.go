package imageUtils

import (
	"context"
	"gorm.io/gorm"
	"personal_blog/internal/model/consts"
	"personal_blog/internal/model/entity"
	"regexp"
	"strings"
)

// InitImagesCategory 初始化图片类别
func InitImagesCategory(ctx context.Context, tx *gorm.DB, urls []string) error {
	return tx.WithContext(ctx).
		Model(&entity.Image{}).
		Where("url IN ?", urls).
		Update("category", consts.Category(0)).
		Error
}

// ChangeImagesCategory 修改图片类别
func ChangeImagesCategory(ctx context.Context, tx *gorm.DB, urls []string, category consts.Category) error {
	return tx.WithContext(ctx).
		Model(&entity.Image{}).
		Where("url IN ?", urls).
		Update("category", category).
		Error
}

// FindIllustrations 获取插图url
// 格式: ![alt text](image_url#pic_center)
// 如：https://i-blog.csdnimg.cn/direct/8a17a2ac6127438590f4bd8c67a43181.png
func FindIllustrations(text string) ([]string, error) {
	// 1、定义正则表达式，匹配 Markdown 图片语法
	// 格式: ![alt text](image_url)
	regex := `!\[([^\]]*)\]\(([^)]+)\)`

	// 编译正则表达式
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	// 2、查找所有符合正则表达式的匹配项
	matches := re.FindAllStringSubmatch(text, -1)

	// 3、存储匹配到的所有图片链接
	var illustrations []string
	// 3.a 循环剔除
	for _, match := range matches {
		if len(match) > 2 {
			url := match[2]
			// 3.b 移除 URL 片段（如 #pic_center）
			if idx := strings.Index(url, "#"); idx != -1 {
				url = url[:idx]
			}
			// 3.c 获取图片链接
			illustrations = append(illustrations, url)
		}
	}

	return illustrations, nil
}
