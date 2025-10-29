import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  ArrowRight,
  BookOpen,
  Code2,
  Trophy,
  Users,
  Star,
  Play,
  CheckCircle,
  Clock,
  Target,
  Zap,
  Globe,
  Award,
  TrendingUp
} from "lucide-react";
import Link from "next/link";

export default function Home() {
  const stats = [
    { label: "Active Learners", value: "10,000+", icon: Users },
    { label: "Lessons Completed", value: "50,000+", icon: BookOpen },
    { label: "Code Challenges", value: "500+", icon: Code2 },
    { label: "Success Rate", value: "94%", icon: Trophy },
  ];

  const features = [
    {
      icon: BookOpen,
      title: "Interactive Lessons",
      description: "Learn Go through hands-on, interactive lessons that adapt to your pace.",
    },
    {
      icon: Code2,
      title: "Real Code Practice",
      description: "Write actual Go code with instant feedback and automated testing.",
    },
    {
      icon: Trophy,
      title: "Project-Based Learning",
      description: "Build real applications from CLI tools to microservices.",
    },
    {
      icon: Users,
      title: "Community Support",
      description: "Join thousands of Go developers in our supportive community.",
    },
    {
      icon: Target,
      title: "Personalized Path",
      description: "AI-powered learning paths tailored to your goals and experience.",
    },
    {
      icon: Award,
      title: "Industry Recognition",
      description: "Earn certificates recognized by top tech companies.",
    },
  ];

  const testimonials = [
    {
      name: "Sarah Chen",
      role: "Backend Developer at Google",
      avatar: "/avatars/sarah.svg",
      content: "GO-PRO transformed my understanding of Go. The hands-on approach made complex concepts click instantly.",
      rating: 5,
    },
    {
      name: "Marcus Rodriguez",
      role: "Senior Engineer at Uber",
      avatar: "/avatars/marcus.svg",
      content: "Best Go learning platform I've used. The project-based approach prepared me for real-world development.",
      rating: 5,
    },
    {
      name: "Emily Johnson",
      role: "DevOps Engineer at Netflix",
      avatar: "/avatars/emily.svg",
      content: "The microservices module was incredible. I went from beginner to building production systems.",
      rating: 5,
    },
  ];

  return (
    <div className="flex flex-col">
      {/* Hero Section */}
      <section
        className="relative overflow-hidden animated-gradient"
        aria-labelledby="hero-title"
      >
        {/* Enhanced Decorative elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-96 h-96 bg-gradient-to-br from-primary/20 to-blue-500/20 rounded-full blur-3xl float-animation" />
          <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-gradient-to-tr from-cyan-500/20 to-primary/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '1s' }} />
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-gradient-to-r from-primary/10 via-transparent to-blue-500/10 rounded-full blur-3xl float-animation" style={{ animationDelay: '2s' }} />
        </div>

        <div className="container-responsive py-16 sm:py-24 lg:py-32 xl:py-40 relative z-10">
          <div className="mx-auto max-w-5xl text-center">
            <Badge variant="secondary" className="mb-6 text-sm pulse-badge animate-in fade-in slide-in-bottom duration-500">
              🚀 New: Advanced Microservices Course Available
            </Badge>

            <h1
              id="hero-title"
              className="text-3xl font-bold tracking-tight sm:text-5xl lg:text-6xl xl:text-7xl mb-6 animate-in fade-in slide-in-bottom duration-700"
            >
              Master{" "}
              <span className="go-gradient-text">Go Programming</span>
              <br />
              Through Practice
            </h1>

            <p className="text-sm sm:text-base lg:text-lg xl:text-xl text-muted-foreground mb-6 sm:mb-8 max-w-xl sm:max-w-2xl mx-auto leading-relaxed animate-in fade-in slide-in-bottom duration-1000">
              Learn Go from basics to microservices with interactive lessons, real code practice,
              and projects that prepare you for production development.
            </p>

            <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center mb-8 sm:mb-12 animate-in fade-in slide-in-bottom duration-1000">
              <Link href="/learn/1">
                <Button size="lg" className="go-gradient text-white text-sm sm:text-base lg:text-lg px-6 sm:px-8 py-4 sm:py-6 min-h-[48px] sm:min-h-[56px] shadow-lg hover:shadow-2xl">
                  <Play className="mr-2 h-4 w-4 sm:h-5 sm:w-5" />
                  Start Learning Free
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </Link>
              <Link href="/curriculum">
                <Button size="lg" variant="outline" className="text-sm sm:text-base lg:text-lg px-6 sm:px-8 py-4 sm:py-6 min-h-[48px] sm:min-h-[56px]">
                  <BookOpen className="mr-2 h-4 w-4 sm:h-5 sm:w-5" />
                  View Curriculum
                </Button>
              </Link>
            </div>

            {/* Enhanced Stats */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 max-w-xs sm:max-w-2xl lg:max-w-4xl mx-auto animate-in fade-in duration-1000">
              {stats.map((stat, index) => (
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

      {/* Features Section */}
      <section className="py-16 sm:py-24 lg:py-32 relative overflow-hidden" aria-labelledby="features-title">
        <div className="absolute inset-0 bg-gradient-to-b from-background via-accent/5 to-background pointer-events-none" />

        <div className="container-responsive relative z-10">
          <div className="mx-auto max-w-3xl text-center margin-responsive">
            <Badge variant="outline" className="mb-4 animate-in fade-in duration-500">
              ✨ Features
            </Badge>
            <h2
              id="features-title"
              className="text-2xl font-bold tracking-tight sm:text-3xl lg:text-4xl mb-4 animate-in fade-in slide-in-bottom duration-700"
            >
              Why Choose <span className="go-gradient-text">GO-PRO</span>?
            </h2>
            <p className="text-lg text-muted-foreground animate-in fade-in slide-in-bottom duration-1000">
              Our platform combines the best of interactive learning with real-world application
            </p>
          </div>

          <div className="grid-responsive-cards gap-responsive">
            {features.map((feature, index) => (
              <Card
                key={`feature-${feature.title}-${index}`}
                className="glass-card h-full group hover:border-primary/50 transition-all duration-500 animate-in fade-in scale-in relative overflow-hidden"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                {/* Animated gradient background on hover */}
                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-blue-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                <CardHeader className="pb-3 sm:pb-4 relative z-10">
                  <div className="flex items-center space-x-3">
                    <div className="flex h-12 w-12 sm:h-14 sm:w-14 items-center justify-center rounded-xl bg-gradient-to-br from-primary/20 to-primary/10 flex-shrink-0 group-hover:scale-110 group-hover:rotate-3 transition-all duration-300 shadow-lg group-hover:shadow-xl group-hover:shadow-primary/30">
                      <feature.icon className="h-6 w-6 sm:h-7 sm:w-7 text-primary group-hover:scale-110 transition-transform duration-300" />
                    </div>
                    <CardTitle className="text-base sm:text-lg lg:text-xl leading-tight group-hover:text-primary transition-colors duration-300">{feature.title}</CardTitle>
                  </div>
                </CardHeader>
                <CardContent className="pt-0 relative z-10">
                  <CardDescription className="text-sm sm:text-base leading-relaxed group-hover:text-foreground/80 transition-colors duration-300">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Learning Path Section */}
      <section className="py-20 sm:py-32 relative overflow-hidden">
        <div className="absolute inset-0 animated-gradient opacity-50" />

        <div className="container max-w-screen-2xl px-4 relative z-10">
          <div className="mx-auto max-w-2xl text-center mb-16">
            <Badge variant="outline" className="mb-4 animate-in fade-in duration-500">
              📚 Curriculum
            </Badge>
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4 animate-in fade-in slide-in-bottom duration-700">
              Your <span className="go-gradient-text">Learning Journey</span>
            </h2>
            <p className="text-lg text-muted-foreground animate-in fade-in slide-in-bottom duration-1000">
              Structured curriculum that takes you from beginner to Go expert
            </p>
          </div>

          <Tabs defaultValue="beginner" className="max-w-4xl mx-auto animate-in fade-in duration-1000">
            <TabsList className="grid w-full grid-cols-3 bg-card/50 backdrop-blur-sm border border-border/50 shadow-lg">
              <TabsTrigger value="beginner" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
                <Zap className="mr-2 h-4 w-4" />
                Beginner
              </TabsTrigger>
              <TabsTrigger value="intermediate" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
                <Target className="mr-2 h-4 w-4" />
                Intermediate
              </TabsTrigger>
              <TabsTrigger value="advanced" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
                <Trophy className="mr-2 h-4 w-4" />
                Advanced
              </TabsTrigger>
            </TabsList>

            <TabsContent value="beginner" className="mt-8">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Zap className="mr-2 h-5 w-5 text-primary" />
                    Go Fundamentals
                  </CardTitle>
                  <CardDescription>
                    Master the basics of Go programming language
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-3">
                      <div className="flex items-center space-x-2">
                        <CheckCircle className="h-4 w-4 text-green-500" />
                        <span className="text-sm">Variables & Data Types</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <CheckCircle className="h-4 w-4 text-green-500" />
                        <span className="text-sm">Control Structures</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <CheckCircle className="h-4 w-4 text-green-500" />
                        <span className="text-sm">Functions & Methods</span>
                      </div>
                    </div>
                    <div className="space-y-3">
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Structs & Interfaces</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Error Handling</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Package Management</span>
                      </div>
                    </div>
                  </div>
                  <Progress value={75} className="mt-4" />
                  <p className="text-sm text-muted-foreground">8 lessons • 2-3 weeks</p>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="intermediate" className="mt-8">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Globe className="mr-2 h-5 w-5 text-primary" />
                    Web Development & APIs
                  </CardTitle>
                  <CardDescription>
                    Build web applications and REST APIs with Go
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-3">
                      <div className="flex items-center space-x-2">
                        <CheckCircle className="h-4 w-4 text-green-500" />
                        <span className="text-sm">HTTP Server Basics</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">REST API Design</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Database Integration</span>
                      </div>
                    </div>
                    <div className="space-y-3">
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Middleware & Auth</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Testing Strategies</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Deployment</span>
                      </div>
                    </div>
                  </div>
                  <Progress value={25} className="mt-4" />
                  <p className="text-sm text-muted-foreground">12 lessons • 4-5 weeks</p>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="advanced" className="mt-8">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <TrendingUp className="mr-2 h-5 w-5 text-primary" />
                    Microservices & Performance
                  </CardTitle>
                  <CardDescription>
                    Build scalable, production-ready Go applications
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-3">
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Concurrency Patterns</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Microservices Architecture</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">gRPC & Protocol Buffers</span>
                      </div>
                    </div>
                    <div className="space-y-3">
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Performance Optimization</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Monitoring & Observability</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-blue-500" />
                        <span className="text-sm">Production Deployment</span>
                      </div>
                    </div>
                  </div>
                  <Progress value={0} className="mt-4" />
                  <p className="text-sm text-muted-foreground">15 lessons • 6-8 weeks</p>
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="py-20 sm:py-32 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-background via-primary/5 to-background pointer-events-none" />

        <div className="container max-w-screen-2xl px-4 relative z-10">
          <div className="mx-auto max-w-2xl text-center mb-16">
            <Badge variant="outline" className="mb-4 animate-in fade-in duration-500">
              💬 Testimonials
            </Badge>
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl mb-4 animate-in fade-in slide-in-bottom duration-700">
              Loved by <span className="go-gradient-text">Developers Worldwide</span>
            </h2>
            <p className="text-lg text-muted-foreground animate-in fade-in slide-in-bottom duration-1000">
              Join thousands of developers who've advanced their careers with GO-PRO
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {testimonials.map((testimonial, index) => (
              <Card
                key={testimonial.name}
                className="glass-card group hover:border-primary/50 transition-all duration-500 animate-in fade-in scale-in relative overflow-hidden"
                style={{ animationDelay: `${index * 150}ms` }}
              >
                {/* Gradient overlay */}
                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-yellow-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                <CardHeader className="relative z-10">
                  <div className="flex items-center space-x-3 mb-4">
                    <Avatar className="h-12 w-12 ring-2 ring-primary/20 group-hover:ring-primary/50 group-hover:scale-110 transition-all duration-300 shadow-lg">
                      <AvatarImage src={testimonial.avatar} alt={testimonial.name} />
                      <AvatarFallback className="bg-gradient-to-br from-primary/20 to-primary/10 text-primary font-semibold">
                        {testimonial.name.split(' ').map(n => n[0]).join('')}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <CardTitle className="text-base group-hover:text-primary transition-colors duration-300">{testimonial.name}</CardTitle>
                      <CardDescription className="text-sm">{testimonial.role}</CardDescription>
                    </div>
                  </div>
                  <div className="flex space-x-1">
                    {[...Array(testimonial.rating)].map((_, i) => (
                      <Star
                        key={`star-${testimonial.name}-${i}`}
                        className="h-4 w-4 fill-yellow-400 text-yellow-400 group-hover:scale-125 transition-transform duration-300"
                        style={{ transitionDelay: `${i * 50}ms` }}
                      />
                    ))}
                  </div>
                </CardHeader>
                <CardContent className="relative z-10">
                  <p className="text-sm text-muted-foreground group-hover:text-foreground/80 italic leading-relaxed transition-colors duration-300">"{testimonial.content}"</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Enhanced CTA Section */}
      <section className="py-20 sm:py-32 relative overflow-hidden">
        <div className="absolute inset-0 animated-gradient" />

        {/* Enhanced Decorative elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-gradient-to-r from-primary/30 via-blue-500/20 to-cyan-500/30 rounded-full blur-3xl float-animation" />
          <div className="absolute top-1/4 right-1/4 w-64 h-64 bg-gradient-to-br from-yellow-500/20 to-orange-500/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '1s' }} />
          <div className="absolute bottom-1/4 left-1/4 w-64 h-64 bg-gradient-to-tr from-green-500/20 to-emerald-500/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '2s' }} />
        </div>

        <div className="container max-w-screen-2xl px-4 relative z-10">
          <div className="mx-auto max-w-3xl text-center glass-card-strong p-12 rounded-2xl shadow-2xl hover:shadow-primary/20 transition-all duration-500 border-2 border-primary/20 hover:border-primary/40">
            <Badge variant="secondary" className="mb-6 pulse-badge shadow-lg">
              🎯 Start Your Journey
            </Badge>
            <h2 className="text-3xl font-bold tracking-tight sm:text-5xl mb-6 animate-in fade-in slide-in-bottom duration-700">
              Ready to Master <span className="go-gradient-text">Go</span>?
            </h2>
            <p className="text-lg text-muted-foreground mb-8 leading-relaxed animate-in fade-in slide-in-bottom duration-1000">
              Join thousands of developers who are already building amazing things with Go.
              Start your journey today - it's completely free!
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-8">
              <Link href="/learn/1">
                <Button size="lg" className="go-gradient text-white text-lg px-8 py-6 shadow-2xl hover:shadow-primary/50">
                  <Play className="mr-2 h-5 w-5" />
                  Start Learning Now
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </Link>
              <Link href="/curriculum">
                <Button size="lg" variant="outline" className="text-lg px-8 py-6 border-2">
                  <BookOpen className="mr-2 h-5 w-5" />
                  View Curriculum
                </Button>
              </Link>
            </div>

            <div className="flex items-center justify-center gap-6 text-sm text-muted-foreground">
              <div className="flex items-center gap-2">
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span>No credit card required</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span>Free forever</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
