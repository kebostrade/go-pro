"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { 
  TrendingUp, 
  Clock, 
  Target, 
  Award,
  Calendar,
  Zap,
  BookOpen,
  Code2,
  Trophy
} from "lucide-react";

interface AnalyticsData {
  totalHours: number;
  streak: number;
  completedCourses: number;
  skillLevel: string;
  weeklyProgress: Array<{
    day: string;
    hours: number;
  }>;
  skillBreakdown: Array<{
    skill: string;
    level: number;
    maxLevel: number;
  }>;
  achievements: Array<{
    id: string;
    title: string;
    description: string;
    icon: string;
    earned: boolean;
    earnedDate?: string;
  }>;
}

interface LearningAnalyticsProps {
  data: AnalyticsData;
}

const LearningAnalytics = ({ data }: LearningAnalyticsProps) => {
  const getSkillColor = (level: number, maxLevel: number) => {
    const percentage = (level / maxLevel) * 100;
    if (percentage >= 80) return 'text-green-500';
    if (percentage >= 60) return 'text-blue-500';
    if (percentage >= 40) return 'text-yellow-500';
    return 'text-gray-500';
  };

  const getAchievementIcon = (iconName: string) => {
    const icons: { [key: string]: any } = {
      'first-lesson': BookOpen,
      'streak-7': Calendar,
      'first-project': Trophy,
      'speed-demon': Zap,
      'code-master': Code2,
    };
    return icons[iconName] || Award;
  };

  return (
    <div className="space-y-6">
      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Hours</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.totalHours}h</div>
            <p className="text-xs text-muted-foreground">
              +2.5h from last week
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Current Streak</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.streak} days</div>
            <p className="text-xs text-muted-foreground">
              Keep it up! 🔥
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Completed</CardTitle>
            <Trophy className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.completedCourses}</div>
            <p className="text-xs text-muted-foreground">
              courses finished
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Skill Level</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.skillLevel}</div>
            <p className="text-xs text-muted-foreground">
              Go Developer
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Weekly Progress */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Calendar className="mr-2 h-5 w-5" />
            Weekly Activity
          </CardTitle>
          <CardDescription>
            Your learning activity over the past 7 days
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {data.weeklyProgress.map((day, index) => (
              <div key={index} className="flex items-center space-x-4">
                <div className="w-12 text-sm text-muted-foreground">
                  {day.day}
                </div>
                <div className="flex-1">
                  <Progress 
                    value={(day.hours / 4) * 100} 
                    className="h-2" 
                  />
                </div>
                <div className="w-16 text-sm font-medium text-right">
                  {day.hours}h
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Skills Breakdown */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Target className="mr-2 h-5 w-5" />
            Skill Progress
          </CardTitle>
          <CardDescription>
            Your proficiency in different Go concepts
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {data.skillBreakdown.map((skill, index) => (
              <div key={index} className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">{skill.skill}</span>
                  <div className="flex items-center space-x-2">
                    <span className={`text-sm font-bold ${getSkillColor(skill.level, skill.maxLevel)}`}>
                      {skill.level}/{skill.maxLevel}
                    </span>
                    <Badge variant="outline" className="text-xs">
                      {Math.round((skill.level / skill.maxLevel) * 100)}%
                    </Badge>
                  </div>
                </div>
                <Progress 
                  value={(skill.level / skill.maxLevel) * 100} 
                  className="h-2" 
                />
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Achievements */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Award className="mr-2 h-5 w-5" />
            Achievements
          </CardTitle>
          <CardDescription>
            Your learning milestones and accomplishments
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {data.achievements.map((achievement) => {
              const IconComponent = getAchievementIcon(achievement.icon);
              return (
                <div
                  key={achievement.id}
                  className={`flex items-center space-x-3 p-3 rounded-lg border transition-colors ${
                    achievement.earned
                      ? 'bg-primary/5 border-primary/20'
                      : 'bg-muted/50 border-border opacity-60'
                  }`}
                >
                  <div className={`flex h-10 w-10 items-center justify-center rounded-lg ${
                    achievement.earned
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-muted text-muted-foreground'
                  }`}>
                    <IconComponent className="h-5 w-5" />
                  </div>
                  <div className="flex-1">
                    <h4 className={`font-medium ${!achievement.earned ? 'text-muted-foreground' : ''}`}>
                      {achievement.title}
                    </h4>
                    <p className="text-sm text-muted-foreground">
                      {achievement.description}
                    </p>
                    {achievement.earned && achievement.earnedDate && (
                      <p className="text-xs text-primary font-medium mt-1">
                        Earned {achievement.earnedDate}
                      </p>
                    )}
                  </div>
                  {achievement.earned && (
                    <Badge className="go-gradient text-white">
                      Earned
                    </Badge>
                  )}
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default LearningAnalytics;
