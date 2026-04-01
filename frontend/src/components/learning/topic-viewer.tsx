'use client';

import React, { useState } from 'react';
import { type Topic } from '@/lib/topics-data';
import { 
  BookOpen, Clock, Target, CheckCircle, Code, Terminal, 
  Play, Download, ExternalLink, ChevronRight, Award, GitBranch
} from 'lucide-react';
import Link from 'next/link';
import ExerciseCard from './exercise-card';
import CodeEditor from '../workspace/CodeEditor';
import { DockerPanel } from './docker-panel';

interface TopicViewerProps {
  topic: Topic;
  completedExercises?: Set<string>;
  onExerciseComplete?: (exerciseId: string) => void;
}

export const TopicViewer: React.FC<TopicViewerProps> = ({ 
  topic, 
  completedExercises = new Set(),
  onExerciseComplete 
}) => {
  const [activeTab, setActiveTab] = useState<'overview' | 'content' | 'practice'>('overview');

  return (
    <div className="max-w-7xl mx-auto space-y-8">
      {/* Hero Section */}
      <div className={`bg-gradient-to-r ${topic.color} rounded-2xl p-8 md:p-12 text-white shadow-2xl`}>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-3 mb-4">
              <span className="text-6xl">{topic.icon}</span>
              <div>
                <p className="text-sm opacity-90 font-semibold">Phase {topic.phase}: {topic.phaseName}</p>
                <h1 className="text-3xl md:text-4xl font-bold">{topic.title}</h1>
              </div>
            </div>
            <p className="text-lg opacity-95 mb-6 max-w-3xl">{topic.description}</p>
            <div className="flex flex-wrap gap-4">
              <div className="flex items-center gap-2 bg-white/20 backdrop-blur-sm px-4 py-2 rounded-lg">
                <Clock className="w-5 h-5" />
                <span className="font-medium">{topic.duration}</span>
              </div>
              <div className="flex items-center gap-2 bg-white/20 backdrop-blur-sm px-4 py-2 rounded-lg">
                <Target className="w-5 h-5" />
                <span className="font-medium capitalize">{topic.difficulty}</span>
              </div>
              <div className="flex items-center gap-2 bg-white/20 backdrop-blur-sm px-4 py-2 rounded-lg">
                <BookOpen className="w-5 h-5" />
                <span className="font-medium">{topic.exercises.length} Exercises</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
        <div className="flex border-b border-gray-200 dark:border-gray-700">
          <TabButton
            active={activeTab === 'overview'}
            onClick={() => setActiveTab('overview')}
            icon={<BookOpen className="w-5 h-5" />}
            label="Overview"
          />
          <TabButton
            active={activeTab === 'content'}
            onClick={() => setActiveTab('content')}
            icon={<Code className="w-5 h-5" />}
            label="Content"
          />
          <TabButton
            active={activeTab === 'practice'}
            onClick={() => setActiveTab('practice')}
            icon={<Terminal className="w-5 h-5" />}
            label="Practice"
          />
        </div>

        <div className="p-8">
          {activeTab === 'overview' && <OverviewTab topic={topic} />}
          {activeTab === 'content' && <ContentTab topic={topic} />}
          {activeTab === 'practice' && (
            <PracticeTab 
              topic={topic} 
              completedExercises={completedExercises}
              onComplete={onExerciseComplete}
            />
          )}
        </div>
      </div>
    </div>
  );
};

const TabButton: React.FC<{
  active: boolean;
  onClick: () => void;
  icon: React.ReactNode;
  label: string;
}> = ({ active, onClick, icon, label }) => (
  <button
    onClick={onClick}
    className={`flex items-center gap-2 px-6 py-4 font-medium transition-colors ${
      active
        ? 'text-blue-600 border-b-2 border-blue-600 bg-blue-50 dark:bg-gray-700 dark:text-blue-400 dark:border-blue-400'
        : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-50 dark:hover:bg-gray-700'
    }`}
  >
    {icon}
    <span>{label}</span>
  </button>
);

