'use client';

import React, { useState, useMemo } from 'react';
import { tutorials, tutorialCategories, type Tutorial } from '@/lib/tutorials-data';
import { Search, Filter, Clock, BookOpen, Star, ChevronRight } from 'lucide-react';
import Link from 'next/link';

export const TutorialBrowser: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [selectedDifficulty, setSelectedDifficulty] = useState<string>('all');

  const filteredTutorials = useMemo(() => {
    return tutorials.filter(tutorial => {
      const matchesSearch = tutorial.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          tutorial.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          tutorial.topics.some(topic => topic.toLowerCase().includes(searchQuery.toLowerCase()));
      
      const matchesCategory = selectedCategory === 'all' || tutorial.category === selectedCategory;
      const matchesDifficulty = selectedDifficulty === 'all' || tutorial.difficulty === selectedDifficulty;

      return matchesSearch && matchesCategory && matchesDifficulty;
    });
  }, [searchQuery, selectedCategory, selectedDifficulty]);

  return (
    <div className="w-full space-y-8">
      {/* Header */}
      <div className="text-center space-y-4">
        <h1 className="text-4xl md:text-5xl font-bold bg-gradient-to-r from-go-blue to-go-cyan bg-clip-text text-transparent">
          Go Learning Tutorials
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
          Master Go programming with {tutorials.length} comprehensive tutorials covering everything from basics to advanced topics
        </p>
      </div>

      {/* Search and Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 space-y-4">
        {/* Search Bar */}
        <div className="relative">
          <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
          <input
            type="text"
            placeholder="Search tutorials, topics, or technologies..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-12 pr-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-go-blue focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
          />
        </div>

        {/* Filters */}
        <div className="flex flex-wrap gap-4">
          {/* Category Filter */}
          <div className="flex-1 min-w-[200px]">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              <Filter className="inline w-4 h-4 mr-1" />
              Category
            </label>
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-go-blue bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="all">All Categories</option>
              {Object.entries(tutorialCategories).map(([key, cat]) => (
                <option key={key} value={key}>
                  {cat.icon} {cat.name}
                </option>
              ))}
            </select>
          </div>

          {/* Difficulty Filter */}
          <div className="flex-1 min-w-[200px]">
            <label htmlFor="difficulty-filter" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Difficulty
            </label>
            <select
              id="difficulty-filter"
              value={selectedDifficulty}
              onChange={(e) => setSelectedDifficulty(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-go-blue bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="all">All Levels</option>
              <option value="beginner">Beginner</option>
              <option value="intermediate">Intermediate</option>
              <option value="advanced">Advanced</option>
            </select>
          </div>
        </div>

        {/* Results Count */}
        <div className="text-sm text-gray-600 dark:text-gray-400">
          Showing {filteredTutorials.length} of {tutorials.length} tutorials
        </div>
      </div>

      {/* Tutorial Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredTutorials.map((tutorial) => (
          <TutorialCard key={tutorial.id} tutorial={tutorial} />
        ))}
      </div>

      {/* No Results */}
      {filteredTutorials.length === 0 && (
        <div className="text-center py-12">
          <BookOpen className="w-16 h-16 mx-auto text-gray-400 mb-4" />
          <h3 className="text-xl font-semibold text-gray-700 dark:text-gray-300 mb-2">
            No tutorials found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            Try adjusting your search or filters
          </p>
        </div>
      )}
    </div>
  );
};

const TutorialCard: React.FC<{ tutorial: Tutorial }> = ({ tutorial }) => {
  const difficultyColors = {
    beginner: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
    intermediate: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
    advanced: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  };

  return (
    <Link href={`/tutorials/${tutorial.id}`}>
      <div className="group bg-white dark:bg-gray-800 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border border-gray-200 dark:border-gray-700 hover:border-go-blue cursor-pointer h-full flex flex-col">
        {/* Header with Soft Gradient */}
        <div className="relative overflow-hidden">
          <div className={`absolute inset-0 bg-gradient-to-r ${tutorial.color} opacity-80`}></div>
          <div className="relative bg-gradient-to-br from-white/10 to-transparent backdrop-blur-sm p-6 text-white">
            <div className="absolute top-0 right-0 text-8xl opacity-10 transform translate-x-4 -translate-y-4">
              {tutorial.icon}
            </div>
            <div className="relative z-10">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-semibold opacity-90">Tutorial {tutorial.number}</span>
                {tutorial.featured && (
                  <Star className="w-5 h-5 fill-current" />
                )}
              </div>
              <h3 className="text-xl font-bold mb-2 group-hover:scale-105 transition-transform">
                {tutorial.title}
              </h3>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 flex-1 flex flex-col">
          <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-3">
            {tutorial.description}
          </p>

          {/* Topics */}
          <div className="flex flex-wrap gap-2 mb-4">
            {tutorial.topics.slice(0, 3).map((topic) => (
              <span
                key={topic}
                className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded-full"
              >
                {topic}
              </span>
            ))}
            {tutorial.topics.length > 3 && (
              <span className="px-2 py-1 text-gray-500 text-xs">
                +{tutorial.topics.length - 3} more
              </span>
            )}
          </div>

          {/* Footer */}
          <div className="mt-auto space-y-3">
            <div className="flex items-center justify-between text-sm">
              <span className={`px-3 py-1 rounded-full font-medium ${difficultyColors[tutorial.difficulty]}`}>
                {tutorial.difficulty.charAt(0).toUpperCase() + tutorial.difficulty.slice(1)}
              </span>
              <div className="flex items-center text-gray-600 dark:text-gray-400">
                <Clock className="w-4 h-4 mr-1" />
                {tutorial.duration}
              </div>
            </div>

            <div className="flex items-center justify-between pt-3 border-t border-gray-200 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">
                {tutorial.learningOutcomes.length} outcomes
              </span>
              <ChevronRight className="w-5 h-5 text-go-blue group-hover:translate-x-1 transition-transform" />
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
};

export default TutorialBrowser;

