package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// CloudStorageClient wraps GCP Cloud Storage operations
type CloudStorageClient struct {
	client     *storage.Client
	bucketName string
	ctx        context.Context
}

// NewCloudStorageClient creates a new Cloud Storage client
func NewCloudStorageClient(ctx context.Context, bucketName string) (*CloudStorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &CloudStorageClient{
		client:     client,
		bucketName: bucketName,
		ctx:        ctx,
	}, nil
}

// Close closes the client connection
func (c *CloudStorageClient) Close() error {
	return c.client.Close()
}

// CreateBucket creates a new bucket
func (c *CloudStorageClient) CreateBucket(projectID, location string) error {
	bucket := c.client.Bucket(c.bucketName)
	
	if err := bucket.Create(c.ctx, projectID, &storage.BucketAttrs{
		Location: location,
	}); err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}
	
	log.Printf("✅ Bucket %s created successfully", c.bucketName)
	return nil
}

// UploadFile uploads a file to Cloud Storage
func (c *CloudStorageClient) UploadFile(objectName string, content []byte) error {
	obj := c.client.Bucket(c.bucketName).Object(objectName)
	writer := obj.NewWriter(c.ctx)
	
	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("failed to write to object: %w", err)
	}
	
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	
	log.Printf("✅ File %s uploaded successfully", objectName)
	return nil
}

// UploadFileFromPath uploads a file from local path
func (c *CloudStorageClient) UploadFileFromPath(objectName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	obj := c.client.Bucket(c.bucketName).Object(objectName)
	writer := obj.NewWriter(c.ctx)
	
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	
	log.Printf("✅ File %s uploaded from %s", objectName, filePath)
	return nil
}

// DownloadFile downloads a file from Cloud Storage
func (c *CloudStorageClient) DownloadFile(objectName string) ([]byte, error) {
	obj := c.client.Bucket(c.bucketName).Object(objectName)
	reader, err := obj.NewReader(c.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()
	
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	
	log.Printf("✅ File %s downloaded (%d bytes)", objectName, len(data))
	return data, nil
}

// DownloadFileToPath downloads a file to local path
func (c *CloudStorageClient) DownloadFileToPath(objectName, destPath string) error {
	data, err := c.DownloadFile(objectName)
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
func (c *CloudStorageClient) ListFiles(prefix string) ([]string, error) {
	var files []string
	
	query := &storage.Query{Prefix: prefix}
	it := c.client.Bucket(c.bucketName).Objects(c.ctx, query)
	
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate objects: %w", err)
		}
		files = append(files, attrs.Name)
	}
	
	log.Printf("✅ Found %d files with prefix '%s'", len(files), prefix)
	return files, nil
}

// DeleteFile deletes a file from Cloud Storage
func (c *CloudStorageClient) DeleteFile(objectName string) error {
	obj := c.client.Bucket(c.bucketName).Object(objectName)
	
	if err := obj.Delete(c.ctx); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	
	log.Printf("✅ File %s deleted successfully", objectName)
	return nil
}

// GetFileMetadata retrieves file metadata
func (c *CloudStorageClient) GetFileMetadata(objectName string) (*storage.ObjectAttrs, error) {
	obj := c.client.Bucket(c.bucketName).Object(objectName)
	attrs, err := obj.Attrs(c.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %w", err)
	}
	
	log.Printf("✅ Metadata for %s:", objectName)
	log.Printf("   Size: %d bytes", attrs.Size)
	log.Printf("   Content-Type: %s", attrs.ContentType)
	log.Printf("   Created: %s", attrs.Created)
	log.Printf("   Updated: %s", attrs.Updated)
	
	return attrs, nil
}

// MakePublic makes a file publicly accessible
func (c *CloudStorageClient) MakePublic(objectName string) (string, error) {
	obj := c.client.Bucket(c.bucketName).Object(objectName)
	
	if err := obj.ACL().Set(c.ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set ACL: %w", err)
	}
	
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectName)
	log.Printf("✅ File is now public: %s", publicURL)
	
	return publicURL, nil
}

// GenerateSignedURL generates a signed URL for temporary access
func (c *CloudStorageClient) GenerateSignedURL(objectName string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}
	
	url, err := c.client.Bucket(c.bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	
	log.Printf("✅ Signed URL generated (expires in %v)", expiration)
	return url, nil
}

func main() {
	ctx := context.Background()
	
	bucketName := os.Getenv("GCP_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("GCP_BUCKET_NAME environment variable is required")
	}
	
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable is required")
	}
	
	// Create client
	client, err := NewCloudStorageClient(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	log.Println("🚀 Cloud Storage Demo")
	log.Println("=" + string(make([]byte, 50)))
	
	// Example 1: Upload file
	log.Println("\n📤 Example 1: Upload file")
	content := []byte("Hello from Cloud Storage!")
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
	
	// Example 5: Generate signed URL
	log.Println("\n🔗 Example 5: Generate signed URL")
	if url, err := client.GenerateSignedURL("demo/hello.txt", 15*time.Minute); err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("URL: %s", url)
	}
	
	log.Println("\n✅ All examples completed!")
}

