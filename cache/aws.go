package cache

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"time"
)

func AwscCche(awsAccessKeyID string, awsSecretAccessKey string, awsRegion string, id string) {

	// 创建一个新的 AWS 会话，使用硬编码的凭证
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	}))

	// 创建 CloudFront 服务客户端
	svc := cloudfront.New(sess)

	// 设置 CloudFront 分配的 ID 和要清除的路径
	distributionID := id   // 替换为你的 CloudFront 分配 ID
	paths := []string{"/"} // 替换为你想要清除的路径

	// 构建清除请求
	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(fmt.Sprintf("invalidation-%d", time.Now().Unix())),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(int64(len(paths))),
				Items:    aws.StringSlice(paths),
			},
		},
	}

	// 发送清除请求
	result, err := svc.CreateInvalidation(input)
	if err != nil {
		fmt.Println("Error creating invalidation:", err)
		return
	}

	fmt.Printf("Invalidation ID: %s\n", *result.Invalidation.Id)
}
