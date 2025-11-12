package system

import (
    "fmt"
    "mime/multipart"
    "personal_blog/global"
    "personal_blog/internal/model/consts"
    "personal_blog/internal/model/entity"
    serviceSystem "personal_blog/internal/service/system"
    "personal_blog/pkg/imageUtils"
    "personal_blog/pkg/jwt"
    "personal_blog/pkg/response"
    "strconv"
    "strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageCtrl struct {
	imageService *serviceSystem.ImageService
	jwtService   *serviceSystem.JWTService
}

// Upload 处理图片上传
func (i *ImageCtrl) Upload(c *gin.Context) {
	uid := jwt.GetUserID(c)
	if uid == 0 {
		response.NewResponse[any, any](c).
			SetCode(global.StatusUnauthorized).
			Failed("未认证或凭证无效", nil)
		return
	}

	// 解析 driver 与 bind；默认 driver 使用当前配置，bind 默认 false
	driver := strings.ToLower(strings.TrimSpace(c.Query("driver")))
	bind := parseBind(c.Query("bind"))
	category := parseCategory(c)

	// 收集所有文件：支持 images 与 file 两个字段
	files := collectFiles(c)
	if len(files) == 0 {
		response.NewResponse[any, any](c).
			SetCode(global.StatusBadRequest).
			Failed("缺少文件: form字段为 'images' 或 'file'", nil)
		return
	}
	// 上传数量限制
	if max := global.Config.Static.MaxUploads; max > 0 && len(files) > max {
		response.NewResponse[any, any](c).
			SetCode(global.StatusBadRequest).
			Failed(fmt.Sprintf("单次最多上传 %d 个文件", max), nil)
		return
	}

	urls := make([]string, 0, len(files))
	ids := make([]uint, 0, len(files))

	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			global.Log.Error("打开文件失败", zap.Error(err))
			response.NewResponse[any, any](c).
				SetCode(global.StatusInternalServerError).
				Failed(fmt.Sprintf("打开文件失败: %v", err), nil)
			return
		}
		// 每个文件独立处理
		if bind {
			img, err := i.imageService.UploadImageWithDriver(
				c.Request.Context(), uid, fh.Filename, f, category, driver,
			)
			_ = f.Close()
			if err != nil {
				global.Log.Error("图片上传并绑定失败", zap.Error(err))
				response.NewResponse[any, any](c).
					SetCode(global.StatusInternalServerError).
					Failed(fmt.Sprintf("上传失败: %v", err), nil)
				return
			}
			urls = append(urls, img.URL)
			ids = append(ids, img.ID)
        } else {
            // 不入库：使用工具包上传，并规范化返回URL
            if strings.TrimSpace(driver) == "" {
                obj, err := imageUtils.UploadViaCurrent(c.Request.Context(), f, fh.Filename)
                _ = f.Close()
                if err != nil {
                    global.Log.Error("图片上传失败", zap.Error(err))
                    response.NewResponse[any, any](c).
                        SetCode(global.StatusInternalServerError).
                        Failed(fmt.Sprintf("上传失败: %v", err), nil)
                    return
                }
                tmp := &entity.Image{Key: obj.Key, URL: obj.URL}
                urls = append(urls, imageUtils.ImageURL(tmp))
            } else {
                obj, err := imageUtils.UploadWithDriver(c.Request.Context(), driver, &entity.Image{}, f, fh.Filename)
                _ = f.Close()
                if err != nil {
                    global.Log.Error("图片上传失败", zap.Error(err))
                    response.NewResponse[any, any](c).
                        SetCode(global.StatusInternalServerError).
                        Failed(fmt.Sprintf("上传失败: %v", err), nil)
                    return
                }
                tmp := &entity.Image{Key: obj.Key, URL: obj.URL}
                urls = append(urls, imageUtils.ImageURL(tmp))
            }
        }
    }

	// 响应数据格式按 bind 选择
	if bind {
		response.NewResponse[any, any](c).
			SetCode(global.StatusOK).
			Success("上传成功", gin.H{"urls": urls, "ids": ids})
	} else {
		response.NewResponse[any, any](c).
			SetCode(global.StatusOK).
			Success("上传成功", gin.H{"urls": urls})
	}
}

