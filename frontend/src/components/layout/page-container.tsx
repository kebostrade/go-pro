"use client";

import { cn } from "@/lib/utils";

interface PageContainerProps {
  children: React.ReactNode;
  className?: string;
  size?: "default" | "narrow" | "wide" | "full";
  padding?: "default" | "compact" | "spacious" | "none";
  background?: "default" | "muted" | "accent" | "gradient";
}

const PageContainer = ({
  children,
  className,
  size = "default",
  padding = "default",
  background = "default"
}: PageContainerProps) => {
  const sizeClasses = {
    narrow: "max-w-4xl",
    default: "max-w-7xl",
    wide: "max-w-screen-2xl",
    full: "max-w-none"
  };

  const paddingClasses = {
    none: "",
    compact: "px-4 py-4 sm:px-6 sm:py-6",
    default: "px-4 py-6 sm:px-6 sm:py-8 lg:px-8 lg:py-10",
    spacious: "px-4 py-8 sm:px-6 sm:py-12 lg:px-8 lg:py-16"
  };

  const backgroundClasses = {
    default: "",
    muted: "bg-muted/30",
    accent: "bg-accent/5",
    gradient: "bg-gradient-to-br from-background via-background to-accent/5"
  };

  return (
    <div className={cn(backgroundClasses[background], "min-h-full")}>
      <div className={cn(
        "container mx-auto",
        sizeClasses[size],
        paddingClasses[padding],
        className
      )}>
        {children}
      </div>
    </div>
  );
};

export default PageContainer;
