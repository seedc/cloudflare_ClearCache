package zoneid

import (
	"cloudflare_ClearCache/cache"
	"cloudflare_ClearCache/settings"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"strings"
)

func CloudFrontid(domainToFind string) {
	// 硬编码 AWS 凭证
	awsAccessKeyID := settings.Conf.AwsAccessKeyID         // 替换为你的 Access Key ID
	awsSecretAccessKey := settings.Conf.AwsSecretAccessKey // 替换为你的 Secret Access Key
	awsRegion := settings.Conf.AwsRegion                   // 替换为你的区域

	// 创建一个新的 AWS 会话，使用硬编码的凭证
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	}))

	// 创建 CloudFront 服务客户端
	svc := cloudfront.New(sess)

	// 列出所有分配
	input := &cloudfront.ListDistributionsInput{}
	result, err := svc.ListDistributions(input)
	if err != nil {
		fmt.Println("Error listing distributions:", err)
		return
	}

	// 查找匹配的域名
	for _, dist := range result.DistributionList.Items {
		for _, alias := range dist.Aliases.Items {
			if strings.Contains(*alias, domainToFind) {
				fmt.Printf("Found distribution ID for domain %s: %s\n", domainToFind, *dist.Id)
				cache.AwscCche(awsAccessKeyID, awsSecretAccessKey, awsRegion, *dist.Id)
				return
			}
		}
	}

	fmt.Println("No distribution found for domain", domainToFind)
}
