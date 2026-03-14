package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client wraps the S3 service
type S3Client struct {
	client *s3.Client
}

// NewS3Client creates a new S3 client
func NewS3Client(ctx context.Context) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3Client{
		client: s3.NewFromConfig(cfg),
	}, nil
}

// ProcessImage processes an uploaded image
func (sc *S3Client) ProcessImage(ctx context.Context, bucket, key string) error {
	// Get object from S3
	result, err := sc.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	// Read object content
	content, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("failed to read object: %w", err)
	}

	log.Printf("Processed %s from bucket %s (size: %d bytes)", key, bucket, len(content))

	// Example processing: create thumbnail
	// In real application, use image processing library
	thumbnailKey := fmt.Sprintf("thumbnails/%s", key)

	_, err = sc.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(thumbnailKey),
		Body:   strings.NewReader("thumbnail content"),
	})

	if err != nil {
		return fmt.Errorf("failed to put thumbnail: %w", err)
	}

	log.Printf("Created thumbnail: %s", thumbnailKey)

	return nil
}

// Example 1: Simple S3 event handler
func s3EventHandler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3Entity := record.S3
		log.Printf("Bucket: %s, Key: %s", s3Entity.Bucket.Name, s3Entity.Object.Key)

		// Process the file
		if err := processFile(ctx, s3Entity.Bucket.Name, s3Entity.Object.Key); err != nil {
			log.Printf("Error processing file: %v", err)
			return err
		}
	}

	return nil
}

// Example 2: S3 event handler with S3 client
func s3ProcessHandler(ctx context.Context, s3Event events.S3Event) error {
	s3Client, err := NewS3Client(ctx)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	for _, record := range s3Event.Records {
		s3Entity := record.S3
		log.Printf("Processing: %s/%s", s3Entity.Bucket.Name, s3Entity.Object.Key)

		// Process the image
		if err := s3Client.ProcessImage(ctx, s3Entity.Bucket.Name, s3Entity.Object.Key); err != nil {
			log.Printf("Error processing image: %v", err)
			return err
		}
	}

	return nil
}

// Example 3: Filter specific file types
func imageFilterHandler(ctx context.Context, s3Event events.S3Event) error {
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	for _, record := range s3Event.Records {
		s3Entity := record.S3
		key := s3Entity.Object.Key

		// Check if file is an image
		ext := strings.ToLower(key[strings.LastIndex(key, "."):])
		if !imageExtensions[ext] {
			log.Printf("Skipping non-image file: %s", key)
			continue
		}

		log.Printf("Processing image: %s", key)

		if err := processFile(ctx, s3Entity.Bucket.Name, key); err != nil {
			log.Printf("Error processing image: %v", err)
			return err
		}
	}

	return nil
}

// Example 4: Process multiple files in parallel
func parallelProcessHandler(ctx context.Context, s3Event events.S3Event) error {
	s3Client, err := NewS3Client(ctx)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	// Process files in parallel
	errChan := make(chan error, len(s3Event.Records))

	for _, record := range s3Event.Records {
		go func(r events.S3EntityRecord) {
			s3Entity := r.S3
			errChan <- s3Client.ProcessImage(ctx, s3Entity.Bucket.Name, s3Entity.Object.Key)
		}(record)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(s3Event.Records); i++ {
		if err := <-errChan; err != nil {
			log.Printf("Error in parallel processing: %v", err)
		}
	}

	return nil
}

// Example 5: Handle S3 event with metadata
func metadataHandler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3Entity := record.S3
		log.Printf("Event Name: %s", record.EventName)
		log.Printf("Event Source: %s", record.EventSource)
		log.Printf("Event Time: %s", record.EventTime)
		log.Printf("Bucket: %s", s3Entity.Bucket.Name)
		log.Printf("Key: %s", s3Entity.Object.Key)
		log.Printf("Size: %d", s3Entity.Object.Size)
		log.Printf("ETag: %s", s3Entity.Object.ETag)
		log.Printf("Version ID: %s", aws.ToString(s3Entity.Object.VersionID))

		// Process file
		if err := processFile(ctx, s3Entity.Bucket.Name, s3Entity.Object.Key); err != nil {
			log.Printf("Error processing file: %v", err)
			return err
		}
	}

	return nil
}

// processFile is a helper function to process a file
func processFile(ctx context.Context, bucket, key string) error {
	log.Printf("Processing file: s3://%s/%s", bucket, key)

	// Add your processing logic here
	// Examples:
	// - Resize images
	// - Extract text from documents
	// - Validate file format
	// - Generate thumbnails
	// - Update database

	return nil
}

