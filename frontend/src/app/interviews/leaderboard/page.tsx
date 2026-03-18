'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  Trophy,
  Medal,
  Crown,
  Flame,
  Star,
  TrendingUp,
  Code2,
  MessageSquare,
  Network,
  ChevronUp,
  ChevronDown,
  Minus,
} from 'lucide-react';

interface LeaderboardEntry {
  rank: number;
  name: string;
  avatar: string;
  points: number;
  solved: number;
  streak: number;
  change: number; // positive = up, negative = down, 0 = no change
  badges: string[];
  isCurrentUser?: boolean;
}

const mockLeaderboard: LeaderboardEntry[] = [
  { rank: 1, name: 'Sarah Chen', avatar: 'SC', points: 4850, solved: 156, streak: 12, change: 0, badges: ['Legend', 'Speed Demon'] },
  { rank: 2, name: 'Alex Kumar', avatar: 'AK', points: 4720, solved: 148, streak: 8, change: 1, badges: ['On Fire'] },
  { rank: 3, name: 'Mike Johnson', avatar: 'MJ', points: 4580, solved: 142, streak: 15, change: -1, badges: ['Consistent'] },
  { rank: 4, name: 'Emma Wilson', avatar: 'EW', points: 4350, solved: 138, streak: 5, change: 2, badges: ['Rising Star'] },
  { rank: 5, name: 'David Lee', avatar: 'DL', points: 4120, solved: 132, streak: 3, change: 0, badges: [] },
  { rank: 6, name: 'Lisa Park', avatar: 'LP', points: 3980, solved: 128, streak: 7, change: 3, badges: ['Week Warrior'] },
  { rank: 7, name: 'James Brown', avatar: 'JB', points: 3850, solved: 124, streak: 2, change: -2, badges: [] },
  { rank: 8, name: 'Anna Smith', avatar: 'AS', points: 3720, solved: 120, streak: 4, change: 1, badges: [] },
  { rank: 9, name: 'Chris Davis', avatar: 'CD', points: 3580, solved: 115, streak: 6, change: -1, badges: [] },
  { rank: 10, name: 'You', avatar: 'YO', points: 760, solved: 5, streak: 5, change: 4, badges: ['First Solve'], isCurrentUser: true },
];

const badgeColors: Record<string, string> = {
  'Legend': 'bg-yellow-100 text-yellow-700',
  'Speed Demon': 'bg-blue-100 text-blue-700',
  'On Fire': 'bg-orange-100 text-orange-700',
  'Consistent': 'bg-green-100 text-green-700',
  'Rising Star': 'bg-purple-100 text-purple-700',
  'Week Warrior': 'bg-indigo-100 text-indigo-700',
  'First Solve': 'bg-pink-100 text-pink-700',
};

function getRankIcon(rank: number) {
  if (rank === 1) return <Crown className="h-6 w-6 text-yellow-500" />;
  if (rank === 2) return <Medal className="h-6 w-6 text-gray-400" />;
  if (rank === 3) return <Medal className="h-6 w-6 text-amber-600" />;
  return <span className="text-lg font-bold text-gray-500">#{rank}</span>;
}

function getChangeIcon(change: number) {
  if (change > 0) return <ChevronUp className="h-4 w-4 text-green-500" />;
  if (change < 0) return <ChevronDown className="h-4 w-4 text-red-500" />;
  return <Minus className="h-4 w-4 text-gray-400" />;
}

