import React from 'react';
import { TutorialBrowser } from '@/components/learning/tutorial-browser';
import { getTutorialStats } from '@/lib/tutorials-data';
import { BookOpen, Code, Rocket, Award, Target } from 'lucide-react';
import { Badge } from '@/components/ui/badge';

export const metadata = {
  title: 'Go Tutorials | Go Pro Learning Platform',
  description: 'Comprehensive Go programming tutorials from basics to advanced topics',
};

export default function TutorialsPage() {
  const stats = getTutorialStats();

  const tutorialStats = [
    {
      icon: BookOpen,
      value: stats.total,
      label: "Total Tutorials"
    },
    {
      icon: Code,
      value: stats.byDifficulty.beginner + stats.byDifficulty.intermediate,
      label: "Beginner-Friendly"
    },
    {
      icon: Rocket,
      value: stats.byDifficulty.advanced,
      label: "Advanced Topics"
    },
    {
      icon: Award,
      value: "100+",
      label: "Learning Hours"
    },
  ];

  return (
    <div className="flex flex-col min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden animated-gradient">
        {/* Decorative elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-96 h-96 bg-gradient-to-br from-primary/20 to-blue-500/20 rounded-full blur-3xl float-animation" />
          <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-gradient-to-tr from-cyan-500/20 to-primary/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '1s' }} />
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-gradient-to-r from-primary/10 via-transparent to-blue-500/10 rounded-full blur-3xl float-animation" style={{ animationDelay: '2s' }} />
        </div>

        <div className="container-responsive py-16 sm:py-24 lg:py-32 relative z-10">
          <div className="mx-auto max-w-5xl text-center">
            <Badge variant="secondary" className="mb-6 text-sm pulse-badge animate-in fade-in slide-in-bottom duration-500">
              📚 Learning Resources
            </Badge>

            <h1 className="text-3xl font-bold tracking-tight sm:text-5xl lg:text-6xl xl:text-7xl mb-6 animate-in fade-in slide-in-bottom duration-700">
              Comprehensive{" "}
              <span className="go-gradient-text">Go Tutorials</span>
              <br />
              From Basics to Advanced
            </h1>

            <p className="text-sm sm:text-base lg:text-lg xl:text-xl text-muted-foreground mb-8 max-w-xl sm:max-w-2xl mx-auto leading-relaxed animate-in fade-in slide-in-bottom duration-1000">
              Master Go programming with our structured tutorials covering everything from fundamentals to advanced patterns, microservices, and real-world applications.
            </p>

            {/* Stats */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 max-w-xs sm:max-w-2xl lg:max-w-4xl mx-auto animate-in fade-in duration-1000">
              {tutorialStats.map((stat, index) => (
                <div
                  key={`stat-${stat.label}-${index}`}
                  className="group text-center p-4 sm:p-6 rounded-xl glass-card hover:shadow-2xl hover:shadow-primary/20 transition-all duration-500 hover:-translate-y-2 relative overflow-hidden"
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  {/* Gradient overlay on hover */}
                  <div className="absolute inset-0 bg-gradient-to-br from-primary/10 via-transparent to-blue-500/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                  <div className="relative z-10">
                    <div className="flex justify-center mb-2 sm:mb-3">
                      <div className="p-3 rounded-xl bg-gradient-to-br from-primary/20 to-primary/10 group-hover:from-primary/30 group-hover:to-primary/20 transition-all duration-300 group-hover:scale-110">
                        <stat.icon className="h-5 w-5 sm:h-6 sm:w-6 text-primary group-hover:animate-pulse" />
                      </div>
                    </div>
                    <div className="text-xl sm:text-2xl lg:text-3xl font-bold bg-gradient-to-br from-foreground via-primary/80 to-foreground/70 bg-clip-text text-transparent group-hover:scale-110 transition-transform duration-300">{stat.value}</div>
                    <div className="text-xs sm:text-sm text-muted-foreground mt-1 group-hover:text-foreground/80 transition-colors">{stat.label}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Tutorial Browser Section */}
      <section className="py-16 sm:py-24 lg:py-32 relative overflow-hidden" aria-labelledby="tutorials-title">
        <div className="absolute inset-0 bg-gradient-to-b from-background via-accent/5 to-background pointer-events-none" />

        <div className="container-responsive relative z-10">
          <div className="mx-auto max-w-3xl text-center margin-responsive mb-12">
            <Badge variant="outline" className="mb-4 animate-in fade-in duration-500">
              <Target className="mr-2 h-3 w-3" />
              Browse All Topics
            </Badge>
            <h2
              id="tutorials-title"
              className="text-2xl font-bold tracking-tight sm:text-3xl lg:text-4xl mb-4 animate-in fade-in slide-in-bottom duration-700"
            >
              Find the <span className="go-gradient-text">Perfect Tutorial</span>
            </h2>
            <p className="text-lg text-muted-foreground animate-in fade-in slide-in-bottom duration-1000">
              Filter by topic, difficulty level, or search for specific concepts
            </p>
          </div>

          {/* Tutorial Browser */}
          <div className="animate-in fade-in duration-1000">
            <TutorialBrowser />
          </div>
        </div>
      </section>
    </div>
  );
}

