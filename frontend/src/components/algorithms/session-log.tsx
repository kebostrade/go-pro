"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Calendar,
  Clock,
  CheckCircle,
  Plus,
  Play,
  BookOpen,
  ChevronRight,
  ChevronLeft,
  MoreHorizontal,
  Trash2,
} from "lucide-react";
import { Session } from "@/types/algorithms";

interface SessionLogProps {
  sessions: Session[];
}

const SessionLog = ({ sessions }: SessionLogProps) => {
  const [selectedSession, setSelectedSession] = useState<Session | null>(null);

  const formatTime = (minutes: number) => {
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
    });
  };

  if (sessions.length === 0) {
    return (
      <div className="text-center py-12">
        <Calendar className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold mb-2">No sessions yet</h3>
        <p className="text-muted-foreground mb-4">Start your first practice session</p>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          New Session
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold">Session Log</h2>
          <p className="text-muted-foreground">{sessions.length} sessions recorded</p>
        </div>
        <Button className="flex items-center gap-2">
          <Plus className="h-4 w-4" />
          New Session
        </Button>
      </div>

      {/* Sessions Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {sessions.map((session) => {
          const isSelected = selectedSession?.id === session.id;

          return (
            <Card
              key={session.id}
              className={`cursor-pointer transition-all hover:shadow-sm ${
                isSelected ? 'ring-2 ring-primary' : ''
              }`}
              onClick={() => setSelectedSession(isSelected ? null : session)}
            >
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <Calendar className="h-5 w-5 text-muted-foreground" />
                    <div>
                      <CardTitle className="text-base">{formatDate(session.date)}</CardTitle>
                      <CardDescription>Session #{session.sessionNumber}</CardDescription>
                    </div>
                  </div>
                  <Badge variant="outline">
                    {session.topic}
                  </Badge>
                </div>
                <div className="flex items-center gap-4 text-sm text-muted-foreground mt-2">
                  <div className="flex items-center gap-1">
                    <Clock className="h-4 w-4" />
                    <span>{formatTime(session.timeSpent)}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <CheckCircle className="h-4 w-4" />
                    <span>{session.problemsSolved} solved</span>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                {session.notes && (
                  <p className="text-sm text-muted-foreground mb-2">{session.notes}</p>
                )}
                {session.problems.length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    {session.problems.slice(0, 5).map((problem, i) => (
                      <Badge key={i} variant="secondary" className="text-xs">
                        {problem}
                      </Badge>
                    ))}
                    {session.problems.length > 5 && (
                      <span className="text-xs text-muted-foreground ml-2">
                        +{session.problems.length - 5} more
                      </span>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Session Detail */}
      {selectedSession && (
        <Card className="mt-4">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>Session Details</CardTitle>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSelectedSession(null)}
              >
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <h4 className="font-medium mb-2">Problems Worked On</h4>
                <div className="flex flex-wrap gap-2">
                  {selectedSession.problems.map((problem) => (
                    <Badge key={problem} variant="outline">
                      {problem}
                    </Badge>
                  ))}
                </div>
              </div>
              <div>
                <h4 className="font-medium mb-2">Session Notes</h4>
                <p className="text-sm text-muted-foreground">
                  {selectedSession.notes}
                </p>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm">
                  <Play className="h-4 w-4 mr-1" />
                  Resume Session
                </Button>
                <Button variant="outline" size="sm">
                  <BookOpen className="h-4 w-4 mr-1" />
                  View Notes
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default SessionLog;
