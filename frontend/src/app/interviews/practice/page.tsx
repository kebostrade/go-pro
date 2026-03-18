'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { Input } from '@/components/ui/input';
import {
  Code2,
  MessageSquare,
  Network,
  Clock,
  Trophy,
  Flame,
  CheckCircle,
  Circle,
  Play,
  Search,
  Filter,
  ChevronRight,
  Zap,
  Target,
  Star,
  TrendingUp,
  BarChart3,
  Award,
} from 'lucide-react';
import Link from 'next/link';
import type { InterviewType, Difficulty } from '@/types/interview';

// Problem sets like HackerRank
const problems = [
  // Coding - Easy
  { id: 'c1', title: 'Sum of Even Numbers', type: 'coding' as InterviewType, difficulty: 'beginner' as Difficulty, category: 'Arrays', status: 'solved', score: 100, time: '5:23' },
  { id: 'c2', title: 'String Reversal', type: 'coding' as InterviewType, difficulty: 'beginner' as Difficulty, category: 'Strings', status: 'solved', score: 95, time: '4:15' },
  { id: 'c3', title: 'Find Maximum Value', type: 'coding' as InterviewType, difficulty: 'beginner' as Difficulty, category: 'Arrays', status: 'solved', score: 100, time: '3:45' },
  { id: 'c4', title: 'Merge Sorted Slices', type: 'coding' as InterviewType, difficulty: 'intermediate' as Difficulty, category: 'Algorithms', status: 'attempted', score: 60, time: '12:30' },
  { id: 'c5', title: 'Detect Linked List Cycle', type: 'coding' as InterviewType, difficulty: 'intermediate' as Difficulty, category: 'Data Structures', status: 'unsolved', score: 0, time: '-' },
  { id: 'c6', title: 'Thread-Safe LRU Cache', type: 'coding' as InterviewType, difficulty: 'advanced' as Difficulty, category: 'System Design', status: 'unsolved', score: 0, time: '-' },
  // Behavioral
  { id: 'b1', title: 'Working with Difficult Team Members', type: 'behavioral' as InterviewType, difficulty: 'beginner' as Difficulty, category: 'Teamwork', status: 'solved', score: 85, time: '4:50' },
  { id: 'b2', title: 'Taking Initiative', type: 'behavioral' as InterviewType, difficulty: 'beginner' as Difficulty, category: 'Leadership', status: 'solved', score: 90, time: '5:10' },
  { id: 'b3', title: 'Technical Decision Making', type: 'behavioral' as InterviewType, difficulty: 'intermediate' as Difficulty, category: 'Problem Solving', status: 'unsolved', score: 0, time: '-' },
  // System Design
  { id: 's1', title: 'URL Shortener Service', type: 'system_design' as InterviewType, difficulty: 'beginner' as Difficulty, category: 'Web Services', status: 'attempted', score: 70, time: '15:00' },
  { id: 's2', title: 'Real-time Chat Application', type: 'system_design' as InterviewType, difficulty: 'intermediate' as Difficulty, category: 'Distributed Systems', status: 'unsolved', score: 0, time: '-' },
  { id: 's3', title: 'Distributed Key-Value Store', type: 'system_design' as InterviewType, difficulty: 'advanced' as Difficulty, category: 'Databases', status: 'unsolved', score: 0, time: '-' },
];

const categories = {
  coding: ['Arrays', 'Strings', 'Algorithms', 'Data Structures', 'System Design'],
  behavioral: ['Teamwork', 'Leadership', 'Problem Solving', 'Communication'],
  system_design: ['Web Services', 'Distributed Systems', 'Databases', 'Caching'],
};

const tracks = [
  { id: 'coding', name: 'Coding Interview', icon: Code2, color: 'text-blue-500', bg: 'bg-blue-500/10', count: 6, solved: 3 },
  { id: 'behavioral', name: 'Behavioral Interview', icon: MessageSquare, color: 'text-purple-500', bg: 'bg-purple-500/10', count: 3, solved: 2 },
  { id: 'system_design', name: 'System Design', icon: Network, color: 'text-orange-500', bg: 'bg-orange-500/10', count: 3, solved: 0 },
];

