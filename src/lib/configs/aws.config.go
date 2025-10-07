package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func InitS3() {
	c := EnvModule()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.S3.AccessKeyID,
			c.S3.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config: %v", err)
	}

	S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(c.S3.URL)
	})

	log.Println("✅ Connected to S3 (LocalStack or AWS)")
}
