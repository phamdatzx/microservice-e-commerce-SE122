package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Client *s3.S3
var S3BucketName string

func InitS3() {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	endpoint := os.Getenv("AWS_S3_ENDPOINT") // Optional: for MinIO or LocalStack

	fmt.Printf("🔍 S3 Config - Region: %s, Bucket: %s, Endpoint: %s\n", region, bucketName, endpoint)

	if region == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		fmt.Println("⚠️  Missing AWS S3 configuration (AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_S3_BUCKET_NAME)")
		return
	}

	S3BucketName = bucketName

	// Build AWS config with timeout
	awsConfig := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second, // 5 second timeout to prevent hanging
		},
	}

	// Add endpoint if specified (for MinIO, LocalStack, etc.)
	if endpoint != "" {
		awsConfig.Endpoint = aws.String(endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true)       // Required for MinIO and some S3-compatible services
		awsConfig.DisableSSL = aws.Bool(true)             // Disable SSL for local development
		fmt.Printf("🔧 Using custom S3 endpoint: %s (SSL disabled)\n", endpoint)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Printf("❌ Failed to create AWS session: %v\n", err)
		return
	}

	S3Client = s3.New(sess)
	fmt.Println("✅ S3 client initialized")

	// Test connection by listing buckets
	fmt.Println("🧪 Testing S3 connection...")
	buckets, err := S3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("❌ Failed to list S3 buckets: %v\n", err)
		fmt.Println("💡 Common issues:")
		fmt.Println("   - Check if MinIO/S3 service is running")
		fmt.Println("   - Verify AWS_S3_ENDPOINT URL is correct")
		fmt.Println("   - Ensure credentials are valid")
		return
	}
	
	fmt.Printf("✅ Successfully connected to S3! Found %d bucket(s)\n", len(buckets.Buckets))
	for _, bucket := range buckets.Buckets {
		fmt.Printf("   📦 %s\n", *bucket.Name)
	}
}
