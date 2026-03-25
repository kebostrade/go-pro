"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight, BookOpen, Bookmark, Circle, Cloud, Code, Database, DollarSign, Activity, GitBranch, HardHat, Lock, MapPin, Monitor, Server, Settings, Shield, ShieldCheck, Terminal, TrendingUp, Users, Zap } from "lucide-react";
import Link from "next/link";
import MarkdownRenderer from "@/components/learning/markdown-renderer";

export default function DevOpsEngineeringPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container mx-auto px-4 py-8">
        {/* Breadcrumb Navigation */}
        <div className="mb-6">
          <div className="flex items-center space-x-2 text-sm text-muted-foreground mb-4">
            <Link href="/" className="hover:text-primary transition-colors">
              <ArrowLeft className="h-4 w-4" />
            </Link>
            <span className="mx-2">/</span>
            <Link href="/learn" className="hover:text-primary transition-colors">
              Curriculum
            </Link>
            <span className="mx-2">/</span>
            <span className="text-foreground font-medium">DevOps Engineering</span>
          </div>
        </div>

        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent mb-4">
            DevOps Engineering
          </h1>
          <p className="text-muted-foreground text-lg max-w-3xl">
            Master the practices and tools that bridge development and operations to deliver software faster and more reliably.
          </p>
        </div>

        {/* Main Content */}
        <div className="space-y-8">
          {/* Overview Section */}
          <Card className="border-2">
            <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <BookOpen className="h-5 w-5 text-primary" />
                </div>
                Overview
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
               <MarkdownRenderer
                 content={`# DevOps Engineering
 
 DevOps Engineering is a set of practices that combines software development (Dev) and IT operations (Ops) to shorten the systems development life cycle and provide continuous delivery with high software quality. It emphasizes collaboration, automation, and integration between developers and operations teams.
 
 ## Core Principles
 
 - **Collaboration**: Breaking down silos between development and operations
 - **Automation**: Automating repetitive tasks in the software delivery process
 - **Continuous Improvement**: Constantly measuring and improving processes
 - **Customer Focus**: Delivering value to customers quickly and reliably
 - **End-to-End Responsibility**: Teams own the entire lifecycle of applications
 
 ## Key Practices
 
 - Continuous Integration and Continuous Deployment (CI/CD)
 - Infrastructure as Code (IaC)
 - Monitoring and Logging
 - Microservices Architecture
 - Containerization
 - Cloud Computing
 - Site Reliability Engineering (SRE) principles`}
                 enableCodeHighlight={true}
                 enableInteractiveExamples={true}
               />
            </CardContent>
          </Card>

          {/* Skills Section */}
          <Card className="border-2">
            <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <Code className="h-5 w-5 text-primary" />
                </div>
                Essential Skills
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              <div className="grid gap-6 md:grid-cols-2">
                <div>
                  <h3 className="text-lg font-semibold mb-4">Technical Skills</h3>
                  <ul className="space-y-2 text-sm">
                    <li className="flex items-center">
                      <Zap className="h-4 w-4 text-primary mr-2" />
                      <span>Linux/Unix Systems Administration</span>
                    </li>
                    <li className="flex items-center">
                      <Terminal className="h-4 w-4 text-primary mr-2" />
                      <span>Shell Scripting (Bash, PowerShell)</span>
                    </li>
                    <li className="flex items-center">
                      <Server className="h-4 w-4 text-primary mr-2" />
                      <span>Networking & Protocols</span>
                    </li>
                    <li className="flex items-center">
                      <Cloud className="h-4 w-4 text-primary mr-2" />
                      <span>Cloud Platforms (AWS, Azure, GCP)</span>
                    </li>
                    <li className="flex items-center">
                      <HardHat className="h-4 w-4 text-primary mr-2" />
                      <span>Infrastructure as Code</span>
                    </li>
                    <li className="flex items-center">
                      <Database className="h-4 w-4 text-primary mr-2" />
                      <span>Configuration Management</span>
                    </li>
                    <li className="flex items-center">
                      <Activity className="h-4 w-4 text-primary mr-2" />
                      <span>Monitoring & Logging</span>
                    </li>
                    <li className="flex items-center">
                      <Shield className="h-4 w-4 text-primary mr-2" />
                      <span>Security Best Practices</span>
                    </li>
                  </ul>
                </div>
                 <div>
                   <h3 className="text-lg font-semibold mb-4">Tools & Technologies</h3>
                   <ul className="space-y-2 text-sm">
                     <li className="flex items-center">
                       <GitBranch className="h-4 w-4 text-primary mr-2" />
                       <span>Version Control (Git)</span>
                     </li>
                     <li className="flex items-center">
                       <Monitor className="h-4 w-4 text-primary mr-2" />
                       <span>CI/CD Tools (Jenkins, GitLab CI, GitHub Actions)</span>
                     </li>
                     <li className="flex items-center">
                       <Users className="h-4 w-4 text-primary mr-2" />
                       <span>Containerization (Docker, Kubernetes)</span>
                     </li>
                     <li className="flex items-center">
                       <DollarSign className="h-4 w-4 text-primary mr-2" />
                       <span>Cloud Services (EC2, S3, Lambda, etc.)</span>
                     </li>
                     <li className="flex items-center">
                       <TrendingUp className="h-4 w-4 text-primary mr-2" />
                       <span>Infrastructure Tools (Terraform, Ansible, Chef)</span>
                     </li>
                     <li className="flex items-center">
                       <Code className="h-4 w-4 text-primary mr-2" />
                       <span>Programming & Scripting</span>
                     </li>
                     <li className="flex items-center">
                       <Lock className="h-4 w-4 text-primary mr-2" />
                       <span>Security & Compliance Tools</span>
                     </li>
                     <li className="flex items-center">
                       <Settings className="h-4 w-4 text-primary mr-2" />
                       <span>Observability (Prometheus, Grafana, ELK)</span>
                     </li>
                   </ul>
                 </div>
              </div>
            </CardContent>
          </Card>

          {/* Learning Path Section */}
          <Card className="border-2">
            <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <MapPin className="h-5 w-5 text-primary" />
                </div>
                Suggested Learning Path
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              <div className="space-y-6">
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Foundations</h4>
                    <p className="text-sm text-muted-foreground">
                      Start with Linux, networking, and networking fundamentals. Learn basic shell scripting and version control with Git.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Automation</h4>
                    <p className="text-sm text-muted-foreground">
                      Learn configuration management tools (Ansible, Puppet, Chef) and basic scripting for automation.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Cloud & Containers</h4>
                    <p className="text-sm text-muted-foreground">
                      Master Docker, Kubernetes, and at least one major cloud platform (AWS, Azure, or GCP).
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">CI/CD & Monitoring</h4>
                    <p className="text-sm text-muted-foreground">
                      Study CI/CD pipelines, infrastructure as code, and monitoring/logging solutions.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Specialization</h4>
                    <p className="text-sm text-muted-foreground">
                      Focus on areas like site reliability engineering, cloud architecture, or DevSecOps.
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Resources Section */}
          <Card className="border-2">
            <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <BookOpen className="h-5 w-5 text-primary" />
                </div>
                Recommended Resources
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              <div className="space-y-4">
                <div className="flex items-start space-x-3">
                  <Bookmark className="h-4 w-4 text-primary mt-1 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Books</h4>
                    <p className="text-sm text-muted-foreground">
                      "The Phoenix Project" by Gene Kim, "Accelerate" by Nicole Forsgren, "Site Reliability Engineering" by Google
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-3">
                  <Monitor className="h-4 w-4 text-primary mt-1 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Courses</h4>
                    <p className="text-sm text-muted-foreground">
                      AWS DevOps Engineer Professional, Google Cloud DevOps Engineer, Docker and Kubernetes Fundamentals
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-3">
                  <Server className="h-4 w-4 text-primary mt-1 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Tools to Explore</h4>
                    <p className="text-sm text-muted-foreground">
                      Git, Docker, Kubernetes, Jenkins, Terraform, Ansible, Prometheus, Grafana, AWS/Azure/GCP
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Navigation Buttons */}
        <div className="mt-10 flex flex-col sm:flex-row gap-4">
          <Link href="/learn/ai-platform-engineering" className="flex-1 sm:flex-none">
            <Button variant="outline" size="lg">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Previous: AI Platform Engineering
            </Button>
          </Link>
          <Link href="/learn" className="flex-1 sm:flex-none">
            <Button size="lg">
              Back to Curriculum
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}


