"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { 
  Trophy, 
  Star, 
  Award, 
  Target, 
  Zap, 
  Calendar,
  BookOpen,
  Code2,
  Users,
  TrendingUp,
  Lock,
  CheckCircle,
  Sparkles
} from "lucide-react";

interface Achievement {
  id: string;
  title: string;
  description: string;
  category: 'learning' | 'streak' | 'social' | 'mastery' | 'special';
  tier: 'bronze' | 'silver' | 'gold' | 'platinum';
  icon: string;
  earned: boolean;
  earnedDate?: string;
  progress?: number;
  maxProgress?: number;
  rarity?: 'common' | 'rare' | 'epic' | 'legendary';
  points: number;
}

interface AchievementBadgeProps {
  achievement: Achievement;
  size?: 'sm' | 'md' | 'lg';
  showProgress?: boolean;
  onClick?: () => void;
}

const AchievementBadge = ({ 
  achievement, 
  size = 'md', 
  showProgress = true,
  onClick 
}: AchievementBadgeProps) => {
  const [isHovered, setIsHovered] = useState(false);

  const getIcon = (iconName: string) => {
    const icons: { [key: string]: any } = {
      'trophy': Trophy,
      'star': Star,
      'award': Award,
      'target': Target,
      'zap': Zap,
      'calendar': Calendar,
      'book': BookOpen,
      'code': Code2,
      'users': Users,
      'trending': TrendingUp,
    };
    return icons[iconName] || Award;
  };

  const getTierColor = (tier: string) => {
    switch (tier) {
      case 'bronze': return 'from-amber-600 to-amber-800';
      case 'silver': return 'from-gray-400 to-gray-600';
      case 'gold': return 'from-yellow-400 to-yellow-600';
      case 'platinum': return 'from-purple-400 to-purple-600';
      default: return 'from-gray-400 to-gray-600';
    }
  };

  const getRarityColor = (rarity: string) => {
    switch (rarity) {
      case 'common': return 'border-gray-300 dark:border-gray-700';
      case 'rare': return 'border-blue-400 dark:border-blue-600';
      case 'epic': return 'border-purple-400 dark:border-purple-600';
      case 'legendary': return 'border-yellow-400 dark:border-yellow-600';
      default: return 'border-gray-300 dark:border-gray-700';
    }
  };

  const getSizeClasses = (size: string) => {
    switch (size) {
      case 'sm': return {
        card: 'p-3',
        icon: 'h-8 w-8',
        iconContainer: 'h-12 w-12',
        title: 'text-sm',
        description: 'text-xs',
      };
      case 'lg': return {
        card: 'p-6',
        icon: 'h-12 w-12',
        iconContainer: 'h-20 w-20',
        title: 'text-xl',
        description: 'text-base',
      };
      default: return {
        card: 'p-4',
        icon: 'h-10 w-10',
        iconContainer: 'h-16 w-16',
        title: 'text-base',
        description: 'text-sm',
      };
    }
  };

  const IconComponent = getIcon(achievement.icon);
  const sizeClasses = getSizeClasses(size);
  const progressPercentage = achievement.progress && achievement.maxProgress 
    ? (achievement.progress / achievement.maxProgress) * 100 
    : 0;

  return (
    <Card 
      className={`
        relative transition-all duration-300 cursor-pointer
        ${achievement.earned 
          ? `${getRarityColor(achievement.rarity || 'common')} border-2 shadow-lg` 
          : 'opacity-60 border-dashed'
        }
        ${isHovered && achievement.earned ? 'scale-105 shadow-xl' : ''}
        ${!achievement.earned ? 'bg-muted/50' : ''}
      `}
      onClick={onClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Sparkle effect for earned achievements */}
      {achievement.earned && achievement.rarity === 'legendary' && (
        <div className="absolute -top-1 -right-1">
          <Sparkles className="h-4 w-4 text-yellow-400 animate-pulse" />
        </div>
      )}

      <CardContent className={sizeClasses.card}>
        <div className="flex flex-col items-center text-center space-y-3">
          {/* Icon */}
          <div className={`
            ${sizeClasses.iconContainer} 
            rounded-full flex items-center justify-center relative
            ${achievement.earned 
              ? `bg-gradient-to-br ${getTierColor(achievement.tier)} text-white shadow-lg` 
              : 'bg-muted text-muted-foreground'
            }
          `}>
            {achievement.earned ? (
              <IconComponent className={sizeClasses.icon} />
            ) : (
              <Lock className={sizeClasses.icon} />
            )}
            
            {/* Tier indicator */}
            {achievement.earned && (
              <div className="absolute -bottom-1 -right-1">
                <div className={`
                  w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold
                  ${achievement.tier === 'bronze' ? 'bg-amber-600 text-white' :
                    achievement.tier === 'silver' ? 'bg-gray-500 text-white' :
                    achievement.tier === 'gold' ? 'bg-yellow-500 text-white' :
                    'bg-purple-500 text-white'
                  }
                `}>
                  {achievement.tier.charAt(0).toUpperCase()}
                </div>
              </div>
            )}
          </div>

          {/* Title and Description */}
          <div className="space-y-1">
            <h3 className={`font-semibold ${sizeClasses.title} ${!achievement.earned ? 'text-muted-foreground' : ''}`}>
              {achievement.title}
            </h3>
            <p className={`${sizeClasses.description} text-muted-foreground line-clamp-2`}>
              {achievement.description}
            </p>
          </div>

          {/* Progress Bar */}
          {showProgress && achievement.progress !== undefined && achievement.maxProgress && !achievement.earned && (
            <div className="w-full space-y-1">
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>Progress</span>
                <span>{achievement.progress}/{achievement.maxProgress}</span>
              </div>
              <Progress value={progressPercentage} className="h-1" />
            </div>
          )}

          {/* Badges */}
          <div className="flex flex-wrap gap-1 justify-center">
            {achievement.earned && (
              <Badge variant="secondary" className="text-xs bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                <CheckCircle className="mr-1 h-3 w-3" />
                Earned
              </Badge>
            )}
            
            {achievement.rarity && achievement.rarity !== 'common' && (
              <Badge 
                variant="outline" 
                className={`text-xs capitalize ${
                  achievement.rarity === 'rare' ? 'border-blue-400 text-blue-600' :
                  achievement.rarity === 'epic' ? 'border-purple-400 text-purple-600' :
                  'border-yellow-400 text-yellow-600'
                }`}
              >
                {achievement.rarity}
              </Badge>
            )}
            
            <Badge variant="outline" className="text-xs">
              {achievement.points} pts
            </Badge>
          </div>

          {/* Earned Date */}
          {achievement.earned && achievement.earnedDate && (
            <p className="text-xs text-muted-foreground">
              Earned {achievement.earnedDate}
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

// Achievement Gallery Component
interface AchievementGalleryProps {
  achievements: Achievement[];
  title?: string;
  showEarnedOnly?: boolean;
  onAchievementClick?: (achievement: Achievement) => void;
}

export const AchievementGallery = ({ 
  achievements, 
  title = "Achievements",
  showEarnedOnly = false,
  onAchievementClick 
}: AchievementGalleryProps) => {
  const [filter, setFilter] = useState<'all' | 'earned' | 'unearned'>('all');
  
  const filteredAchievements = achievements.filter(achievement => {
    if (showEarnedOnly) return achievement.earned;
    if (filter === 'earned') return achievement.earned;
    if (filter === 'unearned') return !achievement.earned;
    return true;
  });

  const earnedCount = achievements.filter(a => a.earned).length;
  const totalPoints = achievements.filter(a => a.earned).reduce((sum, a) => sum + a.points, 0);

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center">
              <Trophy className="mr-2 h-5 w-5" />
              {title}
            </CardTitle>
            <CardDescription>
              {earnedCount}/{achievements.length} earned • {totalPoints} points
            </CardDescription>
          </div>
          
          {!showEarnedOnly && (
            <div className="flex space-x-1">
              <Button
                variant={filter === 'all' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setFilter('all')}
              >
                All
              </Button>
              <Button
                variant={filter === 'earned' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setFilter('earned')}
              >
                Earned
              </Button>
              <Button
                variant={filter === 'unearned' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setFilter('unearned')}
              >
                Locked
              </Button>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {filteredAchievements.map((achievement) => (
            <AchievementBadge
              key={achievement.id}
              achievement={achievement}
              onClick={() => onAchievementClick?.(achievement)}
            />
          ))}
        </div>
        
        {filteredAchievements.length === 0 && (
          <div className="text-center py-8 text-muted-foreground">
            No achievements found for the selected filter.
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default AchievementBadge;
