"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  User,
  Calendar,
  MessageSquare,
  Trophy,
  Star,
  Award,
  TrendingUp,
  MapPin,
  Link as LinkIcon,
  Github,
  Twitter,
  Mail,
  Users,
  BookOpen
} from "lucide-react";
import Link from "next/link";

interface UserProfileCardProps {
  id: string;
  name: string;
  avatar: string;
  badge: string;
  reputation: number;
  contributions: number;
  joinedAt: string;
  location?: string;
  bio?: string;
  website?: string;
  github?: string;
  twitter?: string;
  email?: string;
  specialties: string[];
  stats: {
    posts: number;
    answers: number;
    helpfulAnswers: number;
    followers: number;
    following: number;
  };
  isFollowing?: boolean;
  isCurrentUser?: boolean;
  onFollow?: () => void;
  onMessage?: () => void;
}

const UserProfileCard = ({
  id,
  name,
  avatar,
  badge,
  reputation,
  contributions,
  joinedAt,
  location,
  bio,
  website,
  github,
  twitter,
  email,
  specialties,
  stats,
  isFollowing = false,
  isCurrentUser = false,
  onFollow,
  onMessage
}: UserProfileCardProps) => {
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

  const getReputationLevel = (reputation: number) => {
    if (reputation >= 5000) return { level: "Expert", color: "text-purple-600" };
    if (reputation >= 2000) return { level: "Advanced", color: "text-blue-600" };
    if (reputation >= 500) return { level: "Intermediate", color: "text-green-600" };
    return { level: "Beginner", color: "text-gray-600" };
  };

  const reputationLevel = getReputationLevel(reputation);

  return (
    <Card className="transition-all hover:shadow-md">
      <CardHeader className="text-center">
        {/* Avatar and Basic Info */}
        <div className="flex flex-col items-center space-y-3">
          <div className="w-16 h-16 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-xl">
            {avatar}
          </div>
          
          <div>
            <CardTitle className="text-xl">{name}</CardTitle>
            {isCurrentUser && (
              <Badge variant="secondary" className="mt-1">You</Badge>
            )}
          </div>

          <div className="flex items-center space-x-2">
            <Badge className={getBadgeColor(badge)}>
              {badge}
            </Badge>
            <Badge variant="outline" className={reputationLevel.color}>
              {reputationLevel.level}
            </Badge>
          </div>
        </div>

        {/* Bio */}
        {bio && (
          <CardDescription className="text-center mt-3">
            {bio}
          </CardDescription>
        )}
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Reputation and Stats */}
        <div className="grid grid-cols-2 gap-4 text-center">
          <div>
            <div className="text-2xl font-bold text-primary">{reputation.toLocaleString()}</div>
            <div className="text-sm text-muted-foreground">Reputation</div>
          </div>
          <div>
            <div className="text-2xl font-bold text-green-600">{contributions}</div>
            <div className="text-sm text-muted-foreground">Contributions</div>
          </div>
        </div>

        {/* Detailed Stats */}
        <div className="space-y-2 text-sm">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <MessageSquare className="h-4 w-4 text-muted-foreground" />
              <span>Posts</span>
            </div>
            <span className="font-medium">{stats.posts}</span>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <BookOpen className="h-4 w-4 text-muted-foreground" />
              <span>Answers</span>
            </div>
            <span className="font-medium">{stats.answers}</span>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Trophy className="h-4 w-4 text-muted-foreground" />
              <span>Helpful Answers</span>
            </div>
            <span className="font-medium">{stats.helpfulAnswers}</span>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Users className="h-4 w-4 text-muted-foreground" />
              <span>Followers</span>
            </div>
            <span className="font-medium">{stats.followers}</span>
          </div>
        </div>

        {/* Specialties */}
        {specialties.length > 0 && (
          <div>
            <div className="text-sm font-medium mb-2">Specialties</div>
            <div className="flex flex-wrap gap-1">
              {specialties.map((specialty) => (
                <Badge key={specialty} variant="secondary" className="text-xs">
                  {specialty}
                </Badge>
              ))}
            </div>
          </div>
        )}

        {/* Additional Info */}
        <div className="space-y-2 text-sm text-muted-foreground">
          <div className="flex items-center space-x-2">
            <Calendar className="h-4 w-4" />
            <span>Joined {joinedAt}</span>
          </div>
          {location && (
            <div className="flex items-center space-x-2">
              <MapPin className="h-4 w-4" />
              <span>{location}</span>
            </div>
          )}
        </div>

        {/* Social Links */}
        {(website || github || twitter || email) && (
          <div className="flex items-center justify-center space-x-2">
            {website && (
              <Button variant="outline" size="sm" asChild>
                <a href={website} target="_blank" rel="noopener noreferrer">
                  <LinkIcon className="h-4 w-4" />
                </a>
              </Button>
            )}
            {github && (
              <Button variant="outline" size="sm" asChild>
                <a href={`https://github.com/${github}`} target="_blank" rel="noopener noreferrer">
                  <Github className="h-4 w-4" />
                </a>
              </Button>
            )}
            {twitter && (
              <Button variant="outline" size="sm" asChild>
                <a href={`https://twitter.com/${twitter}`} target="_blank" rel="noopener noreferrer">
                  <Twitter className="h-4 w-4" />
                </a>
              </Button>
            )}
            {email && (
              <Button variant="outline" size="sm" asChild>
                <a href={`mailto:${email}`}>
                  <Mail className="h-4 w-4" />
                </a>
              </Button>
            )}
          </div>
        )}

        {/* Action Buttons */}
        {!isCurrentUser && (
          <div className="flex space-x-2">
            <Button
              variant={isFollowing ? "outline" : "default"}
              className="flex-1"
              onClick={onFollow}
            >
              {isFollowing ? "Following" : "Follow"}
            </Button>
            <Button variant="outline" className="flex-1" onClick={onMessage}>
              <MessageSquare className="mr-2 h-4 w-4" />
              Message
            </Button>
          </div>
        )}

        {isCurrentUser && (
          <Link href="/profile/edit">
            <Button variant="outline" className="w-full">
              <User className="mr-2 h-4 w-4" />
              Edit Profile
            </Button>
          </Link>
        )}
      </CardContent>
    </Card>
  );
};

export default UserProfileCard;
