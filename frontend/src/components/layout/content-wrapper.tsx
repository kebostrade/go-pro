"use client";

import { cn } from "@/lib/utils";

interface ContentWrapperProps {
  children: React.ReactNode;
  className?: string;
  layout?: "single" | "sidebar" | "two-column" | "three-column";
  gap?: "sm" | "md" | "lg" | "xl";
  align?: "start" | "center" | "stretch";
}

const ContentWrapper = ({
  children,
  className,
  layout = "single",
  gap = "lg",
  align = "stretch"
}: ContentWrapperProps) => {
  const layoutClasses = {
    single: "grid grid-cols-1",
    sidebar: "grid grid-cols-1 lg:grid-cols-4",
    "two-column": "grid grid-cols-1 lg:grid-cols-2",
    "three-column": "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3"
  };

  const gapClasses = {
    sm: "gap-4",
    md: "gap-6",
    lg: "gap-8",
    xl: "gap-12"
  };

  const alignClasses = {
    start: "items-start",
    center: "items-center",
    stretch: "items-stretch"
  };

  return (
    <div className={cn(
      layoutClasses[layout],
      gapClasses[gap],
      alignClasses[align],
      className
    )}>
      {children}
    </div>
  );
};

// Specialized content area components
export const MainContent = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <div className={cn("lg:col-span-3", className)}>
    {children}
  </div>
);

export const Sidebar = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <div className={cn("lg:col-span-1", className)}>
    {children}
  </div>
);

export const TwoColumnLeft = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <div className={cn("lg:col-span-1", className)}>
    {children}
  </div>
);

export const TwoColumnRight = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <div className={cn("lg:col-span-1", className)}>
    {children}
  </div>
);

export default ContentWrapper;
