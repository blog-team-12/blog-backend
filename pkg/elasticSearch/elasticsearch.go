package elasticSearch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"personal_blog/global"
	"personal_blog/internal/model/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
)

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
func IndexCreate(indexName string, mapping *types.TypeMapping) error {
	_, err := global.ESClient.Indices.Create(indexName).Mappings(mapping).Do(context.TODO())
	return err
}

// IndexDelete 删除指定的 Elasticsearch 索引
func IndexDelete(indexName string) error {
	_, err := global.ESClient.Indices.Delete(indexName).Do(context.TODO())
	return err
}

// IndexExists 检查指定的 Elasticsearch 索引是否存在
func IndexExists(indexName string) (bool, error) {
	return global.ESClient.Indices.Exists(indexName).Do(context.TODO())
}

// Get 用于通过ID从 Elasticsearch 获取文章
func Get(
	ctx context.Context,
	id string,
) (elasticsearch.Article, error) {
	// 1、从Elasticsearch获取文章
	// 1.a 获取文章
	var a elasticsearch.Article
	res, err := global.ESClient.
		Get(elasticsearch.ArticleIndex(), id).
		Do(ctx)
	if err != nil {
		return elasticsearch.Article{}, err
	}
	// 1.b、如果找不到该文档，则返回错误
	if !res.Found {
		return elasticsearch.Article{}, errors.New("document not found")
	}
	// 2、将返回的源数据反序列化为 Article 对象
	err = json.Unmarshal(res.Source_, &a)
	return a, err
}

// Update 用于更新文章数据
func Update(
	ctx context.Context,
	articleID string,
	v any,
) error {
	// 1、将待更新的值转换为 JSON
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// 2、执行更新请求，并设置刷新操作为 true
	_, err = global.ESClient.
		Update(elasticsearch.ArticleIndex(), articleID).
		Request(&update.Request{Doc: bytes}).
		Refresh(refresh.True).
		Do(ctx)

	return err
}

// Delete 用于删除 Elasticsearch 中的文章
func Delete(
	ctx context.Context,
	ids []string,
) error {
	var request bulk.Request
	// 遍历文章ID，构建批量删除请求
	for _, id := range ids {
		request = append(request, types.OperationContainer{Delete: &types.DeleteOperation{Id_: &id}})
	}
	// 执行批量删除请求，并设置刷新操作为 true
	_, err := global.ESClient.Bulk(). // 创建批量操作，用于一次性执行多个索引、更新、删除操作
						Request(&request).
						Index(elasticsearch.ArticleIndex()).
						Refresh(refresh.True).Do(ctx)
	return err
}
