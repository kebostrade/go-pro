'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Code2,
  MessageSquare,
  Network,
  Clock,
  Trophy,
  ChevronRight,
  Calendar,
  Filter,
  TrendingUp,
  BarChart3,
  RefreshCw,
} from 'lucide-react';
import Link from 'next/link';
import { api } from '@/lib/api';

interface SessionListItem {
  id: string;
  user_id?: string;
  type: 'coding' | 'behavioral' | 'system_design';
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  status: 'in_progress' | 'completed' | 'abandoned';
  score?: number;
  created_at: string;
  completed_at?: string;
}

const typeConfig = {
  coding: { icon: Code2, color: 'text-blue-500', bg: 'bg-blue-500/10', label: 'Coding' },
  behavioral: { icon: MessageSquare, color: 'text-purple-500', bg: 'bg-purple-500/10', label: 'Behavioral' },
  system_design: { icon: Network, color: 'text-orange-500', bg: 'bg-orange-500/10', label: 'System Design' },
};

const difficultyConfig = {
  beginner: { label: 'Easy', color: 'bg-green-100 text-green-700' },
  intermediate: { label: 'Medium', color: 'bg-yellow-100 text-yellow-700' },
  advanced: { label: 'Hard', color: 'bg-red-100 text-red-700' },
};

// Mock sessions for demo mode
const mockSessions: SessionListItem[] = [
  {
    id: 'session-1',
    type: 'coding',
    difficulty: 'beginner',
    status: 'completed',
    score: 85,
    created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(), // 2 hours ago
    completed_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 'session-2',
    type: 'behavioral',
    difficulty: 'intermediate',
    status: 'completed',
    score: 92,
    created_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(), // 1 day ago
    completed_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 'session-3',
    type: 'system_design',
    difficulty: 'beginner',
    status: 'completed',
    score: 78,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(), // 2 days ago
    completed_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 'session-4',
    type: 'coding',
    difficulty: 'intermediate',
    status: 'completed',
    score: 70,
    created_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(), // 3 days ago
    completed_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
  },
  {
    id: 'session-5',
    type: 'coding',
    difficulty: 'advanced',
    status: 'abandoned',
    score: 0,
    created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(), // 5 days ago
  },
];

