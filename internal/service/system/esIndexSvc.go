package system

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"personal_blog/global"
)

// EsService 提供了对 Elasticsearch 索引的操作方法
type EsService struct{}

// NewEsService 返回一个实例
func NewEsService() *EsService {
	return &EsService{}
}

/*
## 参数说明
	- indexName string ：要创建的索引名称（比如 "blog_articles"）
	- mapping *types.TypeMapping ：索引的字段映射配置，定义了文档中各字段的数据类型和属性
## 工作流程
	1. 调用 ES 客户端 ：使用全局的 global.ESClient 连接 Elasticsearch
	2. 创建索引 ：通过 Indices.Create(indexName) 指定要创建的索引名
	3. 设置映射 ：通过 .Mappings(mapping) 为索引配置字段结构
	4. 执行操作 ： .Do(context.TODO()) 实际执行创建命令
	5. 返回结果 ：如果创建成功返回 nil，失败则返回错误信息
*/

// IndexCreate 创建一个新的 Elasticsearch 索引，带有指定的映射
func (esService *EsService) IndexCreate(indexName string, mapping *types.TypeMapping) error {
	_, err := global.ESClient.Indices.Create(indexName).Mappings(mapping).Do(context.TODO())
	return err
}

// IndexDelete 删除指定的 Elasticsearch 索引
func (esService *EsService) IndexDelete(indexName string) error {
	_, err := global.ESClient.Indices.Delete(indexName).Do(context.TODO())
	return err
}

// IndexExists 检查指定的 Elasticsearch 索引是否存在
func (esService *EsService) IndexExists(indexName string) (bool, error) {
	return global.ESClient.Indices.Exists(indexName).Do(context.TODO())
}
