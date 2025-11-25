package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Client wraps AWS S3 operations
type S3Client struct {
	client     *s3.Client
	bucketName string
	ctx        context.Context
}

// NewS3Client creates a new S3 client
func NewS3Client(ctx context.Context, bucketName string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{
		client:     client,
		bucketName: bucketName,
		ctx:        ctx,
	}, nil
}

// CreateBucket creates a new S3 bucket
func (c *S3Client) CreateBucket(region string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(c.bucketName),
	}

	// For regions other than us-east-1, specify location constraint
	if region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		}
	}

	_, err := c.client.CreateBucket(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	log.Printf("✅ Bucket %s created successfully", c.bucketName)
	return nil
}

// UploadFile uploads a file to S3
func (c *S3Client) UploadFile(key string, content []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	}

	_, err := c.client.PutObject(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("✅ File %s uploaded successfully", key)
	return nil
}

// UploadFileFromPath uploads a file from local path
func (c *S3Client) UploadFileFromPath(key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		Body:   file,
	}

	_, err = c.client.PutObject(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("✅ File %s uploaded from %s", key, filePath)
	return nil
}

// DownloadFile downloads a file from S3
func (c *S3Client) DownloadFile(key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	result, err := c.client.GetObject(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	log.Printf("✅ File %s downloaded (%d bytes)", key, len(data))
	return data, nil
}

// DownloadFileToPath downloads a file to local path
func (c *S3Client) DownloadFileToPath(key, destPath string) error {
	data, err := c.DownloadFile(key)
	if err != nil {
		return err
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("✅ File saved to %s", destPath)
	return nil
}

// ListFiles lists all files in the bucket
func (c *S3Client) ListFiles(prefix string) ([]string, error) {
	var files []string

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(c.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(c.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			files = append(files, *obj.Key)
		}
	}

	log.Printf("✅ Found %d files with prefix '%s'", len(files), prefix)
	return files, nil
}

// DeleteFile deletes a file from S3
func (c *S3Client) DeleteFile(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	_, err := c.client.DeleteObject(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Printf("✅ File %s deleted successfully", key)
	return nil
}

// GetFileMetadata retrieves file metadata
func (c *S3Client) GetFileMetadata(key string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	result, err := c.client.HeadObject(c.ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	log.Printf("✅ Metadata for %s:", key)
	log.Printf("   Size: %d bytes", *result.ContentLength)
	log.Printf("   Content-Type: %s", *result.ContentType)
	log.Printf("   Last Modified: %s", *result.LastModified)
	if result.ETag != nil {
		log.Printf("   ETag: %s", *result.ETag)
	}

	return result, nil
}

// MakePublic makes a file publicly accessible
func (c *S3Client) MakePublic(key string) (string, error) {
	input := &s3.PutObjectAclInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPublicRead,
	}

	_, err := c.client.PutObjectAcl(c.ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to set ACL: %w", err)
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.bucketName, region, key)
	log.Printf("✅ File is now public: %s", publicURL)

	return publicURL, nil
}

// GeneratePresignedURL generates a presigned URL for temporary access
func (c *S3Client) GeneratePresignedURL(key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)

	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	result, err := presignClient.PresignGetObject(c.ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	log.Printf("✅ Presigned URL generated (expires in %v)", expiration)
	return result.URL, nil
}

// CopyFile copies a file within S3
func (c *S3Client) CopyFile(sourceKey, destKey string) error {
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(c.bucketName),
		CopySource: aws.String(fmt.Sprintf("%s/%s", c.bucketName, sourceKey)),
		Key:        aws.String(destKey),
	}

	_, err := c.client.CopyObject(c.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	log.Printf("✅ File copied from %s to %s", sourceKey, destKey)
	return nil
}

func main() {
	ctx := context.Background()

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("AWS_BUCKET_NAME environment variable is required")
	}

	// Create client
	client, err := NewS3Client(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("🚀 AWS S3 Demo")
	log.Println("=" + string(make([]byte, 50)))

	// Example 1: Upload file
	log.Println("\n📤 Example 1: Upload file")
	content := []byte("Hello from AWS S3!")
	if err := client.UploadFile("demo/hello.txt", content); err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 2: Download file
	log.Println("\n📥 Example 2: Download file")
	data, err := client.DownloadFile("demo/hello.txt")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Content: %s", string(data))
	}

	// Example 3: List files
	log.Println("\n📋 Example 3: List files")
	files, err := client.ListFiles("demo/")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, file := range files {
			log.Printf("  - %s", file)
		}
	}

	// Example 4: Get metadata
	log.Println("\n📊 Example 4: Get file metadata")
	if _, err := client.GetFileMetadata("demo/hello.txt"); err != nil {
		log.Printf("Error: %v", err)
	}

	// Example 5: Generate presigned URL
	log.Println("\n🔗 Example 5: Generate presigned URL")
	if url, err := client.GeneratePresignedURL("demo/hello.txt", 15*time.Minute); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("URL: %s", url)
	}

	// Example 6: Copy file
	log.Println("\n📋 Example 6: Copy file")
	if err := client.CopyFile("demo/hello.txt", "demo/hello-copy.txt"); err != nil {
		log.Printf("Error: %v", err)
	}

	log.Println("\n✅ All examples completed!")
}