// Example 6: S3 to DynamoDB integration
func s3ToDynamoDBHandler(ctx context.Context, s3Event events.S3Event) error {
	// This would integrate with DynamoDB
	// See dynamodb.go for full implementation

	for _, record := range s3Event.Records {
		s3Entity := record.S3

		// Create record in DynamoDB
		log.Printf("Creating DynamoDB record for: %s", s3Entity.Object.Key)

		// In real application:
		// db := NewDynamoDBClient(ctx)
		// db.CreateFileRecord(ctx, FileRecord{
		//     Bucket: s3Entity.Bucket.Name,
		//     Key: s3Entity.Object.Key,
		//     Size: s3Entity.Object.Size,
		//     UploadedAt: time.Now(),
		// })
	}

	return nil
}

// Example 7: Error handling with retry logic
func retryHandler(ctx context.Context, s3Event events.S3Event) error {
	maxRetries := 3

	for _, record := range s3Event.Records {
		s3Entity := record.S3
		var err error

		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = processFile(ctx, s3Entity.Bucket.Name, s3Entity.Object.Key)
			if err == nil {
				log.Printf("Successfully processed %s on attempt %d", s3Entity.Object.Key, attempt)
				break
			}

			log.Printf("Attempt %d failed for %s: %v", attempt, s3Entity.Object.Key, err)

			if attempt == maxRetries {
				// Return error to trigger Lambda retry
				return fmt.Errorf("failed to process %s after %d attempts: %w", s3Entity.Object.Key, maxRetries, err)
			}
		}
	}

	return nil
}

func main() {
	// Choose which handler to use
	lambda.Start(imageFilterHandler)
}

/*
S3 Event Testing:

1. Create S3 bucket:
   aws s3 mb s3://my-lambda-bucket

2. Configure S3 event notification:
   aws s3api put-bucket-notification-configuration \
     --bucket my-lambda-bucket \
     --notification-configuration '{
       "LambdaFunctionConfigurations": [{
         "Id": "UploadEvents",
         "LambdaFunctionArn": "arn:aws:lambda:region:account:function:S3Handler",
         "Events": ["s3:ObjectCreated:*"],
         "Filter": {
           "Key": {
             "FilterRules": [
               {"Name": "suffix", "Value": "jpg"}
             ]
           }
         }
       }]
     }'

3. Add Lambda permission for S3 invocation:
   aws lambda add-permission \
     --function-name S3Handler \
     --principal s3.amazonaws.com \
     --statement-id s3invoke \
     --action "lambda:InvokeFunction" \
     --source-arn arn:aws:s3:::my-lambda-bucket \
     --source-account account-id

4. Test by uploading file:
   aws s3 cp test.jpg s3://my-lambda-bucket/

5. View CloudWatch logs:
   aws logs tail /aws/lambda/S3Handler --follow

Local Testing with SAM:

1. Create test event (events/s3_upload.json):
   {
     "Records": [{
       "eventVersion": "2.1",
       "eventSource": "aws:s3",
       "eventName": "ObjectCreated:Put",
       "eventTime": "2024-01-01T00:00:00.000Z",
       "s3": {
         "bucket": {
           "name": "my-lambda-bucket"
         },
         "object": {
           "key": "test.jpg",
           "size": 1024,
           "eTag": "abc123"
         }
       }
     }]
   }

2. Invoke locally:
   sam local invoke S3Handler --event events/s3_upload.json

IAM Permissions Required:

- s3:GetObject
- s3:PutObject (for creating thumbnails/processed files)
- logs:CreateLogGroup
- logs:CreateLogStream
- logs:PutLogEvents

S3 Event Types:

- s3:ObjectCreated:* - Any object creation
- s3:ObjectCreated:Put - PUT operation
- s3:ObjectCreated:Post - POST operation
- s3:ObjectCreated:Copy - COPY operation
- s3:ObjectRemoved:* - Any object deletion
- s3:ObjectReducedRedundancy - Lost redundancy

Best Practices:

1. Filter events by prefix/suffix in S3 notification
2. Use S3 event Time To Live (TTL) for old objects
3. Implement idempotent processing (handle duplicates)
4. Use dead letter queue (DLQ) for failed events
5. Monitor with CloudWatch metrics and alarms
6. Set appropriate Lambda timeout for large files
7. Use Lambda concurrency limits for cost control
*/
