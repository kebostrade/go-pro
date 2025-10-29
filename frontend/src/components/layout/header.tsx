"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { 
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@/components/ui/navigation-menu";
import { cn } from "@/lib/utils";
import {
  BookOpen,
  Code2,
  Trophy,
  Users,
  Menu,
  Sun,
  Moon,
} from "lucide-react";
import { useState, useEffect } from "react";

const Header = () => {
  const [isDark, setIsDark] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    // Check for saved theme preference or default to system preference
    const savedTheme = localStorage.getItem('theme');
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

    if (savedTheme === 'dark' || (!savedTheme && systemPrefersDark)) {
      setIsDark(true);
      document.documentElement.classList.add('dark');
    }
  }, []);

  const toggleTheme = () => {
    const newTheme = !isDark;
    setIsDark(newTheme);

    if (newTheme) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  };

  const navigationItems = [
    {
      title: "Learn",
      href: "/learn",
      description: "Interactive Go lessons and tutorials",
      icon: BookOpen,
    },
    {
      title: "Practice",
      href: "/practice",
      description: "Coding exercises and challenges",
      icon: Code2,
    },
    {
      title: "Projects",
      href: "/projects",
      description: "Real-world Go applications",
      icon: Trophy,
    },
    {
      title: "Community",
      href: "/community",
      description: "Connect with other Go developers",
      icon: Users,
    },
  ];

  return (
    <header
      className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
      role="banner"
    >
      <div className="container flex h-14 sm:h-16 lg:h-18 max-w-screen-2xl items-center justify-between px-3 sm:px-4 lg:px-6 xl:px-8">
        {/* Enhanced Logo */}
        <Link href="/" className="flex items-center space-x-2 sm:space-x-3 group">
          <div className="flex h-7 w-7 sm:h-8 sm:w-8 lg:h-9 lg:w-9 items-center justify-center rounded-lg bg-gradient-to-br from-primary to-primary/80 shadow-lg group-hover:shadow-xl group-hover:shadow-primary/30 transition-all duration-300 group-hover:scale-110 group-hover:rotate-3">
            <span className="text-sm sm:text-lg lg:text-xl font-bold text-primary-foreground">G</span>
          </div>
          <div className="flex flex-col">
            <span className="text-base sm:text-lg lg:text-xl font-bold go-gradient-text group-hover:scale-105 transition-transform duration-300">GO-PRO</span>
            <span className="hidden sm:block text-xs lg:text-sm text-muted-foreground -mt-1 group-hover:text-foreground/80 transition-colors">Learn Go Programming</span>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <NavigationMenu className="hidden lg:flex">
          <NavigationMenuList className="space-x-1">
            {navigationItems.map((item) => (
              <NavigationMenuItem key={item.title}>
                <NavigationMenuTrigger className="h-9 px-3 lg:px-4 py-2 text-sm lg:text-base">
                  <item.icon className="mr-2 h-3 w-3 lg:h-4 lg:w-4" />
                  {item.title}
                </NavigationMenuTrigger>
                <NavigationMenuContent>
                  <div className="grid gap-3 p-4 lg:p-6 w-[350px] lg:w-[400px]">
                    <NavigationMenuLink asChild>
                      <Link
                        href={item.href}
                        className={cn(
                          "block select-none space-y-1 rounded-md p-3 leading-none no-underline outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground"
                        )}
                      >
                        <div className="flex items-center space-x-2">
                          <item.icon className="h-4 w-4" />
                          <div className="text-sm font-medium leading-none">{item.title}</div>
                        </div>
                        <p className="line-clamp-2 text-sm leading-snug text-muted-foreground">
                          {item.description}
                        </p>
                      </Link>
                    </NavigationMenuLink>
                  </div>
                </NavigationMenuContent>
              </NavigationMenuItem>
            ))}
          </NavigationMenuList>
        </NavigationMenu>

        {/* Medium Screen Navigation - Simplified */}
        <nav className="hidden md:flex lg:hidden">
          <div className="flex items-center space-x-1">
            {navigationItems.map((item) => (
              <Link
                key={item.title}
                href={item.href}
                className="flex items-center space-x-1 px-3 py-2 text-sm font-medium rounded-md hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                <item.icon className="h-4 w-4" />
                <span>{item.title}</span>
              </Link>
            ))}
          </div>
        </nav>

        {/* Right side actions */}
        <div className="flex items-center space-x-1 sm:space-x-2">
          {/* Enhanced Theme toggle */}
          {mounted && (
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleTheme}
              className="h-8 w-8 sm:h-9 sm:w-9 relative overflow-hidden group hover:bg-primary/10 transition-all duration-300"
            >
              <div className="absolute inset-0 bg-gradient-to-br from-primary/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
              {isDark ? (
                <Sun className="h-3 w-3 sm:h-4 sm:w-4 text-yellow-500 group-hover:rotate-90 transition-transform duration-500 relative z-10" />
              ) : (
                <Moon className="h-3 w-3 sm:h-4 sm:w-4 text-blue-600 dark:text-blue-400 group-hover:-rotate-12 transition-transform duration-500 relative z-10" />
              )}
              <span className="sr-only">Toggle theme</span>
            </Button>
          )}

          {/* CTA Section - responsive visibility */}
          <div className="hidden sm:flex items-center space-x-2 ml-2">
            <Badge variant="secondary" className="text-xs lg:text-sm">
              Beta
            </Badge>
            <Button size="sm" className="go-gradient text-white text-xs sm:text-sm px-2 sm:px-3 lg:px-4">
              <span className="hidden sm:inline">Get Started</span>
              <span className="sm:hidden">Start</span>
            </Button>
          </div>

          {/* Mobile menu button */}
          <Button
            variant="ghost"
            size="icon"
            className="md:flex lg:hidden h-8 w-8 sm:h-9 sm:w-9 ml-1"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-expanded={isMobileMenuOpen}
            aria-controls="mobile-menu"
            aria-label={isMobileMenuOpen ? "Close menu" : "Open menu"}
          >
            <Menu className="h-4 w-4" />
            <span className="sr-only">Toggle menu</span>
          </Button>
        </div>
      </div>

      {/* Mobile Navigation */}
      {isMobileMenuOpen && (
        <nav
          id="mobile-menu"
          className="md:flex lg:hidden border-t border-border bg-background/95 backdrop-blur shadow-lg"
          role="navigation"
          aria-label="Mobile navigation"
        >
          <div className="container px-3 sm:px-4 py-3 sm:py-4 space-y-2 sm:space-y-3 max-w-screen-2xl">
            {navigationItems.map((item) => (
              <Link
                key={item.title}
                href={item.href}
                className="flex items-center space-x-3 rounded-lg p-3 text-sm sm:text-base hover:bg-accent hover:text-accent-foreground transition-colors active:bg-accent/80"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                <item.icon className="h-4 w-4 sm:h-5 sm:w-5 flex-shrink-0" />
                <div className="flex-1 min-w-0">
                  <div className="font-medium truncate">{item.title}</div>
                  <div className="text-xs sm:text-sm text-muted-foreground line-clamp-1">{item.description}</div>
                </div>
              </Link>
            ))}

            {/* Mobile CTA Section */}
            <div className="pt-3 border-t border-border">
              <Button size="sm" className="w-full go-gradient text-white">
                Get Started
              </Button>
            </div>
          </div>
        </nav>
      )}
    </header>
  );
};

export default Header;
