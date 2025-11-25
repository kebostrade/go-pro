"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import {
  ThumbsUp,
  ThumbsDown,
  MessageCircle,
  Share2,
  Flag,
  CheckCircle,
  Clock,
  MoreHorizontal,
  Reply,
  Edit,
  Trash2,
  Award
} from "lucide-react";
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
  replies?: Reply[];
}

interface DiscussionThreadProps {
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
  hasCode?: boolean;
  codeContent?: string;
  onLike?: () => void;
  onDislike?: () => void;
  onReply?: (content: string) => void;
  onAcceptAnswer?: (replyId: string) => void;
}

const DiscussionThread = ({
  id,
  title,
  content,
  author,
  category,
  tags,
  createdAt,
  replies,
  views,
  likes,
  dislikes,
  solved,
  hasCode,
  codeContent,
  onLike,
  onDislike,
  onReply,
  onAcceptAnswer
}: DiscussionThreadProps) => {
  const [newReply, setNewReply] = useState("");
  const [showReplyForm, setShowReplyForm] = useState(false);
  const [userVote, setUserVote] = useState<"up" | "down" | null>(null);

  const getBadgeColor = (badge: string) => {
    switch (badge) {
      case "Go Expert": return "bg-purple-100 text-purple-800 border-purple-200";
      case "Concurrency Master": return "bg-blue-100 text-blue-800 border-blue-200";
      case "Web Dev Pro": return "bg-green-100 text-green-800 border-green-200";
      case "Go Developer": return "bg-indigo-100 text-indigo-800 border-indigo-200";
      case "Moderator": return "bg-red-100 text-red-800 border-red-200";
      case "Go Learner": return "bg-gray-100 text-gray-800 border-gray-200";
      default: return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const handleVote = (type: "up" | "down") => {
    if (userVote === type) {
      setUserVote(null);
    } else {
      setUserVote(type);
      if (type === "up" && onLike) {
        onLike();
      } else if (type === "down" && onDislike) {
        onDislike();
      }
    }
  };

  const handleReplySubmit = () => {
    if (newReply.trim() && onReply) {
      onReply(newReply);
      setNewReply("");
      setShowReplyForm(false);
    }
  };

  const ReplyComponent = ({ reply, depth = 0 }: { reply: Reply; depth?: number }) => (
    <div className={`${depth > 0 ? 'ml-8 mt-4' : ''}`}>
      <Card className={reply.isAccepted ? "border-green-200 bg-green-50/50" : ""}>
        <CardContent className="p-4">
          <div className="flex items-start space-x-3">
            {/* Author Avatar */}
            <div className="flex-shrink-0">
              <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-sm">
                {reply.author.avatar}
              </div>
            </div>

            {/* Reply Content */}
            <div className="flex-1">
              {reply.isAccepted && (
                <div className="flex items-center space-x-2 mb-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <span className="text-sm font-medium text-green-800">Accepted Answer</span>
                  <Award className="h-4 w-4 text-green-600" />
                </div>
              )}

              <div className="flex items-center space-x-2 mb-2">
                <span className="font-medium text-sm">{reply.author.name}</span>
                <Badge className={getBadgeColor(reply.author.badge)}>
                  {reply.author.badge}
                </Badge>
                <span className="text-xs text-muted-foreground">{reply.author.reputation} rep</span>
                <span className="text-xs text-muted-foreground">•</span>
                <span className="text-xs text-muted-foreground">{reply.createdAt}</span>
              </div>

              <div className="prose dark:prose-invert max-w-none text-sm mb-3">
                <div dangerouslySetInnerHTML={{ 
                  __html: reply.content.replace(/\n/g, '<br>').replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') 
                }} />
              </div>

              {/* Code Block */}
              {reply.hasCode && reply.codeContent && (
                <div className="mb-3">
                  <CodeEditor
                    title="Code Example"
                    description="Suggested solution"
                    initialCode={reply.codeContent}
                    language="go"
                    readOnly={true}
                  />
                </div>
              )}

              {/* Reply Actions */}
              <div className="flex items-center justify-between text-xs">
                <div className="flex items-center space-x-3">
                  <div className="flex items-center space-x-1">
                    <Button variant="ghost" size="sm" className="h-6 px-2">
                      <ThumbsUp className="mr-1 h-3 w-3" />
                      {reply.likes}
                    </Button>
                    <Button variant="ghost" size="sm" className="h-6 px-2">
                      <ThumbsDown className="mr-1 h-3 w-3" />
                      {reply.dislikes}
                    </Button>
                  </div>
                  <Button variant="ghost" size="sm" className="h-6 px-2">
                    <Reply className="mr-1 h-3 w-3" />
                    Reply
                  </Button>
                  {!solved && !reply.isAccepted && onAcceptAnswer && (
                    <Button 
                      variant="ghost" 
                      size="sm" 
                      className="h-6 px-2 text-green-600 hover:text-green-700"
                      onClick={() => onAcceptAnswer(reply.id)}
                    >
                      <CheckCircle className="mr-1 h-3 w-3" />
                      Accept
                    </Button>
                  )}
                </div>
                <Button variant="ghost" size="sm" className="h-6 px-2">
                  <MoreHorizontal className="h-3 w-3" />
                </Button>
              </div>

              {/* Nested Replies */}
              {reply.replies && reply.replies.length > 0 && (
                <div className="mt-4 space-y-4">
                  {reply.replies.map((nestedReply) => (
                    <ReplyComponent key={nestedReply.id} reply={nestedReply} depth={depth + 1} />
                  ))}
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Original Post */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                {author.avatar}
              </div>
              <div>
                <div className="flex items-center space-x-2">
                  <span className="font-medium">{author.name}</span>
                  <Badge className={getBadgeColor(author.badge)}>
                    {author.badge}
                  </Badge>
                  <span className="text-sm text-muted-foreground">{author.reputation} rep</span>
                </div>
                <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                  <Clock className="h-3 w-3" />
                  <span>{createdAt}</span>
                  <span>•</span>
                  <span>{views} views</span>
                </div>
              </div>
            </div>
            {solved && (
              <Badge className="bg-green-100 text-green-800 border-green-200">
                <CheckCircle className="mr-1 h-3 w-3" />
                Solved
              </Badge>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <div className="prose dark:prose-invert max-w-none mb-4">
            <div dangerouslySetInnerHTML={{ 
              __html: content.replace(/\n/g, '<br>').replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') 
            }} />
          </div>

          {/* Code Block */}
          {hasCode && codeContent && (
            <div className="mb-4">
              <CodeEditor
                title="Code Example"
                description="Current implementation"
                initialCode={codeContent}
                language="go"
                readOnly={true}
              />
            </div>
          )}

          {/* Tags */}
          <div className="flex flex-wrap gap-1 mb-4">
            {tags.map((tag) => (
              <Badge key={tag} variant="outline" className="text-xs">
                {tag}
              </Badge>
            ))}
          </div>

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
                  {likes}
                </Button>
                <Button
                  variant={userVote === "down" ? "default" : "outline"}
                  size="sm"
                  onClick={() => handleVote("down")}
                >
                  <ThumbsDown className="mr-1 h-4 w-4" />
                  {dislikes}
                </Button>
              </div>
              <Button variant="outline" size="sm">
                <Share2 className="mr-1 h-4 w-4" />
                Share
              </Button>
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => setShowReplyForm(!showReplyForm)}
              >
                <MessageCircle className="mr-1 h-4 w-4" />
                Reply
              </Button>
            </div>
            <Button variant="outline" size="sm">
              <Flag className="mr-1 h-4 w-4" />
              Report
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Reply Form */}
      {showReplyForm && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Add Your Reply</CardTitle>
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
                <div className="flex space-x-2">
                  <Button variant="outline" onClick={() => setShowReplyForm(false)}>
                    Cancel
                  </Button>
                  <Button onClick={handleReplySubmit} className="go-gradient text-white">
                    Post Reply
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Replies */}
      {replies.length > 0 && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold">
            {replies.length} {replies.length === 1 ? 'Reply' : 'Replies'}
          </h3>
          {replies.map((reply) => (
            <ReplyComponent key={reply.id} reply={reply} />
          ))}
        </div>
      )}
    </div>
  );
};

export default DiscussionThread;
