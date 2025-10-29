"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import {
  ArrowLeft,
  ThumbsUp,
  ThumbsDown,
  MessageCircle,
  Share2,
  Flag,
  CheckCircle,
  Clock,
  Eye,
  Home,
  ChevronRight,
  Send,
  Heart,
  Bookmark,
  MoreHorizontal
} from "lucide-react";
import Link from "next/link";
import CodeEditor from "@/components/learning/code-editor";

interface Reply {
  id: string;
  content: string;
  author: {
    name: string;
    avatar: string;
    reputation: number;
    badge: string;
  };
  createdAt: string;
  likes: number;
  dislikes: number;
  isAccepted: boolean;
  hasCode: boolean;
  codeContent?: string;
}

interface PostData {
  id: string;
  title: string;
  content: string;
  author: {
    name: string;
    avatar: string;
    reputation: number;
    badge: string;
  };
  category: string;
  tags: string[];
  createdAt: string;
  replies: Reply[];
  views: number;
  likes: number;
  dislikes: number;
  solved: boolean;
  pinned: boolean;
  hasCode: boolean;
  codeContent?: string;
}

export default function PostDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [postData, setPostData] = useState<PostData | null>(null);
  const [loading, setLoading] = useState(true);
  const [newReply, setNewReply] = useState("");
  const [userVote, setUserVote] = useState<"up" | "down" | null>(null);
  const [isBookmarked, setIsBookmarked] = useState(false);

  const postId = params.id as string;

  useEffect(() => {
    // Mock data loading
    const mockPost: PostData = {
      id: postId,
      title: "Best practices for error handling in Go microservices",
      content: `I'm building a microservices architecture with Go and I'm struggling with consistent error handling patterns across services. 

Here are my main concerns:

1. **Error Propagation**: How should errors be propagated between services? Should I use gRPC status codes or custom error types?

2. **Logging**: What's the best way to maintain correlation IDs across service boundaries for debugging?

3. **Circuit Breakers**: When should I implement circuit breakers and how do they integrate with error handling?

I've been looking at some patterns but would love to hear from the community about what works in production environments.

Here's a simplified example of what I'm currently doing:`,
      author: {
        name: "Alex Chen",
        avatar: "AC",
        reputation: 2450,
        badge: "Go Expert"
      },
      category: "Architecture",
      tags: ["microservices", "error-handling", "best-practices", "architecture"],
      createdAt: "2 hours ago",
      views: 156,
      likes: 23,
      dislikes: 2,
      solved: false,
      pinned: true,
      hasCode: true,
      codeContent: `package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type ServiceError struct {
    Code    string \`json:"code"\`
    Message string \`json:"message"\`
    Details map[string]interface{} \`json:"details,omitempty"\`
}

func (e *ServiceError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func handleUserRequest(ctx context.Context, userID string) error {
    // Call user service
    user, err := userService.GetUser(ctx, userID)
    if err != nil {
        // How should I handle this error?
        return &ServiceError{
            Code:    "USER_SERVICE_ERROR",
            Message: "Failed to fetch user",
            Details: map[string]interface{}{"user_id": userID},
        }
    }
    
    // Process user...
    return nil
}`,
      replies: [
        {
          id: "1",
          content: `Great question! I've been working with Go microservices for 3 years and here's what I've learned:

**For Error Propagation:**
- Use gRPC status codes for service-to-service communication
- Create a common error package that all services can import
- Implement error wrapping to maintain context

**For Logging:**
- Use OpenTelemetry for distributed tracing
- Pass correlation IDs in context
- Structured logging with consistent fields

Here's how I structure my error handling:`,
          author: {
            name: "Sarah Johnson",
            avatar: "SJ",
            reputation: 1890,
            badge: "Concurrency Master"
          },
          createdAt: "1 hour ago",
          likes: 15,
          dislikes: 0,
          isAccepted: true,
          hasCode: true,
          codeContent: `package errors

import (
    "context"
    "fmt"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type ErrorCode string

const (
    ErrCodeNotFound     ErrorCode = "NOT_FOUND"
    ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
    ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
)

type ServiceError struct {
    Code       ErrorCode
    Message    string
    Cause      error
    TraceID    string
    ServiceID  string
}

func (e *ServiceError) Error() string {
    return fmt.Sprintf("[%s] %s: %s", e.ServiceID, e.Code, e.Message)
}

func NewServiceError(code ErrorCode, message string, cause error) *ServiceError {
    return &ServiceError{
        Code:    code,
        Message: message,
        Cause:   cause,
    }
}

func ToGRPCError(err error) error {
    if serviceErr, ok := err.(*ServiceError); ok {
        switch serviceErr.Code {
        case ErrCodeNotFound:
            return status.Error(codes.NotFound, serviceErr.Message)
        case ErrCodeInvalidInput:
            return status.Error(codes.InvalidArgument, serviceErr.Message)
        default:
            return status.Error(codes.Internal, serviceErr.Message)
        }
    }
    return status.Error(codes.Internal, "Internal server error")
}`
        },
        {
          id: "2",
          content: `I'd also recommend looking into the **Circuit Breaker pattern**. We use hystrix-go in our services and it's been a game changer for handling cascading failures.

The key is to fail fast and provide meaningful fallbacks. Don't let one slow service bring down your entire system.`,
          author: {
            name: "Mike Rodriguez",
            avatar: "MR",
            reputation: 1240,
            badge: "Go Developer"
          },
          createdAt: "45 minutes ago",
          likes: 8,
          dislikes: 0,
          isAccepted: false,
          hasCode: false
        }
      ]
    };

    setPostData(mockPost);
    setLoading(false);
  }, [postId]);

  const handleVote = (type: "up" | "down") => {
    if (userVote === type) {
      setUserVote(null);
    } else {
      setUserVote(type);
    }
  };

  const handleReplySubmit = () => {
    if (newReply.trim()) {
      // Mock reply submission
      console.log("Submitting reply:", newReply);
      setNewReply("");
    }
  };

  const getBadgeColor = (badge: string) => {
    switch (badge) {
      case "Go Expert": return "bg-purple-100 text-purple-800 border-purple-200";
      case "Concurrency Master": return "bg-blue-100 text-blue-800 border-blue-200";
      case "Go Developer": return "bg-green-100 text-green-800 border-green-200";
      case "Moderator": return "bg-red-100 text-red-800 border-red-200";
      default: return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  if (loading || !postData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
          <div className="flex items-center justify-center min-h-[60vh]">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
              <p className="text-responsive text-muted-foreground">Loading post...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Breadcrumb Navigation */}
        <div className="flex items-center space-x-2 text-sm text-muted-foreground mb-6 lg:mb-8">
        <Link href="/" className="hover:text-primary">
          <Home className="h-4 w-4" />
        </Link>
        <ChevronRight className="h-4 w-4" />
        <Link href="/community" className="hover:text-primary">
          Community
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-foreground">Discussion</span>
      </div>

      {/* Post Header */}
      <div className="mb-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-3">
              {postData.pinned && (
                <Badge variant="secondary">Pinned</Badge>
              )}
              {postData.solved && (
                <Badge className="bg-green-100 text-green-800 border-green-200">
                  <CheckCircle className="mr-1 h-3 w-3" />
                  Solved
                </Badge>
              )}
              <Badge variant="outline">{postData.category}</Badge>
            </div>
            <h1 className="text-3xl font-bold tracking-tight mb-4">{postData.title}</h1>
            
            <div className="flex items-center space-x-6 text-sm text-muted-foreground">
              <div className="flex items-center space-x-1">
                <Eye className="h-4 w-4" />
                <span>{postData.views} views</span>
              </div>
              <div className="flex items-center space-x-1">
                <MessageCircle className="h-4 w-4" />
                <span>{postData.replies.length} replies</span>
              </div>
              <div className="flex items-center space-x-1">
                <Clock className="h-4 w-4" />
                <span>{postData.createdAt}</span>
              </div>
            </div>
          </div>
          <div className="ml-6">
            <Link href="/community">
              <Button variant="outline" size="sm">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Community
              </Button>
            </Link>
          </div>
        </div>

        {/* Tags */}
        <div className="flex flex-wrap gap-2">
          {postData.tags.map((tag) => (
            <Badge key={tag} variant="outline" className="text-sm">
              {tag}
            </Badge>
          ))}
        </div>
      </div>

      {/* Original Post */}
      <Card className="mb-6">
        <CardContent className="p-6">
          <div className="flex items-start space-x-4">
            {/* Author Info */}
            <div className="flex-shrink-0 text-center">
              <div className="w-12 h-12 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-lg mb-2">
                {postData.author.avatar}
              </div>
              <div className="text-sm">
                <div className="font-medium">{postData.author.name}</div>
                <div className="text-muted-foreground">{postData.author.reputation}</div>
                <Badge className={getBadgeColor(postData.author.badge)}>
                  {postData.author.badge}
                </Badge>
              </div>
            </div>

            {/* Post Content */}
            <div className="flex-1">
              <div className="prose dark:prose-invert max-w-none mb-4">
                <div dangerouslySetInnerHTML={{ 
                  __html: postData.content.replace(/\n/g, '<br>').replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') 
                }} />
              </div>

              {/* Code Block */}
              {postData.hasCode && postData.codeContent && (
                <div className="mb-4">
                  <CodeEditor
                    title="Current Implementation"
                    description="Here's the code I'm currently using"
                    initialCode={postData.codeContent}
                    language="go"
                    readOnly={true}
                  />
                </div>
              )}

              {/* Post Actions */}
              <div className="flex items-center justify-between pt-4 border-t">
                <div className="flex items-center space-x-4">
                  <div className="flex items-center space-x-2">
                    <Button
                      variant={userVote === "up" ? "default" : "outline"}
                      size="sm"
                      onClick={() => handleVote("up")}
                    >
                      <ThumbsUp className="mr-1 h-4 w-4" />
                      {postData.likes}
                    </Button>
                    <Button
                      variant={userVote === "down" ? "default" : "outline"}
                      size="sm"
                      onClick={() => handleVote("down")}
                    >
                      <ThumbsDown className="mr-1 h-4 w-4" />
                      {postData.dislikes}
                    </Button>
                  </div>
                  <Button variant="outline" size="sm">
                    <Share2 className="mr-1 h-4 w-4" />
                    Share
                  </Button>
                  <Button
                    variant={isBookmarked ? "default" : "outline"}
                    size="sm"
                    onClick={() => setIsBookmarked(!isBookmarked)}
                  >
                    <Bookmark className="mr-1 h-4 w-4" />
                    {isBookmarked ? "Saved" : "Save"}
                  </Button>
                </div>
                <Button variant="outline" size="sm">
                  <Flag className="mr-1 h-4 w-4" />
                  Report
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Replies Section */}
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h2 className="text-2xl font-bold">
            {postData.replies.length} {postData.replies.length === 1 ? 'Reply' : 'Replies'}
          </h2>
          <div className="flex items-center space-x-2">
            <span className="text-sm text-muted-foreground">Sort by:</span>
            <select className="text-sm border rounded px-2 py-1">
              <option>Most helpful</option>
              <option>Newest first</option>
              <option>Oldest first</option>
            </select>
          </div>
        </div>

        {/* Reply List */}
        {postData.replies.map((reply) => (
          <Card key={reply.id} className={reply.isAccepted ? "border-green-200 bg-green-50/50" : ""}>
            <CardContent className="p-6">
              <div className="flex items-start space-x-4">
                {/* Author Info */}
                <div className="flex-shrink-0 text-center">
                  <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold mb-2">
                    {reply.author.avatar}
                  </div>
                  <div className="text-sm">
                    <div className="font-medium">{reply.author.name}</div>
                    <div className="text-muted-foreground">{reply.author.reputation}</div>
                    <Badge className={getBadgeColor(reply.author.badge)}>
                      {reply.author.badge}
                    </Badge>
                  </div>
                </div>

                {/* Reply Content */}
                <div className="flex-1">
                  {reply.isAccepted && (
                    <div className="flex items-center space-x-2 mb-3">
                      <CheckCircle className="h-5 w-5 text-green-600" />
                      <span className="text-sm font-medium text-green-800">Accepted Answer</span>
                    </div>
                  )}

                  <div className="prose dark:prose-invert max-w-none mb-4">
                    <div dangerouslySetInnerHTML={{ 
                      __html: reply.content.replace(/\n/g, '<br>').replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') 
                    }} />
                  </div>

                  {/* Code Block */}
                  {reply.hasCode && reply.codeContent && (
                    <div className="mb-4">
                      <CodeEditor
                        title="Suggested Solution"
                        description="Here's my recommended approach"
                        initialCode={reply.codeContent}
                        language="go"
                        readOnly={true}
                      />
                    </div>
                  )}

                  {/* Reply Actions */}
                  <div className="flex items-center justify-between pt-4 border-t text-sm">
                    <div className="flex items-center space-x-4">
                      <div className="flex items-center space-x-2">
                        <Button variant="outline" size="sm">
                          <ThumbsUp className="mr-1 h-3 w-3" />
                          {reply.likes}
                        </Button>
                        <Button variant="outline" size="sm">
                          <ThumbsDown className="mr-1 h-3 w-3" />
                          {reply.dislikes}
                        </Button>
                      </div>
                      <Button variant="outline" size="sm">
                        Reply
                      </Button>
                      {!postData.solved && !reply.isAccepted && (
                        <Button variant="outline" size="sm" className="text-green-600 border-green-200 hover:bg-green-50">
                          <CheckCircle className="mr-1 h-3 w-3" />
                          Accept Answer
                        </Button>
                      )}
                    </div>
                    <div className="text-muted-foreground">
                      {reply.createdAt}
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}

        {/* Add Reply */}
        <Card>
          <CardHeader>
            <CardTitle>Add Your Reply</CardTitle>
            <CardDescription>
              Share your knowledge and help the community
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <Textarea
                placeholder="Write your reply here... You can include code examples and explanations."
                value={newReply}
                onChange={(e) => setNewReply(e.target.value)}
                rows={6}
              />
              <div className="flex items-center justify-between">
                <div className="text-sm text-muted-foreground">
                  Tip: Use **bold** for emphasis and include code examples to help others
                </div>
                <Button onClick={handleReplySubmit} className="go-gradient text-white">
                  <Send className="mr-2 h-4 w-4" />
                  Post Reply
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
      </div>
    </div>
  );
}