function formatTimeAgo(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 60) return `${diffMins}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;
  return date.toLocaleDateString();
}

export default function HistoryPage() {
  const [sessions, setSessions] = useState<SessionListItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filterType, setFilterType] = useState<string>('all');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  useEffect(() => {
    const fetchSessions = async () => {
      setIsLoading(true);
      try {
        const data = await api.getInterviewSessions();
        if (data && data.length > 0) {
          setSessions(data as SessionListItem[]);
        } else {
          setSessions(mockSessions);
        }
      } catch (err) {
        console.log('Using mock sessions:', err);
        setSessions(mockSessions);
      } finally {
        setIsLoading(false);
      }
    };

    fetchSessions();
  }, []);

  const filteredSessions = sessions.filter(s => {
    if (filterType !== 'all' && s.type !== filterType) return false;
    if (filterStatus !== 'all' && s.status !== filterStatus) return false;
    return true;
  });

  const stats = {
    total: sessions.length,
    completed: sessions.filter(s => s.status === 'completed').length,
    avgScore: Math.round(
      sessions.filter(s => s.status === 'completed').reduce((acc, s) => acc + (s.score || 0), 0) /
      Math.max(1, sessions.filter(s => s.status === 'completed').length)
    ),
    thisWeek: sessions.filter(s => {
      const date = new Date(s.created_at);
      const weekAgo = new Date();
      weekAgo.setDate(weekAgo.getDate() - 7);
      return date > weekAgo;
    }).length,
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold flex items-center gap-2">
                <Calendar className="h-6 w-6 text-blue-500" />
                Interview History
              </h1>
              <p className="text-muted-foreground mt-1">Track your progress over time</p>
            </div>
            <Link href="/interviews/practice">
              <Button>
                <RefreshCw className="h-4 w-4 mr-2" />
                New Interview
              </Button>
            </Link>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-6">
        {/* Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <Card className="border-0 shadow-sm">
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-blue-500/10">
                  <BarChart3 className="h-5 w-5 text-blue-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{stats.total}</div>
                  <div className="text-sm text-muted-foreground">Total Sessions</div>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card className="border-0 shadow-sm">
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-green-500/10">
                  <Trophy className="h-5 w-5 text-green-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{stats.completed}</div>
                  <div className="text-sm text-muted-foreground">Completed</div>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card className="border-0 shadow-sm">
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-purple-500/10">
                  <TrendingUp className="h-5 w-5 text-purple-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{stats.avgScore}%</div>
                  <div className="text-sm text-muted-foreground">Avg Score</div>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card className="border-0 shadow-sm">
            <CardContent className="p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-orange-500/10">
                  <Clock className="h-5 w-5 text-orange-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{stats.thisWeek}</div>
                  <div className="text-sm text-muted-foreground">This Week</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Filters */}
        <Card className="border-0 shadow-sm mb-6">
          <CardContent className="py-4">
            <div className="flex items-center gap-4">
              <Filter className="h-4 w-4 text-gray-400" />
              <select
                value={filterType}
                onChange={(e) => setFilterType(e.target.value)}
                className="px-4 py-2 border rounded-lg text-sm bg-white"
              >
                <option value="all">All Types</option>
                <option value="coding">Coding</option>
                <option value="behavioral">Behavioral</option>
                <option value="system_design">System Design</option>
              </select>
              <select
                value={filterStatus}
                onChange={(e) => setFilterStatus(e.target.value)}
                className="px-4 py-2 border rounded-lg text-sm bg-white"
              >
                <option value="all">All Status</option>
                <option value="completed">Completed</option>
                <option value="abandoned">Abandoned</option>
              </select>
            </div>
          </CardContent>
        </Card>

        {/* Sessions List */}
        <Card className="border-0 shadow-sm">
          <CardHeader className="border-b">
            <CardTitle className="text-base">Sessions</CardTitle>
          </CardHeader>
          <CardContent className="p-0">
            {isLoading ? (
              <div className="p-8 text-center text-muted-foreground">
                Loading sessions...
              </div>
            ) : filteredSessions.length === 0 ? (
              <div className="p-8 text-center text-muted-foreground">
                No sessions found. Start a practice interview to see your history.
              </div>
            ) : (
              <div className="divide-y">
                {filteredSessions.map((session) => {
                  const config = typeConfig[session.type];
                  const diffConfig = difficultyConfig[session.difficulty];
                  const Icon = config.icon;

                  return (
                    <Link
                      key={session.id}
                      href={`/interviews/feedback?type=${session.type}&difficulty=${session.difficulty}`}
                      className="flex items-center gap-4 p-4 hover:bg-gray-50 transition-colors"
                    >
                      <div className={`p-3 rounded-xl ${config.bg}`}>
                        <Icon className={`h-6 w-6 ${config.color}`} />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{config.label} Interview</span>
                          <Badge variant="outline" className={diffConfig.color}>
                            {diffConfig.label}
                          </Badge>
                          <Badge
                            variant={session.status === 'completed' ? 'default' : 'secondary'}
                            className={session.status === 'completed' ? 'bg-green-100 text-green-700' : ''}
                          >
                            {session.status}
                          </Badge>
                        </div>
                        <div className="flex items-center gap-4 mt-1 text-sm text-muted-foreground">
                          <span className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            {formatTimeAgo(session.created_at)}
                          </span>
                          {session.status === 'completed' && (
                            <span className="flex items-center gap-1">
                              <Trophy className="h-3 w-3" />
                              Score: {session.score}%
                            </span>
                          )}
                        </div>
                      </div>
                      <ChevronRight className="h-5 w-5 text-gray-400" />
                    </Link>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
