"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Trophy,
  Medal,
  Award,
  TrendingUp,
  TrendingDown,
  Minus,
  Crown,
  Star,
  Zap
} from "lucide-react";

interface LeaderboardEntry {
  rank: number;
  name: string;
  avatar: string;
  points: number;
  badge: string;
  streak: number;
  change: number; // Position change from last period
  isCurrentUser?: boolean;
}

interface LeaderboardProps {
  title?: string;
  description?: string;
  entries: LeaderboardEntry[];
  period?: "daily" | "weekly" | "monthly" | "all-time";
  showChange?: boolean;
  maxEntries?: number;
}

const Leaderboard = ({
  title = "Leaderboard",
  description = "Top performers this period",
  entries,
  period = "weekly",
  showChange = true,
  maxEntries = 10
}: LeaderboardProps) => {
  const getRankIcon = (rank: number) => {
    switch (rank) {
      case 1:
        return <Crown className="h-5 w-5 text-yellow-500" />;
      case 2:
        return <Medal className="h-5 w-5 text-gray-400" />;
      case 3:
        return <Award className="h-5 w-5 text-amber-600" />;
      default:
        return <div className="w-5 h-5 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-xs font-bold">{rank}</div>;
    }
  };

  const getRankBackground = (rank: number, isCurrentUser: boolean = false) => {
    if (isCurrentUser) {
      return "bg-primary/10 border-primary/20";
    }
    
    switch (rank) {
      case 1:
        return "bg-gradient-to-r from-yellow-50 to-yellow-100 border-yellow-200";
      case 2:
        return "bg-gradient-to-r from-gray-50 to-gray-100 border-gray-200";
      case 3:
        return "bg-gradient-to-r from-amber-50 to-amber-100 border-amber-200";
      default:
        return "bg-muted/50 border-border";
    }
  };

  const getBadgeColor = (badge: string) => {
    switch (badge) {
      case "Go Expert": return "bg-purple-100 text-purple-800 border-purple-200";
      case "Concurrency Master": return "bg-blue-100 text-blue-800 border-blue-200";
      case "Web Dev Pro": return "bg-green-100 text-green-800 border-green-200";
      case "Go Developer": return "bg-indigo-100 text-indigo-800 border-indigo-200";
      case "Go Learner": return "bg-gray-100 text-gray-800 border-gray-200";
      default: return "bg-gray-100 text-gray-800 border-gray-200";
    }
  };

  const getChangeIcon = (change: number) => {
    if (change > 0) {
      return <TrendingUp className="h-4 w-4 text-green-600" />;
    } else if (change < 0) {
      return <TrendingDown className="h-4 w-4 text-red-600" />;
    } else {
      return <Minus className="h-4 w-4 text-gray-400" />;
    }
  };

  const getChangeColor = (change: number) => {
    if (change > 0) return "text-green-600";
    if (change < 0) return "text-red-600";
    return "text-gray-500";
  };

  const displayEntries = entries.slice(0, maxEntries);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center">
              <Trophy className="mr-2 h-5 w-5 text-yellow-500" />
              {title}
            </CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
          <Badge variant="outline" className="capitalize">
            {period.replace('-', ' ')}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          {displayEntries.map((entry) => (
            <div
              key={`${entry.rank}-${entry.name}`}
              className={`flex items-center justify-between p-4 rounded-lg border transition-all hover:shadow-sm ${getRankBackground(entry.rank, entry.isCurrentUser)}`}
            >
              <div className="flex items-center space-x-4">
                {/* Rank */}
                <div className="flex items-center justify-center w-8">
                  {getRankIcon(entry.rank)}
                </div>

                {/* User Info */}
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                    {entry.avatar}
                  </div>
                  <div>
                    <div className="font-medium flex items-center space-x-2">
                      <span>{entry.name}</span>
                      {entry.isCurrentUser && (
                        <Badge variant="secondary" className="text-xs">You</Badge>
                      )}
                    </div>
                    <div className="flex items-center space-x-2 text-sm">
                      <Badge className={getBadgeColor(entry.badge)}>
                        {entry.badge}
                      </Badge>
                      <div className="flex items-center space-x-1 text-muted-foreground">
                        <Zap className="h-3 w-3" />
                        <span>{entry.streak} day streak</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Points and Change */}
              <div className="text-right">
                <div className="font-bold text-lg">{entry.points.toLocaleString()}</div>
                <div className="text-sm text-muted-foreground">points</div>
                {showChange && (
                  <div className={`flex items-center justify-end space-x-1 text-sm ${getChangeColor(entry.change)}`}>
                    {getChangeIcon(entry.change)}
                    <span>
                      {entry.change > 0 ? '+' : ''}{entry.change}
                    </span>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>

        {entries.length > maxEntries && (
          <div className="mt-4 text-center">
            <Button variant="outline" size="sm">
              View Full Leaderboard
            </Button>
          </div>
        )}

        {/* Current User Position (if not in top entries) */}
        {!displayEntries.some(entry => entry.isCurrentUser) && (
          <div className="mt-4 pt-4 border-t">
            <div className="text-sm text-muted-foreground text-center mb-2">Your Position</div>
            <div className={`flex items-center justify-between p-3 rounded-lg border ${getRankBackground(0, true)}`}>
              <div className="flex items-center space-x-3">
                <div className="w-6 h-6 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-xs font-bold">
                  {Math.floor(Math.random() * 50) + 20}
                </div>
                <div>
                  <div className="font-medium">You</div>
                  <Badge className="bg-gray-100 text-gray-800 border-gray-200">
                    Go Learner
                  </Badge>
                </div>
              </div>
              <div className="text-right">
                <div className="font-bold">{Math.floor(Math.random() * 1000) + 500}</div>
                <div className="text-sm text-muted-foreground">points</div>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default Leaderboard;
