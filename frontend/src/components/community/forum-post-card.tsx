"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  MessageCircle,
  Eye,
  ThumbsUp,
  Clock,
  CheckCircle,
  Pin,
  BookOpen,
  Code2,
  Trophy,
  HelpCircle,
  Lightbulb,
  Target
} from "lucide-react";
import Link from "next/link";

interface ForumPostCardProps {
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
  replies: number;
  views: number;
  likes: number;
  solved: boolean;
  pinned: boolean;
  onClick?: () => void;
}

const ForumPostCard = ({
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
  solved,
  pinned,
  onClick
}: ForumPostCardProps) => {
  const getCategoryIcon = (category: string) => {
    switch (category) {
      case "Beginner": return <BookOpen className="h-4 w-4" />;
      case "Architecture": return <Target className="h-4 w-4" />;
      case "Performance": return <Trophy className="h-4 w-4" />;
      case "Challenges": return <Code2 className="h-4 w-4" />;
      case "Help": return <HelpCircle className="h-4 w-4" />;
      case "Tips": return <Lightbulb className="h-4 w-4" />;
      default: return <MessageCircle className="h-4 w-4" />;
    }
  };

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

  const getCategoryColor = (category: string) => {
    switch (category) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Architecture": return "text-blue-600 bg-blue-50 border-blue-200";
      case "Performance": return "text-purple-600 bg-purple-50 border-purple-200";
      case "Challenges": return "text-orange-600 bg-orange-50 border-orange-200";
      case "Help": return "text-red-600 bg-red-50 border-red-200";
      case "Tips": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  const handleClick = () => {
    if (onClick) {
      onClick();
    }
  };

  return (
    <Card 
      className={`transition-all hover:shadow-md cursor-pointer ${
        pinned ? 'border-primary/50 bg-primary/5' : ''
      }`}
      onClick={handleClick}
    >
      <CardContent className="p-6">
        <div className="flex items-start space-x-4">
          {/* Author Avatar */}
          <div className="flex-shrink-0">
            <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
              {author.avatar}
            </div>
          </div>

          {/* Post Content */}
          <div className="flex-1 min-w-0">
            {/* Header Badges */}
            <div className="flex items-center space-x-2 mb-2">
              {pinned && (
                <Badge variant="secondary" className="text-xs">
                  <Pin className="mr-1 h-3 w-3" />
                  Pinned
                </Badge>
              )}
              <Badge className={getBadgeColor(author.badge)}>
                {author.badge}
              </Badge>
              {solved && (
                <Badge className="bg-green-100 text-green-800 border-green-200">
                  <CheckCircle className="mr-1 h-3 w-3" />
                  Solved
                </Badge>
              )}
            </div>

            {/* Title */}
            <Link href={`/community/post/${id}`}>
              <h3 className="text-lg font-semibold mb-2 hover:text-primary cursor-pointer line-clamp-2">
                {title}
              </h3>
            </Link>

            {/* Content Preview */}
            <p className="text-muted-foreground mb-3 line-clamp-2 text-sm">
              {content}
            </p>

            {/* Author and Time Info */}
            <div className="flex items-center space-x-4 text-sm text-muted-foreground mb-3">
              <span className="font-medium text-foreground">{author.name}</span>
              <span>{author.reputation} reputation</span>
              <div className="flex items-center space-x-1">
                <Clock className="h-3 w-3" />
                <span>{createdAt}</span>
              </div>
            </div>

            {/* Stats and Category */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                <div className="flex items-center space-x-1">
                  <MessageCircle className="h-4 w-4" />
                  <span>{replies}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Eye className="h-4 w-4" />
                  <span>{views}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <ThumbsUp className="h-4 w-4" />
                  <span>{likes}</span>
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <div className="flex items-center space-x-1">
                  {getCategoryIcon(category)}
                  <Badge className={getCategoryColor(category)}>
                    {category}
                  </Badge>
                </div>
              </div>
            </div>

            {/* Tags */}
            {tags.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-3">
                {tags.slice(0, 4).map((tag) => (
                  <Badge key={tag} variant="outline" className="text-xs">
                    {tag}
                  </Badge>
                ))}
                {tags.length > 4 && (
                  <Badge variant="outline" className="text-xs">
                    +{tags.length - 4} more
                  </Badge>
                )}
              </div>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default ForumPostCard;
