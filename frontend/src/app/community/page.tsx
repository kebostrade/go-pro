"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Users,
  MessageSquare,
  Trophy,
  Star,
  TrendingUp,
  Calendar,
  Clock,
  ThumbsUp,
  MessageCircle,
  Eye,
  Plus,
  Search,
  Filter,
  BookOpen,
  Code2,
  Lightbulb,
  HelpCircle,
  Zap,
  Award,
  Target,
  Heart,
  Share2
} from "lucide-react";
import Link from "next/link";

interface ForumPost {
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
}

interface CommunityMember {
  id: string;
  name: string;
  avatar: string;
  reputation: number;
  badge: string;
  contributions: number;
  joinedAt: string;
  specialties: string[];
}

interface CommunityEvent {
  id: string;
  title: string;
  description: string;
  type: "workshop" | "challenge" | "meetup" | "webinar";
  date: string;
  participants: number;
  maxParticipants?: number;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
}

export default function CommunityPage() {
  const [activeTab, setActiveTab] = useState("discussions");
  const [selectedCategory, setSelectedCategory] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  // Mock data for forum posts
  const forumPosts: ForumPost[] = [
    {
      id: "1",
      title: "Best practices for error handling in Go microservices",
      content: "I'm building a microservices architecture and wondering about the best patterns for error handling across services...",
      author: {
        name: "Alex Chen",
        avatar: "AC",
        reputation: 2450,
        badge: "Go Expert"
      },
      category: "Architecture",
      tags: ["microservices", "error-handling", "best-practices"],
      createdAt: "2 hours ago",
      replies: 12,
      views: 156,
      likes: 23,
      solved: false,
      pinned: true
    },
    {
      id: "2",
      title: "How to optimize goroutine performance?",
      content: "I have a concurrent application that processes thousands of requests. Looking for tips on goroutine optimization...",
      author: {
        name: "Sarah Johnson",
        avatar: "SJ",
        reputation: 1890,
        badge: "Concurrency Master"
      },
      category: "Performance",
      tags: ["goroutines", "performance", "concurrency"],
      createdAt: "4 hours ago",
      replies: 8,
      views: 89,
      likes: 15,
      solved: true,
      pinned: false
    },
    {
      id: "3",
      title: "Beginner question: Understanding interfaces",
      content: "I'm new to Go and struggling with interfaces. Can someone explain when and how to use them effectively?",
      author: {
        name: "Mike Rodriguez",
        avatar: "MR",
        reputation: 340,
        badge: "Go Learner"
      },
      category: "Beginner",
      tags: ["interfaces", "beginner", "concepts"],
      createdAt: "6 hours ago",
      replies: 15,
      views: 234,
      likes: 28,
      solved: true,
      pinned: false
    },
    {
      id: "4",
      title: "Weekly Challenge: Implement a Rate Limiter",
      content: "This week's coding challenge: Build a rate limiter using Go. Share your solutions and discuss different approaches!",
      author: {
        name: "GO-PRO Team",
        avatar: "GP",
        reputation: 5000,
        badge: "Moderator"
      },
      category: "Challenges",
      tags: ["challenge", "rate-limiter", "algorithms"],
      createdAt: "1 day ago",
      replies: 42,
      views: 567,
      likes: 89,
      solved: false,
      pinned: true
    }
  ];

  // Mock data for community members
  const topMembers: CommunityMember[] = [
    {
      id: "1",
      name: "Alex Chen",
      avatar: "AC",
      reputation: 2450,
      badge: "Go Expert",
      contributions: 156,
      joinedAt: "Jan 2023",
      specialties: ["Microservices", "Performance", "Architecture"]
    },
    {
      id: "2",
      name: "Sarah Johnson",
      avatar: "SJ",
      reputation: 1890,
      badge: "Concurrency Master",
      contributions: 98,
      joinedAt: "Mar 2023",
      specialties: ["Concurrency", "Goroutines", "Channels"]
    },
    {
      id: "3",
      name: "David Kim",
      avatar: "DK",
      reputation: 1650,
      badge: "Web Dev Pro",
      contributions: 87,
      joinedAt: "Feb 2023",
      specialties: ["Web Development", "APIs", "HTTP"]
    }
  ];

  // Mock data for community events
  const upcomingEvents: CommunityEvent[] = [
    {
      id: "1",
      title: "Go Concurrency Workshop",
      description: "Deep dive into goroutines, channels, and concurrent patterns",
      type: "workshop",
      date: "2024-01-15",
      participants: 45,
      maxParticipants: 50,
      difficulty: "Intermediate"
    },
    {
      id: "2",
      title: "Monthly Coding Challenge",
      description: "Build a distributed cache system using Go",
      type: "challenge",
      date: "2024-01-20",
      participants: 128,
      difficulty: "Advanced"
    },
    {
      id: "3",
      title: "Beginner's Go Meetup",
      description: "Q&A session for Go beginners with experienced developers",
      type: "meetup",
      date: "2024-01-18",
      participants: 67,
      maxParticipants: 100,
      difficulty: "Beginner"
    }
  ];

  const categories = ["all", "Beginner", "Architecture", "Performance", "Web Development", "Challenges", "General"];

  const filteredPosts = forumPosts.filter(post => {
    const matchesCategory = selectedCategory === "all" || post.category === selectedCategory;
    const matchesSearch = post.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         post.content.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         post.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()));
    
    return matchesCategory && matchesSearch;
  });

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case "Beginner": return <BookOpen className="h-4 w-4" />;
      case "Architecture": return <Target className="h-4 w-4" />;
      case "Performance": return <Zap className="h-4 w-4" />;
      case "Challenges": return <Trophy className="h-4 w-4" />;
      default: return <MessageSquare className="h-4 w-4" />;
    }
  };

  const getBadgeColor = (badge: string) => {
    switch (badge) {
      case "Go Expert": return "bg-purple-100 text-purple-800 border-purple-200";
      case "Concurrency Master": return "bg-blue-100 text-blue-800 border-blue-200";
      case "Web Dev Pro": return "bg-green-100 text-green-800 border-green-200";
      case "Moderator": return "bg-red-100 text-red-800 border-red-200";
      default: return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const communityStats = {
    totalMembers: 12847,
    activeToday: 1247,
    totalPosts: 8934,
    solvedQuestions: 6721
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container max-w-7xl mx-auto px-4 py-8 sm:px-6 sm:py-10 lg:px-8 lg:py-12">
        {/* Header */}
        <div className="mb-10 lg:mb-12">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6 lg:mb-8">
            <div className="mb-4 lg:mb-0">
              <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold tracking-tight mb-3 bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-transparent">
                Community Hub
              </h1>
              <p className="text-muted-foreground text-lg lg:text-xl max-w-2xl">
                Connect, learn, and grow with fellow Go developers
              </p>
            </div>
            <Link href="/community/new-post">
              <Button className="go-gradient text-white text-base px-6 py-3">
                <Plus className="mr-2 h-5 w-5" />
                New Post
              </Button>
            </Link>
          </div>

          {/* Community Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 lg:gap-6 mb-8">
          <Card>
            <CardContent className="p-4 text-center">
              <Users className="h-6 w-6 text-blue-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{communityStats.totalMembers.toLocaleString()}</div>
              <div className="text-sm text-muted-foreground">Members</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <TrendingUp className="h-6 w-6 text-green-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{communityStats.activeToday.toLocaleString()}</div>
              <div className="text-sm text-muted-foreground">Active Today</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <MessageSquare className="h-6 w-6 text-purple-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{communityStats.totalPosts.toLocaleString()}</div>
              <div className="text-sm text-muted-foreground">Discussions</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <HelpCircle className="h-6 w-6 text-orange-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{communityStats.solvedQuestions.toLocaleString()}</div>
              <div className="text-sm text-muted-foreground">Solved</div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-4 lg:w-[500px]">
          <TabsTrigger value="discussions">Discussions</TabsTrigger>
          <TabsTrigger value="members">Members</TabsTrigger>
          <TabsTrigger value="events">Events</TabsTrigger>
          <TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
        </TabsList>

        {/* Discussions Tab */}
        <TabsContent value="discussions" className="space-y-6">
          {/* Filters */}
          <Card>
            <CardContent className="p-4">
              <div className="flex flex-col md:flex-row gap-4">
                <div className="flex-1">
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <input
                      type="text"
                      placeholder="Search discussions..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                  </div>
                </div>
                <div className="flex gap-2">
                  <select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    className="px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                  >
                    {categories.map(category => (
                      <option key={category} value={category}>
                        {category === "all" ? "All Categories" : category}
                      </option>
                    ))}
                  </select>
                  <Button variant="outline" size="sm">
                    <Filter className="mr-2 h-4 w-4" />
                    Filter
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Forum Posts */}
          <div className="space-y-4">
            {filteredPosts.map((post) => (
              <Card key={post.id} className={`${post.pinned ? 'border-primary/50 bg-primary/5' : ''}`}>
                <CardContent className="p-6">
                  <div className="flex items-start space-x-4">
                    {/* Author Avatar */}
                    <div className="flex-shrink-0">
                      <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                        {post.author.avatar}
                      </div>
                    </div>

                    {/* Post Content */}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          {post.pinned && (
                            <Badge variant="secondary" className="text-xs">
                              Pinned
                            </Badge>
                          )}
                          <Badge className={getBadgeColor(post.author.badge)}>
                            {post.author.badge}
                          </Badge>
                          {post.solved && (
                            <Badge className="bg-green-100 text-green-800 border-green-200">
                              Solved
                            </Badge>
                          )}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {post.createdAt}
                        </div>
                      </div>

                      <Link href={`/community/post/${post.id}`}>
                        <h3 className="text-lg font-semibold mb-2 hover:text-primary cursor-pointer">
                          {post.title}
                        </h3>
                      </Link>

                      <p className="text-muted-foreground mb-3 line-clamp-2">
                        {post.content}
                      </p>

                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                          <div className="flex items-center space-x-1">
                            <MessageCircle className="h-4 w-4" />
                            <span>{post.replies}</span>
                          </div>
                          <div className="flex items-center space-x-1">
                            <Eye className="h-4 w-4" />
                            <span>{post.views}</span>
                          </div>
                          <div className="flex items-center space-x-1">
                            <ThumbsUp className="h-4 w-4" />
                            <span>{post.likes}</span>
                          </div>
                        </div>

                        <div className="flex items-center space-x-2">
                          <div className="flex items-center space-x-1">
                            {getCategoryIcon(post.category)}
                            <span className="text-sm text-muted-foreground">{post.category}</span>
                          </div>
                        </div>
                      </div>

                      {/* Tags */}
                      <div className="flex flex-wrap gap-1 mt-3">
                        {post.tags.map((tag) => (
                          <Badge key={tag} variant="outline" className="text-xs">
                            {tag}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {/* Members Tab */}
        <TabsContent value="members" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Top Contributors */}
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Star className="mr-2 h-5 w-5 text-yellow-500" />
                    Top Contributors
                  </CardTitle>
                  <CardDescription>
                    Most active community members this month
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {topMembers.map((member, index) => (
                      <div key={member.id} className="flex items-center justify-between p-4 rounded-lg bg-muted/50">
                        <div className="flex items-center space-x-4">
                          <div className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold ${
                            index === 0 ? "bg-yellow-500 text-white" :
                            index === 1 ? "bg-gray-400 text-white" :
                            index === 2 ? "bg-amber-600 text-white" :
                            "bg-primary text-primary-foreground"
                          }`}>
                            {index < 3 ? index + 1 : member.avatar}
                          </div>
                          <div>
                            <div className="font-medium">{member.name}</div>
                            <div className="text-sm text-muted-foreground">
                              {member.contributions} contributions • Joined {member.joinedAt}
                            </div>
                            <div className="flex flex-wrap gap-1 mt-1">
                              {member.specialties.slice(0, 2).map((specialty) => (
                                <Badge key={specialty} variant="outline" className="text-xs">
                                  {specialty}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-bold text-primary">{member.reputation}</div>
                          <div className="text-sm text-muted-foreground">reputation</div>
                          <Badge className={getBadgeColor(member.badge)}>
                            {member.badge}
                          </Badge>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Member Stats */}
            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Users className="mr-2 h-5 w-5 text-blue-500" />
                    Community Growth
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="text-center">
                      <div className="text-3xl font-bold text-primary">+247</div>
                      <div className="text-sm text-muted-foreground">New members this week</div>
                    </div>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>Active Members</span>
                        <span className="font-bold">89%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div className="bg-primary h-2 rounded-full" style={{ width: '89%' }}></div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Award className="mr-2 h-5 w-5 text-purple-500" />
                    Member Badges
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Go Expert</span>
                      <Badge className="bg-purple-100 text-purple-800 border-purple-200">23</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Concurrency Master</span>
                      <Badge className="bg-blue-100 text-blue-800 border-blue-200">45</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Web Dev Pro</span>
                      <Badge className="bg-green-100 text-green-800 border-green-200">67</Badge>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Go Learner</span>
                      <Badge className="bg-gray-100 text-gray-800 border-gray-200">1,234</Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        {/* Events Tab */}
        <TabsContent value="events" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {upcomingEvents.map((event) => (
              <Card key={event.id}>
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="text-xl mb-2">{event.title}</CardTitle>
                      <CardDescription>
                        {event.description}
                      </CardDescription>
                    </div>
                    <Badge variant="outline" className="ml-4">
                      {event.type}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-2 mt-3">
                    <Badge className={
                      event.difficulty === "Beginner" ? "bg-green-100 text-green-800 border-green-200" :
                      event.difficulty === "Intermediate" ? "bg-yellow-100 text-yellow-800 border-yellow-200" :
                      "bg-red-100 text-red-800 border-red-200"
                    }>
                      {event.difficulty}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center space-x-4">
                        <div className="flex items-center space-x-1">
                          <Calendar className="h-4 w-4 text-muted-foreground" />
                          <span>{new Date(event.date).toLocaleDateString()}</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <Users className="h-4 w-4 text-muted-foreground" />
                          <span>
                            {event.participants}
                            {event.maxParticipants && ` / ${event.maxParticipants}`}
                          </span>
                        </div>
                      </div>
                    </div>

                    {event.maxParticipants && (
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>Participants</span>
                          <span>{event.participants}/{event.maxParticipants}</span>
                        </div>
                        <div className="w-full bg-gray-200 rounded-full h-2">
                          <div
                            className="bg-primary h-2 rounded-full"
                            style={{ width: `${(event.participants / event.maxParticipants) * 100}%` }}
                          ></div>
                        </div>
                      </div>
                    )}

                    <Button className="w-full" variant="outline">
                      {event.maxParticipants && event.participants >= event.maxParticipants ?
                        "Join Waitlist" : "Join Event"
                      }
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Create Event CTA */}
          <Card className="bg-gradient-to-r from-primary/5 to-primary/10 border-primary/20">
            <CardContent className="p-8 text-center">
              <h3 className="text-2xl font-bold mb-4">Host Your Own Event</h3>
              <p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
                Share your knowledge with the community by hosting workshops, meetups, or coding challenges.
              </p>
              <Button className="go-gradient text-white">
                <Plus className="mr-2 h-5 w-5" />
                Create Event
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Leaderboard Tab */}
        <TabsContent value="leaderboard" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Monthly Leaderboard */}
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Trophy className="mr-2 h-5 w-5 text-yellow-500" />
                    Monthly Leaderboard
                  </CardTitle>
                  <CardDescription>
                    Top contributors this month
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {[
                      { rank: 1, name: "Alex Chen", points: 2450, avatar: "AC", badge: "Go Expert", change: "+12" },
                      { rank: 2, name: "Sarah Johnson", points: 2380, avatar: "SJ", badge: "Concurrency Master", change: "+5" },
                      { rank: 3, name: "David Kim", points: 2290, avatar: "DK", badge: "Web Dev Pro", change: "-1" },
                      { rank: 4, name: "Emma Wilson", points: 2150, avatar: "EW", badge: "Go Learner", change: "+8" },
                      { rank: 5, name: "Mike Rodriguez", points: 2050, avatar: "MR", badge: "Go Learner", change: "+3" },
                    ].map((user) => (
                      <div key={user.rank} className="flex items-center justify-between p-3 rounded-lg bg-muted/50">
                        <div className="flex items-center space-x-3">
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold ${
                            user.rank === 1 ? "bg-yellow-500 text-white" :
                            user.rank === 2 ? "bg-gray-400 text-white" :
                            user.rank === 3 ? "bg-amber-600 text-white" :
                            "bg-primary text-primary-foreground"
                          }`}>
                            {user.rank <= 3 ? user.rank : user.avatar}
                          </div>
                          <div>
                            <div className="font-medium">{user.name}</div>
                            <Badge className={getBadgeColor(user.badge)}>
                              {user.badge}
                            </Badge>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-bold">{user.points}</div>
                          <div className={`text-sm ${
                            user.change.startsWith('+') ? 'text-green-600' :
                            user.change.startsWith('-') ? 'text-red-600' :
                            'text-muted-foreground'
                          }`}>
                            {user.change}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Leaderboard Stats */}
            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <TrendingUp className="mr-2 h-5 w-5 text-green-500" />
                    Your Ranking
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-center space-y-4">
                    <div>
                      <div className="text-3xl font-bold text-primary">#47</div>
                      <div className="text-sm text-muted-foreground">Current Rank</div>
                    </div>
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>This Month</span>
                        <span className="font-bold text-green-600">+12 positions</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span>Points</span>
                        <span className="font-bold">1,247</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Award className="mr-2 h-5 w-5 text-purple-500" />
                    Point Categories
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3 text-sm">
                    <div className="flex justify-between">
                      <span>Helpful Answers</span>
                      <span className="font-bold">450 pts</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Questions Asked</span>
                      <span className="font-bold">120 pts</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Code Reviews</span>
                      <span className="font-bold">300 pts</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Community Events</span>
                      <span className="font-bold">200 pts</span>
                    </div>
                    <div className="border-t pt-2 flex justify-between font-bold">
                      <span>Total</span>
                      <span>1,070 pts</span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>
      </Tabs>
      </div>
    </div>
  );
}
