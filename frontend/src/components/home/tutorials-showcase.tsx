'use client';

import React from 'react';
import Link from 'next/link';
import { getFeaturedTutorials, getTutorialStats } from '@/lib/tutorials-data';
import { GraduationCap, Clock, Star, ChevronRight, BookOpen, Code, Rocket } from 'lucide-react';

export const TutorialsShowcase: React.FC = () => {
  const featuredTutorials = getFeaturedTutorials();
  const stats = getTutorialStats();

  return (
    <section className="py-20 bg-gradient-to-b from-white to-gray-50 dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4">
        {/* Section Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-go-blue/10 rounded-full mb-4">
            <GraduationCap className="w-5 h-5 text-go-blue" />
            <span className="text-go-blue font-semibold">Comprehensive Learning Path</span>
          </div>
          <h2 className="text-4xl md:text-5xl font-bold mb-4 bg-gradient-to-r from-go-blue to-go-cyan bg-clip-text text-transparent">
            19 Expert-Crafted Tutorials
          </h2>
          <p className="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">
            From WebSocket real-time apps to Kubernetes deployments, master Go with production-ready projects
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
          <StatCard
            icon={<BookOpen className="w-8 h-8" />}
            value={stats.total}
            label="Total Tutorials"
            description="Comprehensive coverage"
            color="from-blue-500 to-cyan-500"
          />
          <StatCard
            icon={<Code className="w-8 h-8" />}
            value="10,000+"
            label="Lines of Code"
            description="Production-ready examples"
            color="from-purple-500 to-pink-500"
          />
          <StatCard
            icon={<Rocket className="w-8 h-8" />}
            value="100+"
            label="Learning Hours"
            description="Hands-on practice"
            color="from-green-500 to-teal-500"
          />
        </div>

        {/* Featured Tutorials Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          {featuredTutorials.slice(0, 6).map((tutorial) => (
            <Link
              key={tutorial.id}
              href={`/tutorials/${tutorial.id}`}
              className="group"
            >
              <div className="h-full bg-white dark:bg-gray-800 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border border-gray-200 dark:border-gray-700 hover:border-go-blue">
                {/* Gradient Header */}
                <div className={`bg-gradient-to-r ${tutorial.color} p-6 text-white relative overflow-hidden`}>
                  <div className="absolute top-0 right-0 text-7xl opacity-10 transform translate-x-4 -translate-y-4">
                    {tutorial.icon}
                  </div>
                  <div className="relative z-10">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-semibold opacity-90">
                        Tutorial {tutorial.number}
                      </span>
                      {tutorial.featured && (
                        <Star className="w-5 h-5 fill-current" />
                      )}
                    </div>
                    <h3 className="text-lg font-bold group-hover:scale-105 transition-transform">
                      {tutorial.title}
                    </h3>
                  </div>
                </div>

                {/* Content */}
                <div className="p-6">
                  <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 line-clamp-2">
                    {tutorial.description}
                  </p>

                  {/* Topics */}
                  <div className="flex flex-wrap gap-2 mb-4">
                    {tutorial.topics.slice(0, 2).map((topic, index) => (
                      <span
                        key={index}
                        className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs rounded-full"
                      >
                        {topic}
                      </span>
                    ))}
                    {tutorial.topics.length > 2 && (
                      <span className="px-2 py-1 text-gray-500 text-xs">
                        +{tutorial.topics.length - 2}
                      </span>
                    )}
                  </div>

                  {/* Footer */}
                  <div className="flex items-center justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center text-sm text-gray-600 dark:text-gray-400">
                      <Clock className="w-4 h-4 mr-1" />
                      {tutorial.duration}
                    </div>
                    <ChevronRight className="w-5 h-5 text-go-blue group-hover:translate-x-1 transition-transform" />
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>

        {/* CTA */}
        <div className="text-center">
          <Link
            href="/tutorials"
            className="inline-flex items-center gap-2 px-8 py-4 bg-gradient-to-r from-go-blue to-go-cyan text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transition-all hover:scale-105"
          >
            <GraduationCap className="w-5 h-5" />
            Explore All 19 Tutorials
            <ChevronRight className="w-5 h-5" />
          </Link>
          <p className="mt-4 text-sm text-gray-600 dark:text-gray-400">
            From basics to advanced topics • Production-ready code • Real-world projects
          </p>
        </div>
      </div>
    </section>
  );
};

const StatCard: React.FC<{
  icon: React.ReactNode;
  value: string | number;
  label: string;
  description: string;
  color: string;
}> = ({ icon, value, label, description, color }) => (
  <div className={`bg-gradient-to-r ${color} rounded-xl p-6 text-white shadow-lg hover:shadow-xl transition-all`}>
    <div className="flex items-center justify-between mb-4">
      <div className="p-3 bg-white/20 backdrop-blur-sm rounded-lg">
        {icon}
      </div>
    </div>
    <div className="text-4xl font-bold mb-2">{value}</div>
    <div className="text-lg font-semibold mb-1">{label}</div>
    <div className="text-sm opacity-90">{description}</div>
  </div>
);

export default TutorialsShowcase;