export default function LeaderboardPage() {
  const [timeframe, setTimeframe] = useState<'week' | 'month' | 'all'>('week');
  const [category, setCategory] = useState<'all' | 'coding' | 'behavioral' | 'system_design'>('all');

  const currentUser = mockLeaderboard.find(e => 'isCurrentUser' in e && (e as LeaderboardEntry & { isCurrentUser: boolean }).isCurrentUser);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold flex items-center gap-2">
                <Trophy className="h-6 w-6" />
                Leaderboard
              </h1>
              <p className="text-blue-100 mt-1">See how you rank against other interviewers</p>
            </div>
            <div className="flex gap-2">
              {(['week', 'month', 'all'] as const).map((tf) => (
                <button
                  key={tf}
                  onClick={() => setTimeframe(tf)}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                    timeframe === tf
                      ? 'bg-white text-blue-600'
                      : 'bg-white/20 text-white hover:bg-white/30'
                  }`}
                >
                  {tf === 'week' ? 'This Week' : tf === 'month' ? 'This Month' : 'All Time'}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="lg:col-span-1 space-y-4">
            {/* Category Filter */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="pb-2">
                <CardTitle className="text-base">Category</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                {[
                  { id: 'all', label: 'All Categories', icon: Trophy },
                  { id: 'coding', label: 'Coding', icon: Code2 },
                  { id: 'behavioral', label: 'Behavioral', icon: MessageSquare },
                  { id: 'system_design', label: 'System Design', icon: Network },
                ].map((cat) => (
                  <button
                    key={cat.id}
                    onClick={() => setCategory(cat.id as typeof category)}
                    className={`w-full flex items-center gap-3 p-3 rounded-lg transition-all ${
                      category === cat.id
                        ? 'bg-blue-50 text-blue-700 border border-blue-200'
                        : 'hover:bg-gray-50 border border-transparent'
                    }`}
                  >
                    <cat.icon className="h-4 w-4" />
                    <span className="text-sm font-medium">{cat.label}</span>
                  </button>
                ))}
              </CardContent>
            </Card>

            {/* Your Stats */}
            {currentUser && (
              <Card className="border-0 shadow-sm bg-gradient-to-br from-blue-50 to-purple-50">
                <CardHeader className="pb-2">
                  <CardTitle className="text-base flex items-center gap-2">
                    <Star className="h-4 w-4 text-yellow-500" />
                    Your Ranking
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-center">
                    <div className="text-4xl font-bold text-blue-600">#{currentUser.rank}</div>
                    <div className="text-sm text-muted-foreground mt-1">of 1,234 participants</div>
                  </div>
                  <div className="grid grid-cols-3 gap-2 mt-4 text-center">
                    <div className="p-2 bg-white rounded-lg">
                      <div className="text-lg font-bold">{currentUser.points}</div>
                      <div className="text-xs text-muted-foreground">Points</div>
                    </div>
                    <div className="p-2 bg-white rounded-lg">
                      <div className="text-lg font-bold">{currentUser.solved}</div>
                      <div className="text-xs text-muted-foreground">Solved</div>
                    </div>
                    <div className="p-2 bg-white rounded-lg">
                      <div className="text-lg font-bold flex items-center justify-center gap-1">
                        <Flame className="h-4 w-4 text-orange-500" />
                        {currentUser.streak}
                      </div>
                      <div className="text-xs text-muted-foreground">Streak</div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Top Streaks */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="pb-2">
                <CardTitle className="text-base flex items-center gap-2">
                  <Flame className="h-4 w-4 text-orange-500" />
                  Top Streaks
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {mockLeaderboard
                  .slice()
                  .sort((a, b) => b.streak - a.streak)
                  .slice(0, 5)
                  .map((entry, i) => (
                    <div key={entry.rank} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium text-muted-foreground">{i + 1}.</span>
                        <div className="w-6 h-6 rounded-full bg-gradient-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white text-xs font-medium">
                          {entry.avatar}
                        </div>
                        <span className="text-sm">{entry.name}</span>
                      </div>
                      <div className="flex items-center gap-1 text-orange-500">
                        <Flame className="h-4 w-4" />
                        <span className="font-medium">{entry.streak}</span>
                      </div>
                    </div>
                  ))}
              </CardContent>
            </Card>
          </div>

          {/* Main Leaderboard */}
          <div className="lg:col-span-3">
            {/* Top 3 */}
            <div className="grid grid-cols-3 gap-4 mb-6">
              {mockLeaderboard.slice(0, 3).map((entry, index) => (
                <Card
                  key={entry.rank}
                  className={`border-0 shadow-sm ${
                    index === 0 ? 'bg-gradient-to-br from-yellow-50 to-amber-100' :
                    index === 1 ? 'bg-gradient-to-br from-gray-50 to-gray-100' :
                    'bg-gradient-to-br from-amber-50 to-orange-100'
                  }`}
                >
                  <CardContent className="pt-6 text-center">
                    <div className="relative inline-block mb-3">
                      <div className={`w-16 h-16 rounded-full flex items-center justify-center text-white text-xl font-bold ${
                        index === 0 ? 'bg-gradient-to-br from-yellow-400 to-amber-500' :
                        index === 1 ? 'bg-gradient-to-br from-gray-300 to-gray-400' :
                        'bg-gradient-to-br from-amber-400 to-orange-500'
                      }`}>
                        {entry.avatar}
                      </div>
                      <div className="absolute -top-2 -right-2">
                        {getRankIcon(entry.rank)}
                      </div>
                    </div>
                    <h3 className="font-semibold">{entry.name}</h3>
                    <div className="text-2xl font-bold text-blue-600 mt-1">
                      {entry.points.toLocaleString()}
                    </div>
                    <div className="text-sm text-muted-foreground">points</div>
                    <div className="flex items-center justify-center gap-1 mt-2 text-orange-500">
                      <Flame className="h-4 w-4" />
                      <span className="font-medium">{entry.streak} day streak</span>
                    </div>
                    {entry.badges.length > 0 && (
                      <div className="flex flex-wrap justify-center gap-1 mt-3">
                        {entry.badges.map((badge) => (
                          <Badge key={badge} className={`text-xs ${badgeColors[badge] || 'bg-gray-100 text-gray-700'}`}>
                            {badge}
                          </Badge>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>

            {/* Rest of Leaderboard */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="border-b">
                <CardTitle className="text-base flex items-center gap-2">
                  <TrendingUp className="h-4 w-4 text-blue-500" />
                  Rankings
                </CardTitle>
              </CardHeader>
              <CardContent className="p-0">
                <div className="divide-y">
                  {mockLeaderboard.slice(3).map((entry) => {
                    const isCurrentUser = 'isCurrentUser' in entry && (entry as LeaderboardEntry & { isCurrentUser: boolean }).isCurrentUser;

                    return (
                      <div
                        key={entry.rank}
                        className={`flex items-center gap-4 p-4 ${
                          isCurrentUser ? 'bg-blue-50/50' : 'hover:bg-gray-50'
                        } transition-colors`}
                      >
                        <div className="w-10 flex justify-center">
                          {getRankIcon(entry.rank)}
                        </div>
                        <div className="flex items-center gap-3 flex-1">
                          <div className={`w-10 h-10 rounded-full flex items-center justify-center text-white font-medium ${
                            isCurrentUser ? 'bg-gradient-to-br from-blue-500 to-purple-500' : 'bg-gray-300'
                          }`}>
                            {entry.avatar}
                          </div>
                          <div>
                            <div className="flex items-center gap-2">
                              <span className={`font-medium ${isCurrentUser ? 'text-blue-700' : ''}`}>
                                {entry.name}
                              </span>
                              {isCurrentUser && (
                                <Badge variant="outline" className="text-xs border-blue-300 text-blue-600">
                                  You
                                </Badge>
                              )}
                            </div>
                            <div className="flex items-center gap-3 text-sm text-muted-foreground">
                              <span>{entry.solved} solved</span>
                              <span className="flex items-center gap-1">
                                <Flame className="h-3 w-3 text-orange-500" />
                                {entry.streak} days
                              </span>
                            </div>
                          </div>
                        </div>
                        {entry.badges.length > 0 && (
                          <div className="hidden md:flex gap-1">
                            {entry.badges.map((badge) => (
                              <Badge key={badge} className={`text-xs ${badgeColors[badge] || 'bg-gray-100 text-gray-700'}`}>
                                {badge}
                              </Badge>
                            ))}
                          </div>
                        )}
                        <div className="text-right">
                          <div className="font-bold text-lg">{entry.points.toLocaleString()}</div>
                          <div className="text-xs text-muted-foreground">points</div>
                        </div>
                        <div className="w-6 flex justify-center">
                          {getChangeIcon(entry.change)}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
