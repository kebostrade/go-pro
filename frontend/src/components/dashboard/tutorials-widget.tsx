'use client';

import React from 'react';
import Link from 'next/link';
import { getFeaturedTutorials } from '@/lib/tutorials-data';
import { GraduationCap, Clock, ChevronRight, Star } from 'lucide-react';

export const TutorialsWidget: React.FC = () => {
  const featuredTutorials = getFeaturedTutorials().slice(0, 3);

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-200 dark:border-gray-700">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-gradient-to-r from-go-blue to-go-cyan rounded-lg">
            <GraduationCap className="w-6 h-6 text-white" />
          </div>
          <div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">
              Featured Tutorials
            </h2>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Continue your learning journey
            </p>
          </div>
        </div>
        <Link
          href="/tutorials"
          className="text-go-blue hover:text-go-cyan font-medium text-sm flex items-center gap-1 transition-colors"
        >
          View All
          <ChevronRight className="w-4 h-4" />
        </Link>
      </div>

      {/* Tutorial List */}
      <div className="space-y-3">
        {featuredTutorials.map((tutorial) => (
          <Link
            key={tutorial.id}
            href={`/tutorials/${tutorial.id}`}
            className="block group"
          >
            <div className="p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-go-blue dark:hover:border-go-blue transition-all hover:shadow-md">
              <div className="flex items-start gap-3">
                <div className={`text-3xl flex-shrink-0 group-hover:scale-110 transition-transform`}>
                  {tutorial.icon}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-go-blue transition-colors truncate">
                      {tutorial.title}
                    </h3>
                    {tutorial.featured && (
                      <Star className="w-4 h-4 text-yellow-500 fill-current flex-shrink-0" />
                    )}
                  </div>
                  <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2 mb-2">
                    {tutorial.description}
                  </p>
                  <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-500">
                    <span className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      {tutorial.duration}
                    </span>
                    <span className="px-2 py-0.5 bg-gray-100 dark:bg-gray-700 rounded-full capitalize">
                      {tutorial.difficulty}
                    </span>
                  </div>
                </div>
                <ChevronRight className="w-5 h-5 text-gray-400 group-hover:text-go-blue group-hover:translate-x-1 transition-all flex-shrink-0" />
              </div>
            </div>
          </Link>
        ))}
      </div>

      {/* CTA */}
      <Link
        href="/tutorials"
        className="mt-4 block w-full py-3 bg-gradient-to-r from-go-blue to-go-cyan text-white text-center font-medium rounded-lg hover:shadow-lg transition-all"
      >
        Explore All 19 Tutorials
      </Link>
    </div>
  );
};

export default TutorialsWidget;

