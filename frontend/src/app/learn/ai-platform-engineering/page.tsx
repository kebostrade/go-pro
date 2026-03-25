"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight, BookOpen, Bookmark, Brain, Circle, Code, Database, Monitor, Server, Settings, Terminal, ShieldCheck, Zap, Activity, GitBranch, Cloud, Lock, MapPin, Users } from "lucide-react";
import Link from "next/link";
import MarkdownRenderer from "@/components/learning/markdown-renderer";

export default function AIPlatformEngineeringPage() {
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
            <span className="text-foreground font-medium">AI Platform Engineering</span>
          </div>
        </div>

        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent mb-4">
            AI Platform Engineering
          </h1>
          <p className="text-muted-foreground text-lg max-w-3xl">
            Learn how to build, deploy, and maintain scalable AI platforms that power modern applications.
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
                 content={`# AI Platform Engineering

 AI Platform Engineering focuses on building the infrastructure, tools, and systems that enable organizations to develop, deploy, and manage AI applications at scale. This discipline combines software engineering, DevOps, data engineering, and machine learning expertise to create robust AI platforms.

 ## Key Responsibilities

 - Designing scalable AI infrastructure
 - Building MLOps pipelines for model training and deployment
 - Creating developer tools and APIs for AI services
 - Ensuring reliability, monitoring, and observability of AI systems
 - Implementing security and compliance measures for AI workloads
 - Optimizing performance and cost of AI operations

 ## Core Technologies

 AI Platform Engineers work with a variety of technologies including:

 - Container orchestration (Kubernetes, Docker)
 - Cloud platforms (AWS, GCP, Azure)
 - Machine learning frameworks (TensorFlow, PyTorch, Scikit-learn)
 - MLflow, Kubeflow, and other MLOps tools
 - APIs and microservices architecture
 - Monitoring and logging systems (Prometheus, Grafana, ELK)
 - Infrastructure as Code (Terraform, CloudFormation)`}
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
                      <span>Cloud Computing & Infrastructure</span>
                    </li>
                    <li className="flex items-center">
                      <Terminal className="h-4 w-4 text-primary mr-2" />
                      <span>Linux/Unix Systems Administration</span>
                    </li>
                    <li className="flex items-center">
                      <Server className="h-4 w-4 text-primary mr-2" />
                      <span>Containerization & Orchestration</span>
                    </li>
                    <li className="flex items-center">
                      <Database className="h-4 w-4 text-primary mr-2" />
                      <span>Data Engineering & ETL Pipelines</span>
                    </li>
                    <li className="flex items-center">
                      <ShieldCheck className="h-4 w-4 text-primary mr-2" />
                      <span>Security & Compliance</span>
                    </li>
                    <li className="flex items-center">
                      <Activity className="h-4 w-4 text-primary mr-2" />
                      <span>Monitoring & Observability</span>
                    </li>
                  </ul>
                </div>
                <div>
                  <h3 className="text-lg font-semibold mb-4">AI/ML Specific Skills</h3>
                  <ul className="space-y-2 text-sm">
                    <li className="flex items-center">
                      <Brain className="h-4 w-4 text-primary mr-2" />
                      <span>Machine Learning Fundamentals</span>
                    </li>
                    <li className="flex items-center">
                      <GitBranch className="h-4 w-4 text-primary mr-2" />
                      <span>MLOps & Model Versioning</span>
                    </li>
                    <li className="flex items-center">
                      <Users className="h-4 w-4 text-primary mr-2" />
                      <span>Collaboration with Data Scientists</span>
                    </li>
                    <li className="flex items-center">
                      <Code className="h-4 w-4 text-primary mr-2" />
                      <span>Programming (Python, Go, Java)</span>
                    </li>
                    <li className="flex items-center">
                      <Cloud className="h-4 w-4 text-primary mr-2" />
                      <span>Cloud AI Services</span>
                    </li>
                    <li className="flex items-center">
                      <Lock className="h-4 w-4 text-primary mr-2" />
                      <span>AI Ethics & Governance</span>
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
                      Start with Linux, networking, and cloud fundamentals. Learn Python and basic shell scripting.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Infrastructure</h4>
                    <p className="text-sm text-muted-foreground">
                      Master Docker, Kubernetes, and cloud platforms. Learn Infrastructure as Code with Terraform.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">MLOps</h4>
                    <p className="text-sm text-muted-foreground">
                      Study model lifecycle management, CI/CD for ML, and experiment tracking tools.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <Circle className="h-4 w-4 text-primary mt-2.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Specialization</h4>
                    <p className="text-sm text-muted-foreground">
                      Focus on areas like AI infrastructure, model serving, or AI platform product development.
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
                      "Machine Learning Engineering" by Andriy Burkov, "Building Machine Learning Powered Applications" by Emmanuel Ameisen
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-3">
                  <Monitor className="h-4 w-4 text-primary mt-1 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Courses</h4>
                    <p className="text-sm text-muted-foreground">
                      Coursera MLOps Specialization, AWS Machine Learning Specialty, Google Cloud Professional ML Engineer
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-3">
                  <Server className="h-4 w-4 text-primary mt-1 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium">Tools to Explore</h4>
                    <p className="text-sm text-muted-foreground">
                      Kubernetes, Kubeflow, MLflow, Prometheus, Grafana, Airflow, Terraform
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Navigation Buttons */}
        <div className="mt-10 flex flex-col sm:flex-row gap-4">
          <Link href="/learn" className="flex-1 sm:flex-none">
            <Button variant="outline" size="lg">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Curriculum
            </Button>
          </Link>
          <Link href="/learn/devops-engineering" className="flex-1 sm:flex-none">
            <Button size="lg">
              Next: DevOps Engineering →
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}