const OverviewTab: React.FC<{ topic: Topic }> = ({ topic }) => (
  <div className="space-y-8">
    {/* Topics */}
    <section>
      <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Topics Covered</h2>
      <div className="flex flex-wrap gap-3">
        {topic.topics.map((t, index) => (
          <span
            key={index}
            className="px-4 py-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-lg font-medium shadow-md"
          >
            {t}
          </span>
        ))}
      </div>
    </section>

    {/* Learning Outcomes */}
    <section>
      <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-white flex items-center gap-2">
        <Award className="w-6 h-6 text-blue-600 dark:text-blue-400" />
        Learning Outcomes
      </h2>
      <div className="grid gap-3">
        {topic.learningOutcomes.map((outcome, index) => (
          <div
            key={index}
            className="flex items-start gap-3 p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800"
          >
            <CheckCircle className="w-5 h-5 text-green-600 dark:text-green-400 flex-shrink-0 mt-0.5" />
            <span className="text-gray-700 dark:text-gray-300">{outcome}</span>
          </div>
        ))}
      </div>
    </section>

    {/* Prerequisites */}
    {topic.prerequisites.length > 0 && (
      <section>
        <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Prerequisites</h2>
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6">
          <p className="text-sm text-yellow-800 dark:text-yellow-200 mb-3 font-medium">
            Before starting this topic, you should be familiar with:
          </p>
          <ul className="space-y-2">
            {topic.prerequisites.map((prereq, index) => (
              <li key={index} className="flex items-center gap-2 text-gray-700 dark:text-gray-300">
                <ChevronRight className="w-4 h-4 text-yellow-600 dark:text-yellow-400" />
                {prereq}
              </li>
            ))}
          </ul>
        </div>
      </section>
    )}

    {/* Quick Actions */}
    <section>
      <h2 className="text-2xl font-bold mb-4 text-gray-900 dark:text-white">Get Started</h2>
      <div className="grid md:grid-cols-2 gap-4">
        <button
          onClick={() => {
            const practiceTab = document.querySelector('[data-tab="practice"]');
            if (practiceTab) (practiceTab as HTMLElement).click();
          }}
          className="flex items-center justify-between p-6 bg-gradient-to-r from-blue-500 to-cyan-500 text-white rounded-xl shadow-lg hover:shadow-xl transition-all group cursor-pointer"
        >
          <div>
            <h3 className="font-bold text-lg mb-1">Start Learning</h3>
            <p className="text-sm opacity-90">Begin with exercises</p>
          </div>
          <Play className="w-8 h-8 group-hover:scale-110 transition-transform" />
        </button>

        <a
          href={topic.githubUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between p-6 bg-gray-100 dark:bg-gray-700 rounded-xl shadow-lg hover:shadow-xl transition-all group"
        >
          <div>
            <h3 className="font-bold text-lg mb-1 text-gray-900 dark:text-white">View Code</h3>
            <p className="text-sm text-gray-600 dark:text-gray-400">Browse on GitHub</p>
          </div>
          <ExternalLink className="w-8 h-8 text-gray-600 dark:text-gray-400 group-hover:scale-110 transition-transform" />
        </a>
      </div>
    </section>
  </div>
);

const ContentTab: React.FC<{ topic: Topic }> = ({ topic }) => (
  <div className="space-y-6">
    <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
      <h3 className="font-bold text-lg mb-2 text-gray-900 dark:text-white">Project Path</h3>
      <p className="text-gray-700 dark:text-gray-300 mb-4">
        The complete project code is available in the following directory:
      </p>
      <code className="block bg-gray-900 text-green-400 p-4 rounded-lg font-mono text-sm">
        {topic.projectPath}
      </code>
    </div>

    <div className="grid gap-4">
      <ActionCard
        icon={<Download className="w-6 h-6" />}
        title="Clone Repository"
        description="Get the full project code"
        command="git clone https://github.com/DimaJoyti/go-pro.git"
      />
      <ActionCard
        icon={<GitBranch className="w-6 h-6" />}
        title="Navigate to Project"
        description="Change to the project directory"
        command={`cd go-pro/${topic.projectPath}`}
      />
      <ActionCard
        icon={<Terminal className="w-6 h-6" />}
        title="Run the Project"
        description="Build and execute the code"
        command="go run ."
      />
    </div>

    <div className="mt-6">
      <a
        href={topic.githubUrl}
        target="_blank"
        rel="noopener noreferrer"
        className="flex items-center justify-center gap-2 w-full p-4 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
      >
        <ExternalLink className="w-5 h-5" />
        <span className="font-medium text-gray-900 dark:text-white">View on GitHub</span>
      </a>
    </div>

    {/* Docker Environment Section */}
    <div className="mt-8 pt-8 border-t border-gray-200 dark:border-gray-700">
      <h3 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white flex items-center gap-2">
        <Terminal className="w-5 h-5" />
        Docker Environment
      </h3>
      <DockerPanel topicId={topic.id} />
      <p className="mt-3 text-sm text-slate-600 dark:text-slate-400">
        Start a Docker environment with all services needed for this topic.
        Make sure Docker Desktop is running before starting.
      </p>
    </div>
  </div>
);

interface PracticeTabProps {
  topic: Topic;
  completedExercises?: Set<string>;
  onComplete?: (id: string) => void;
}

const PracticeTab: React.FC<PracticeTabProps> = ({ 
  topic, 
  completedExercises = new Set(),
  onComplete 
}) => (
  <div className="space-y-6">
    {/* Code Editor Section */}
    <div className="space-y-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-2xl font-bold">Practice</h2>
      </div>
      <CodeEditor 
        topicId={topic.id} 
        initialCode={`package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`}
      />
    </div>

    {/* Exercise Cards */}
    <div className="space-y-6 mt-8">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-2xl font-bold">Exercises</h2>
        <span className="text-muted-foreground">
          {completedExercises.size} of {topic.exercises.length} completed
        </span>
      </div>
      <div className="grid gap-4">
        {topic.exercises.map((exercise, index) => (
          <ExerciseCard 
            key={exercise.id} 
            exercise={exercise}
            topic={topic}
            index={index + 1}
            completed={completedExercises.has(exercise.id)}
            onComplete={onComplete}
          />
        ))}
      </div>
    </div>
  </div>
);

const ActionCard: React.FC<{
  icon: React.ReactNode;
  title: string;
  description: string;
  command: string;
}> = ({ icon, title, description, command }) => (
  <div className="flex items-start gap-4 p-6 bg-white dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600">
    <div className="p-3 bg-blue-500/10 rounded-lg text-blue-600 dark:text-blue-400">
      {icon}
    </div>
    <div className="flex-1">
      <h4 className="font-bold text-gray-900 dark:text-white mb-1">{title}</h4>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">{description}</p>
      <code className="block bg-gray-900 text-green-400 p-3 rounded font-mono text-sm">
        {command}
      </code>
    </div>
  </div>
);

export default TopicViewer;