const difficultyConfig = {
  beginner: { label: 'Easy', color: 'bg-green-100 text-green-700 border-green-200', dot: 'bg-green-500' },
  intermediate: { label: 'Medium', color: 'bg-yellow-100 text-yellow-700 border-yellow-200', dot: 'bg-yellow-500' },
  advanced: { label: 'Hard', color: 'bg-red-100 text-red-700 border-red-200', dot: 'bg-red-500' },
};

export default function InterviewPracticePage() {
  const [selectedType, setSelectedType] = useState<InterviewType | 'all'>('all');
  const [selectedDifficulty, setSelectedDifficulty] = useState<Difficulty | 'all'>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [isStarting, setIsStarting] = useState(false);

  const filteredProblems = problems.filter(p => {
    if (selectedType !== 'all' && p.type !== selectedType) return false;
    if (selectedDifficulty !== 'all' && p.difficulty !== selectedDifficulty) return false;
    if (searchQuery && !p.title.toLowerCase().includes(searchQuery.toLowerCase())) return false;
    return true;
  });

  const stats = {
    solved: problems.filter(p => p.status === 'solved').length,
    attempted: problems.filter(p => p.status === 'attempted').length,
    total: problems.length,
    streak: 5,
    points: problems.reduce((acc, p) => acc + p.score, 0),
  };

  const handleStartInterview = (type: InterviewType, difficulty: Difficulty) => {
    setIsStarting(true);
    setTimeout(() => {
      window.location.href = `/interviews/session?type=${type}&difficulty=${difficulty}`;
    }, 300);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Navigation Bar */}
      <div className="bg-white border-b sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-6">
            <Link href="/interviews" className="flex items-center gap-2">
              <Code2 className="h-6 w-6 text-blue-600" />
              <span className="font-bold text-lg">Go-Pro Mock Interviews</span>
            </Link>
            <nav className="hidden md:flex items-center gap-4">
              <Link href="/interviews/practice" className="text-blue-600 font-medium">Practice</Link>
              <Link href="/interviews/history" className="text-gray-600 hover:text-gray-900">History</Link>
              <Link href="/interviews/leaderboard" className="text-gray-600 hover:text-gray-900">Leaderboard</Link>
            </nav>
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 px-3 py-1.5 bg-orange-50 rounded-full">
              <Flame className="h-4 w-4 text-orange-500" />
              <span className="font-semibold text-orange-600">{stats.streak} day streak</span>
            </div>
            <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-50 rounded-full">
              <Star className="h-4 w-4 text-blue-500" />
              <span className="font-semibold text-blue-600">{stats.points} pts</span>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="lg:col-span-1 space-y-4">
            {/* Quick Start */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="pb-2">
                <CardTitle className="text-base flex items-center gap-2">
                  <Play className="h-4 w-4 text-green-500" />
                  Quick Start
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button
                  className="w-full bg-green-600 hover:bg-green-700"
                  onClick={() => handleStartInterview('coding', 'beginner')}
                  disabled={isStarting}
                >
                  <Play className="h-4 w-4 mr-2" />
                  Start Coding Interview
                </Button>
                <div className="grid grid-cols-2 gap-2">
                  <Button variant="outline" size="sm" onClick={() => handleStartInterview('behavioral', 'beginner')}>
                    <MessageSquare className="h-3 w-3 mr-1" />
                    Behavioral
                  </Button>
                  <Button variant="outline" size="sm" onClick={() => handleStartInterview('system_design', 'beginner')}>
                    <Network className="h-3 w-3 mr-1" />
                    System Design
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Progress Stats */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="pb-2">
                <CardTitle className="text-base flex items-center gap-2">
                  <BarChart3 className="h-4 w-4 text-blue-500" />
                  Your Progress
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <div className="flex justify-between text-sm mb-1">
                      <span className="text-gray-600">Problems Solved</span>
                      <span className="font-semibold">{stats.solved}/{stats.total}</span>
                    </div>
                    <Progress value={(stats.solved / stats.total) * 100} className="h-2" />
                  </div>

                  <div className="grid grid-cols-3 gap-2 text-center">
                    <div className="p-2 bg-green-50 rounded-lg">
                      <div className="text-lg font-bold text-green-600">{stats.solved}</div>
                      <div className="text-xs text-gray-500">Solved</div>
                    </div>
                    <div className="p-2 bg-yellow-50 rounded-lg">
                      <div className="text-lg font-bold text-yellow-600">{stats.attempted}</div>
                      <div className="text-xs text-gray-500">Attempted</div>
                    </div>
                    <div className="p-2 bg-gray-50 rounded-lg">
                      <div className="text-lg font-bold text-gray-600">{stats.total - stats.solved - stats.attempted}</div>
                      <div className="text-xs text-gray-500">Remaining</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Tracks */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="pb-2">
                <CardTitle className="text-base flex items-center gap-2">
                  <Target className="h-4 w-4 text-purple-500" />
                  Interview Tracks
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                {tracks.map((track) => (
                  <button
                    key={track.id}
                    onClick={() => setSelectedType(selectedType === track.id ? 'all' : track.id as InterviewType)}
                    className={`w-full flex items-center justify-between p-3 rounded-lg transition-all ${
                      selectedType === track.id
                        ? `${track.bg} border-2 border-current ${track.color}`
                        : 'hover:bg-gray-50 border border-gray-200'
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      <div className={`p-2 rounded-lg ${track.bg}`}>
                        <track.icon className={`h-4 w-4 ${track.color}`} />
                      </div>
                      <div className="text-left">
                        <div className="font-medium text-sm">{track.name}</div>
                        <div className="text-xs text-gray-500">{track.solved}/{track.count} solved</div>
                      </div>
                    </div>
                    <ChevronRight className="h-4 w-4 text-gray-400" />
                  </button>
                ))}
              </CardContent>
            </Card>

            {/* Badges */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="pb-2">
                <CardTitle className="text-base flex items-center gap-2">
                  <Award className="h-4 w-4 text-yellow-500" />
                  Badges
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  <div className="flex items-center gap-1 px-2 py-1 bg-yellow-100 rounded-full">
                    <Trophy className="h-3 w-3 text-yellow-600" />
                    <span className="text-xs font-medium text-yellow-700">First Solve</span>
                  </div>
                  <div className="flex items-center gap-1 px-2 py-1 bg-blue-100 rounded-full">
                    <Zap className="h-3 w-3 text-blue-600" />
                    <span className="text-xs font-medium text-blue-700">Speed Demon</span>
                  </div>
                  <div className="flex items-center gap-1 px-2 py-1 bg-purple-100 rounded-full">
                    <Flame className="h-3 w-3 text-purple-600" />
                    <span className="text-xs font-medium text-purple-700">On Fire</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Content */}
          <div className="lg:col-span-3 space-y-4">
            {/* Filters */}
            <Card className="border-0 shadow-sm">
              <CardContent className="py-4">
                <div className="flex flex-col md:flex-row gap-4">
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                    <Input
                      placeholder="Search problems..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-10"
                    />
                  </div>
                  <div className="flex gap-2">
                    <select
                      value={selectedDifficulty}
                      onChange={(e) => setSelectedDifficulty(e.target.value as Difficulty | 'all')}
                      className="px-4 py-2 border rounded-lg text-sm bg-white"
                    >
                      <option value="all">All Difficulties</option>
                      <option value="beginner">Easy</option>
                      <option value="intermediate">Medium</option>
                      <option value="advanced">Hard</option>
                    </select>
                    <select
                      value={selectedType}
                      onChange={(e) => setSelectedType(e.target.value as InterviewType | 'all')}
                      className="px-4 py-2 border rounded-lg text-sm bg-white"
                    >
                      <option value="all">All Types</option>
                      <option value="coding">Coding</option>
                      <option value="behavioral">Behavioral</option>
                      <option value="system_design">System Design</option>
                    </select>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Problems List */}
            <Card className="border-0 shadow-sm">
              <CardHeader className="border-b bg-gray-50/50">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-base">Practice Problems</CardTitle>
                  <span className="text-sm text-gray-500">{filteredProblems.length} problems</span>
                </div>
              </CardHeader>
              <CardContent className="p-0">
                <div className="divide-y">
                  {filteredProblems.map((problem) => (
                    <div
                      key={problem.id}
                      className="flex items-center gap-4 p-4 hover:bg-gray-50 cursor-pointer transition-colors"
                      onClick={() => window.location.href = `/interviews/session?type=${problem.type}&difficulty=${problem.difficulty}`}
                    >
                      {/* Status Icon */}
                      <div className="w-8 h-8 flex items-center justify-center">
                        {problem.status === 'solved' ? (
                          <CheckCircle className="h-5 w-5 text-green-500" />
                        ) : problem.status === 'attempted' ? (
                          <div className="h-5 w-5 rounded-full border-2 border-yellow-500 flex items-center justify-center">
                            <div className="h-2 w-2 rounded-full bg-yellow-500" />
                          </div>
                        ) : (
                          <Circle className="h-5 w-5 text-gray-300" />
                        )}
                      </div>

                      {/* Problem Info */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <h3 className="font-medium text-gray-900 truncate">{problem.title}</h3>
                          {problem.status === 'solved' && (
                            <Badge variant="secondary" className="text-xs bg-green-100 text-green-700">
                              {problem.score}%
                            </Badge>
                          )}
                        </div>
                        <div className="flex items-center gap-3 mt-1">
                          <Badge
                            variant="outline"
                            className={`text-xs ${difficultyConfig[problem.difficulty].color}`}
                          >
                            <span className={`w-1.5 h-1.5 rounded-full ${difficultyConfig[problem.difficulty].dot} mr-1.5`} />
                            {difficultyConfig[problem.difficulty].label}
                          </Badge>
                          <span className="text-xs text-gray-500">{problem.category}</span>
                          {problem.time !== '-' && (
                            <span className="text-xs text-gray-500 flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              {problem.time}
                            </span>
                          )}
                        </div>
                      </div>

                      {/* Type Badge */}
                      <Badge
                        variant="outline"
                        className={`hidden sm:flex text-xs ${
                          problem.type === 'coding' ? 'border-blue-200 text-blue-700' :
                          problem.type === 'behavioral' ? 'border-purple-200 text-purple-700' :
                          'border-orange-200 text-orange-700'
                        }`}
                      >
                        {problem.type === 'coding' && <Code2 className="h-3 w-3 mr-1" />}
                        {problem.type === 'behavioral' && <MessageSquare className="h-3 w-3 mr-1" />}
                        {problem.type === 'system_design' && <Network className="h-3 w-3 mr-1" />}
                        {problem.type.replace('_', ' ')}
                      </Badge>

                      {/* Arrow */}
                      <ChevronRight className="h-5 w-5 text-gray-400" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Skill Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Card className="border-0 shadow-sm bg-gradient-to-br from-blue-50 to-blue-100/50">
                <CardContent className="p-4">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="p-2 bg-blue-500 rounded-lg">
                      <Code2 className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-semibold">Coding Skills</h3>
                      <p className="text-xs text-gray-500">Algorithms & Data Structures</p>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Arrays & Strings</span>
                      <span className="text-blue-600 font-medium">2/3</span>
                    </div>
                    <Progress value={66} className="h-1.5" />
                  </div>
                </CardContent>
              </Card>

              <Card className="border-0 shadow-sm bg-gradient-to-br from-purple-50 to-purple-100/50">
                <CardContent className="p-4">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="p-2 bg-purple-500 rounded-lg">
                      <MessageSquare className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-semibold">Soft Skills</h3>
                      <p className="text-xs text-gray-500">Communication & Leadership</p>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Behavioral Questions</span>
                      <span className="text-purple-600 font-medium">2/3</span>
                    </div>
                    <Progress value={66} className="h-1.5" />
                  </div>
                </CardContent>
              </Card>

              <Card className="border-0 shadow-sm bg-gradient-to-br from-orange-50 to-orange-100/50">
                <CardContent className="p-4">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="p-2 bg-orange-500 rounded-lg">
                      <Network className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="font-semibold">System Design</h3>
                      <p className="text-xs text-gray-500">Architecture & Scalability</p>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Design Patterns</span>
                      <span className="text-orange-600 font-medium">0/3</span>
                    </div>
                    <Progress value={0} className="h-1.5" />
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
