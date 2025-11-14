package flag

import (
	"bufio"
	"fmt"
	"os"
	esModel "personal_blog/internal/model/elasticsearch"
	esutil "personal_blog/pkg/elasticSearch"
)

// Elasticsearch 创建 ES 索引
func Elasticsearch() error {

	// 检查索引是否已存在
	indexExists, err := esutil.IndexExists(esModel.ArticleIndex())
	if err != nil {
		return err
	}

	if indexExists {
		// 打印提示信息
		fmt.Println("The index already exists. Do you want to delete the data and recreate the index? (y/n)")

		// 读取用户输入
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "y":
			// 如果用户输入 y，删除索引
			fmt.Println("Proceeding to delete the data and recreate the index...")
			if err := esutil.IndexDelete(esModel.ArticleIndex()); err != nil {
				return err
			}
		case "n":
			// 如果用户输入 n，退出程序
			fmt.Println("Exiting the program.")
			os.Exit(0)
		default:
			// 如果用户输入无效，提示重新输入
			fmt.Println("Invalid input. Please enter 'y' to delete and recreate the index, or 'n' to exit.")
			return Elasticsearch() // 递归调用，重新输入
		}
	}

	// 创建索引
	return esutil.IndexCreate(esModel.ArticleIndex(), esModel.ArticleMapping())
}
