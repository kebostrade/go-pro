module github.com/yourusername/cloud-cicd-go

go 1.22

require (
	// GCP SDKs
	cloud.google.com/go/storage v1.36.0
	cloud.google.com/go/pubsub v1.33.0
	cloud.google.com/go/firestore v1.14.0
	cloud.google.com/go/functions v1.15.4
	
	// AWS SDKs
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.26.8
	github.com/aws/aws-sdk-go-v2/service/sqs v1.29.7
	github.com/aws/aws-sdk-go-v2/service/sns v1.26.7
	github.com/aws/aws-lambda-go v1.46.0
	
	// Web frameworks
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/mux v1.8.1
	
	// Utilities
	github.com/joho/godotenv v1.5.1
	google.golang.org/api v0.155.0
)

