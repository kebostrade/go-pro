"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { 
  List, 
  ChevronRight, 
  ChevronDown,
  BookOpen,
  CheckCircle,
  Circle,
  Clock,
  Target,
  Eye,
  Hash
} from "lucide-react";

interface TOCItem {
  id: string;
  title: string;
  level: number;
  completed?: boolean;
  timeEstimate?: string;
  children?: TOCItem[];
}

interface TableOfContentsProps {
  items: TOCItem[];
  currentSection?: string;
  onSectionClick?: (sectionId: string) => void;
  showProgress?: boolean;
  className?: string;
}

export default function TableOfContents({
  items,
  currentSection,
  onSectionClick,
  showProgress = true,
  className = ""
}: TableOfContentsProps) {
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set());
  const [completedSections, setCompletedSections] = useState<Set<string>>(new Set());

  // Auto-expand sections that contain the current section
  useEffect(() => {
    if (currentSection) {
      const newExpanded = new Set(expandedSections);
      
      const findAndExpandParents = (items: TOCItem[], targetId: string, parentId?: string): boolean => {
        for (const item of items) {
          if (item.id === targetId) {
            if (parentId) newExpanded.add(parentId);
            return true;
          }
          if (item.children && findAndExpandParents(item.children, targetId, item.id)) {
            newExpanded.add(item.id);
            return true;
          }
        }
        return false;
      };

      findAndExpandParents(items, currentSection);
      setExpandedSections(newExpanded);
    }
  }, [currentSection]);

  const toggleExpanded = (sectionId: string) => {
    const newExpanded = new Set(expandedSections);
    if (newExpanded.has(sectionId)) {
      newExpanded.delete(sectionId);
    } else {
      newExpanded.add(sectionId);
    }
    setExpandedSections(newExpanded);
  };

  const markCompleted = (sectionId: string) => {
    const newCompleted = new Set(completedSections);
    if (newCompleted.has(sectionId)) {
      newCompleted.delete(sectionId);
    } else {
      newCompleted.add(sectionId);
    }
    setCompletedSections(newCompleted);
  };

  const calculateProgress = () => {
    const totalItems = countAllItems(items);
    const completedItems = completedSections.size;
    return totalItems > 0 ? (completedItems / totalItems) * 100 : 0;
  };

  const countAllItems = (items: TOCItem[]): number => {
    return items.reduce((count, item) => {
      return count + 1 + (item.children ? countAllItems(item.children) : 0);
    }, 0);
  };

  const renderTOCItem = (item: TOCItem, depth: number = 0) => {
    const isExpanded = expandedSections.has(item.id);
    const isCompleted = completedSections.has(item.id) || item.completed;
    const isCurrent = currentSection === item.id;
    const hasChildren = item.children && item.children.length > 0;

    const indentClass = depth === 0 ? '' : `ml-${Math.min(depth * 4, 12)}`;
    
    return (
      <div key={item.id} className="space-y-1">
        <div
          className={`flex items-center space-x-2 p-2 rounded-lg cursor-pointer transition-all hover:bg-muted/50 group ${
            isCurrent ? 'bg-primary/10 border border-primary/20' : ''
          } ${indentClass}`}
          onClick={() => onSectionClick?.(item.id)}
        >
          {/* Expand/Collapse Button */}
          {hasChildren && (
            <Button
              variant="ghost"
              size="sm"
              className="h-6 w-6 p-0"
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(item.id);
              }}
            >
              {isExpanded ? (
                <ChevronDown className="h-3 w-3" />
              ) : (
                <ChevronRight className="h-3 w-3" />
              )}
            </Button>
          )}

          {/* Completion Status */}
          <Button
            variant="ghost"
            size="sm"
            className="h-6 w-6 p-0"
            onClick={(e) => {
              e.stopPropagation();
              markCompleted(item.id);
            }}
          >
            {isCompleted ? (
              <CheckCircle className="h-4 w-4 text-green-500" />
            ) : (
              <Circle className="h-4 w-4 text-muted-foreground group-hover:text-primary" />
            )}
          </Button>

          {/* Section Icon */}
          <div className="flex-shrink-0">
            {item.level === 1 && <BookOpen className="h-4 w-4 text-primary" />}
            {item.level === 2 && <Target className="h-4 w-4 text-blue-500" />}
            {item.level >= 3 && <Hash className="h-3 w-3 text-muted-foreground" />}
          </div>

          {/* Title and Metadata */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between">
              <span className={`text-sm font-medium truncate ${
                isCurrent ? 'text-primary' : 'text-foreground group-hover:text-primary'
              }`}>
                {item.title}
              </span>
              
              {item.timeEstimate && (
                <Badge variant="outline" className="text-xs ml-2">
                  <Clock className="mr-1 h-3 w-3" />
                  {item.timeEstimate}
                </Badge>
              )}
            </div>
          </div>

          {/* Current Section Indicator */}
          {isCurrent && (
            <div className="w-2 h-2 rounded-full bg-primary animate-pulse" />
          )}
        </div>

        {/* Children */}
        {hasChildren && isExpanded && (
          <div className="space-y-1 animate-in slide-in-from-top-2 duration-200">
            {item.children!.map(child => renderTOCItem(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  const progress = calculateProgress();

  return (
    <Card className={`glass-card border-2 ${className}`}>
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center text-lg">
          <List className="mr-2 h-5 w-5 text-primary" />
          Table of Contents
        </CardTitle>
        
        {showProgress && (
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Progress</span>
              <span className="font-medium">{Math.round(progress)}%</span>
            </div>
            <Progress value={progress} className="h-2" />
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>{completedSections.size} completed</span>
              <span>{countAllItems(items)} total sections</span>
            </div>
          </div>
        )}
      </CardHeader>

      <CardContent className="pt-0">
        <div className="space-y-1 max-h-96 overflow-y-auto custom-scrollbar">
          {items.length > 0 ? (
            items.map(item => renderTOCItem(item))
          ) : (
            <div className="text-center py-8 text-muted-foreground">
              <List className="mx-auto h-12 w-12 mb-4 opacity-50" />
              <p>No sections available</p>
            </div>
          )}
        </div>

        {/* Quick Actions */}
        <div className="flex items-center justify-between mt-4 pt-4 border-t">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              const allIds = new Set<string>();
              const collectIds = (items: TOCItem[]) => {
                items.forEach(item => {
                  allIds.add(item.id);
                  if (item.children) collectIds(item.children);
                });
              };
              collectIds(items);
              setExpandedSections(allIds);
            }}
            className="text-xs"
          >
            Expand All
          </Button>
          
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setExpandedSections(new Set())}
            className="text-xs"
          >
            Collapse All
          </Button>
          
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              const allIds = new Set<string>();
              const collectIds = (items: TOCItem[]) => {
                items.forEach(item => {
                  allIds.add(item.id);
                  if (item.children) collectIds(item.children);
                });
              };
              collectIds(items);
              setCompletedSections(allIds);
            }}
            className="text-xs"
          >
            <CheckCircle className="mr-1 h-3 w-3" />
            Mark All
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