// List 列出当前用户的图片（分页）
func (i *ImageCtrl) List(c *gin.Context) {
	uid := jwt.GetUserID(c)
	if uid == 0 {
		response.NewResponse[any, any](c).
			SetCode(global.StatusUnauthorized).
			Failed("未认证或凭证无效", nil)
		return
	}

	page, pageSize := parsePagination(c)
	items, total, err := i.imageService.ListUserImages(c.Request.Context(), uid, page, pageSize)
	if err != nil {
		global.Log.Error("获取图片列表失败", zap.Error(err))
		response.NewResponse[any, any](c).
			SetCode(global.StatusInternalServerError).
			Failed(fmt.Sprintf("获取列表失败: %v", err), nil)
		return
	}

	// 精简返回结构
	list := make([]map[string]any, 0, len(items))
	for _, img := range items {
		list = append(list, map[string]any{
			"id":         img.ID,
			"name":       img.Name,
			"url":        img.URL,
			"key":        img.Key,
			"category":   img.Category,
			"storage":    img.Storage,
			"created_at": img.CreatedAt,
		})
	}

	response.NewResponse[any, any](c).
		SetCode(global.StatusOK).
		Success("获取成功", gin.H{
			"items":     list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
}

// Delete 删除当前用户的指定图片
func (i *ImageCtrl) Delete(c *gin.Context) {
	uid := jwt.GetUserID(c)
	if uid == 0 {
		response.NewResponse[any, any](c).
			SetCode(global.StatusUnauthorized).
			Failed("未认证或凭证无效", nil)
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.NewResponse[any, any](c).
			SetCode(global.StatusBadRequest).
			Failed("参数错误: id 应为数字", nil)
		return
	}

	if err := i.imageService.DeleteImage(c.Request.Context(), uid, uint(id64)); err != nil {
		global.Log.Error("删除图片失败", zap.Error(err))
		response.NewResponse[any, any](c).
			SetCode(global.StatusInternalServerError).
			Failed(fmt.Sprintf("删除失败: %v", err), nil)
		return
	}
	response.NewResponse[any, any](c).
		SetCode(global.StatusOK).
		Success("删除成功", nil)
}

func parseCategory(c *gin.Context) consts.Category {
	// 支持从 query 或 form 读取
	raw := c.DefaultPostForm("category", c.Query("category"))
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return consts.Null
	}
	if n, err := strconv.Atoi(raw); err == nil {
		return consts.Category(n)
	}
	cat := consts.ToCategory(raw)
	if cat < 0 {
		return consts.Null
	}
	return cat
}

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	v := c.Query("page_size")
	if v == "" {
		v = c.Query("pageSize")
	}
	if v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}
	return page, pageSize
}

// parseBind 解析 bind=1/true 为 true，默认 false
func parseBind(v string) bool {
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "1" || v == "true" || v == "t" || v == "yes" || v == "y" {
		return true
	}
	return false
}

// collectFiles 收集 images 与 file 字段的所有文件
func collectFiles(c *gin.Context) []*multipart.FileHeader {
	var out []*multipart.FileHeader
	if form, err := c.MultipartForm(); err == nil && form != nil {
		if arr := form.File["images"]; len(arr) > 0 {
			out = append(out, arr...)
		}
		if arr := form.File["file"]; len(arr) > 0 {
			out = append(out, arr...)
		}
	}
	// 兼容只有一个文件且未使用 MultipartForm 的情况
	if len(out) == 0 {
		if fh, err := c.FormFile("images"); err == nil && fh != nil {
			out = append(out, fh)
		}
	}
	if len(out) == 0 {
		if fh, err := c.FormFile("file"); err == nil && fh != nil {
			out = append(out, fh)
		}
	}
	return out
}

// Ensure entity.Image is referenced to avoid unused import in some builds
var _ = entity.Image{}
