"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  CheckCircle,
  Circle,
  LoaderCircle,
  Search,
  Filter,
  ExternalLink,
  Clock,
  Calendar,
  ChevronDown,
  ChevronUp,
  Code2
} from "lucide-react";
import { Problem, Category, Difficulty, ProblemStatus, DIFFICULTY_COLORS, STATUS_CONFIG } from "@/types/algorithms";

interface ProblemsListProps {
  problems: Problem[];
  selectedCategory: string | null;
  categories: Category[];
}

const ProblemsList = ({ problems, selectedCategory, categories }: ProblemsListProps) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<ProblemStatus | "all">("all");
  const [difficultyFilter, setDifficultyFilter] = useState<Difficulty | "all">("all");
  const [sortBy, setSortBy] = useState<"number" | "difficulty" | "date">("number");
  const [expandedProblem, setExpandedProblem] = useState<string | null>(null);

  const filteredProblems = problems
    .filter((p) => {
      if (selectedCategory && p.category !== selectedCategory) return false;
      if (searchQuery && !p.title.toLowerCase().includes(searchQuery.toLowerCase())) return false;
      if (statusFilter !== "all" && p.status !== statusFilter) return false;
      if (difficultyFilter !== "all" && p.difficulty !== difficultyFilter) return false;
      return true;
    })
    .sort((a, b) => {
      switch (sortBy) {
        case "difficulty":
          const order = { Easy: 1, Medium: 2, Hard: 3 };
          return order[a.difficulty] - order[b.difficulty];
        case "date":
          return (b.dateCompleted || "").localeCompare(a.dateCompleted || "");
        default:
          return a.number - b.number;
      }
    });

  const getStatusIcon = (status: ProblemStatus) => {
    switch (status) {
      case "completed":
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case "in_progress":
        return <LoaderCircle className="h-5 w-5 text-blue-500 animate-spin" />;
      default:
        return <Circle className="h-5 w-5 text-gray-300" />;
    }
  };

  const formatTime = (minutes: number) => {
    if (minutes === 0) return "-";
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
  };

  const getCategoryName = (categoryId: string) => {
    return categories.find((c) => c.id === categoryId)?.name || categoryId;
  };

  const completedCount = filteredProblems.filter((p) => p.status === "completed").length;

  return (
    <Card>
      <CardHeader>
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <CardTitle className="flex items-center gap-2">
              <Code2 className="h-5 w-5" />
              Problems
              {selectedCategory && (
                <Badge variant="secondary" className="ml-2">
                  {getCategoryName(selectedCategory)}
                </Badge>
              )}
            </CardTitle>
            <CardDescription>
              {filteredProblems.length} problems • {completedCount} completed
            </CardDescription>
          </div>

          {/* Filters */}
          <div className="flex flex-wrap gap-2">
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search problems..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-8 w-48"
              />
            </div>

            <Select value={statusFilter} onValueChange={(v: string) => setStatusFilter(v as ProblemStatus | "all")}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="pending">Not Started</SelectItem>
                <SelectItem value="in_progress">In Progress</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
              </SelectContent>
            </Select>

            <Select value={difficultyFilter} onValueChange={(v: string) => setDifficultyFilter(v as Difficulty | "all")}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Difficulty" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Levels</SelectItem>
                <SelectItem value="Easy">Easy</SelectItem>
                <SelectItem value="Medium">Medium</SelectItem>
                <SelectItem value="Hard">Hard</SelectItem>
              </SelectContent>
            </Select>

            <Select value={sortBy} onValueChange={(v: string) => setSortBy(v as "number" | "difficulty" | "date")}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="number">Number</SelectItem>
                <SelectItem value="difficulty">Difficulty</SelectItem>
                <SelectItem value="date">Date</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        <div className="space-y-2">
          {filteredProblems.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              <Filter className="h-12 w-12 mx-auto mb-3 opacity-50" />
              <p>No problems match your filters</p>
            </div>
          ) : (
            filteredProblems.map((problem) => (
              <div
                key={problem.id}
                className={`border rounded-lg transition-all ${
                  expandedProblem === problem.id ? 'ring-2 ring-primary' : 'hover:border-gray-300'
                }`}
              >
                <div
                  className="flex items-center justify-between p-4 cursor-pointer"
                  onClick={() => setExpandedProblem(expandedProblem === problem.id ? null : problem.id)}
                >
                  <div className="flex items-center gap-4 flex-1 min-w-0">
                    {getStatusIcon(problem.status)}

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-medium truncate">{problem.title}</span>
                        {problem.leetcodeNumber && (
                          <Badge variant="outline" className="text-xs">
                            #{problem.leetcodeNumber}
                          </Badge>
                        )}
                      </div>
                      <div className="flex items-center gap-3 mt-1">
                        <Badge className={DIFFICULTY_COLORS[problem.difficulty]}>
                          {problem.difficulty}
                        </Badge>
                        {problem.status === "completed" && (
                          <>
                            <span className="flex items-center gap-1 text-xs text-muted-foreground">
                              <Clock className="h-3 w-3" />
                              {formatTime(problem.timeSpent)}
                            </span>
                            {problem.dateCompleted && (
                              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                                <Calendar className="h-3 w-3" />
                                {problem.dateCompleted}
                              </span>
                            )}
                          </>
                        )}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    {problem.patterns.length > 0 && (
                      <div className="hidden md:flex gap-1">
                        {problem.patterns.slice(0, 2).map((pattern) => (
                          <Badge key={pattern} variant="secondary" className="text-xs">
                            {pattern}
                          </Badge>
                        ))}
                        {problem.patterns.length > 2 && (
                          <Badge variant="secondary" className="text-xs">
                            +{problem.patterns.length - 2}
                          </Badge>
                        )}
                      </div>
                    )}
                    {expandedProblem === problem.id ? (
                      <ChevronUp className="h-5 w-5 text-muted-foreground" />
                    ) : (
                      <ChevronDown className="h-5 w-5 text-muted-foreground" />
                    )}
                  </div>
                </div>

                {expandedProblem === problem.id && (
                  <div className="border-t px-4 py-3 bg-muted/30">
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-3">
                      <div>
                        <p className="text-xs text-muted-foreground">Category</p>
                        <p className="font-medium text-sm">{getCategoryName(problem.category)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Attempts</p>
                        <p className="font-medium text-sm">{problem.attempts}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Reviews</p>
                        <p className="font-medium text-sm">{problem.reviewCount}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Next Review</p>
                        <p className="font-medium text-sm">{problem.nextReview || "Not scheduled"}</p>
                      </div>
                    </div>

                    {problem.patterns.length > 0 && (
                      <div className="mb-3">
                        <p className="text-xs text-muted-foreground mb-1">Patterns</p>
                        <div className="flex flex-wrap gap-1">
                          {problem.patterns.map((pattern) => (
                            <Badge key={pattern} variant="outline" className="text-xs">
                              {pattern}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}

                    {problem.notes && (
                      <div className="mb-3">
                        <p className="text-xs text-muted-foreground mb-1">Notes</p>
                        <p className="text-sm bg-background p-2 rounded border">{problem.notes}</p>
                      </div>
                    )}

                    <div className="flex gap-2">
                      <Button size="sm" variant="outline">
                        <ExternalLink className="h-4 w-4 mr-1" />
                        Open Problem
                      </Button>
                      <Button size="sm" variant="outline">
                        <Code2 className="h-4 w-4 mr-1" />
                        View Solution
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default ProblemsList;
