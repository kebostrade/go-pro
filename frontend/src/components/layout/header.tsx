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
  GraduationCap,
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

  // Close mobile menu on Escape key press
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isMobileMenuOpen) {
        setIsMobileMenuOpen(false);
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isMobileMenuOpen]);

  // Prevent body scroll when mobile menu is open
  useEffect(() => {
    if (isMobileMenuOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isMobileMenuOpen]);

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
      title: "Tutorials",
      href: "/tutorials",
      description: "19 comprehensive tutorials from basics to advanced",
      icon: GraduationCap,
      badge: "19 Tutorials",
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
      <div className="container flex h-14 sm:h-16 lg:h-18 max-w-screen-2xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Enhanced Logo - Responsive sizing */}
        <Link href="/" className="flex items-center space-x-2 sm:space-x-3 group flex-shrink-0">
          <div className="flex h-8 w-8 sm:h-9 sm:w-9 lg:h-10 lg:w-10 items-center justify-center rounded-lg bg-gradient-to-br from-primary to-primary/80 shadow-lg group-hover:shadow-xl group-hover:shadow-primary/30 transition-all duration-300 group-hover:scale-110 group-hover:rotate-3">
            <span className="text-base sm:text-lg lg:text-xl font-bold text-primary-foreground">G</span>
          </div>
          <div className="flex flex-col">
            <span className="text-base sm:text-lg lg:text-xl font-bold go-gradient-text group-hover:scale-105 transition-transform duration-300 whitespace-nowrap">GO-PRO</span>
            <span className="hidden sm:block text-[10px] sm:text-xs lg:text-sm text-muted-foreground -mt-1 group-hover:text-foreground/80 transition-colors whitespace-nowrap">Learn Go Programming</span>
          </div>
        </Link>

        {/* Desktop Navigation - Full menu */}
        <NavigationMenu className="hidden lg:flex">
          <NavigationMenuList className="gap-1">
            {navigationItems.map((item) => (
              <NavigationMenuItem key={item.title}>
                <NavigationMenuTrigger className="h-9 px-3 xl:px-4 text-sm font-medium">
                  <item.icon className="mr-1.5 h-4 w-4" />
                  {item.title}
                </NavigationMenuTrigger>
                <NavigationMenuContent>
                  <div className="grid gap-3 p-5 w-[380px]">
                    <NavigationMenuLink asChild>
                      <Link
                        href={item.href}
                        className={cn(
                          "block select-none space-y-2 rounded-lg p-3 leading-none no-underline outline-none transition-all hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground hover:shadow-sm"
                        )}
                      >
                        <div className="flex items-center gap-2">
                          <item.icon className="h-5 w-5 flex-shrink-0" />
                          <div className="text-sm font-semibold">{item.title}</div>
                          {item.badge && (
                            <Badge variant="secondary" className="text-xs ml-auto">
                              {item.badge}
                            </Badge>
                          )}
                        </div>
                        <p className="text-sm leading-relaxed text-muted-foreground">
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

        {/* Tablet Navigation - Compact links */}
        <nav className="hidden md:flex lg:hidden" aria-label="Tablet navigation">
          <div className="flex items-center space-x-1 overflow-x-auto">
            {navigationItems.map((item) => (
              <Link
                key={item.title}
                href={item.href}
                className="flex items-center space-x-1.5 px-2.5 py-2 text-sm font-medium rounded-md hover:bg-accent hover:text-accent-foreground transition-colors whitespace-nowrap"
                title={item.description}
              >
                <item.icon className="h-4 w-4 flex-shrink-0" />
                <span className="hidden md:inline">{item.title}</span>
              </Link>
            ))}
          </div>
        </nav>

        {/* Right side actions */}
        <div className="flex items-center gap-1 sm:gap-2 flex-shrink-0">
          {/* Enhanced Theme toggle */}
          {mounted && (
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleTheme}
              className="h-9 w-9 sm:h-10 sm:w-10 relative overflow-hidden group hover:bg-primary/10 transition-all duration-300"
              aria-label="Toggle theme"
            >
              <div className="absolute inset-0 bg-gradient-to-br from-primary/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
              {isDark ? (
                <Sun className="h-4 w-4 sm:h-[18px] sm:w-[18px] text-yellow-500 group-hover:rotate-90 transition-transform duration-500 relative z-10" />
              ) : (
                <Moon className="h-4 w-4 sm:h-[18px] sm:w-[18px] text-blue-600 dark:text-blue-400 group-hover:-rotate-12 transition-transform duration-500 relative z-10" />
              )}
              <span className="sr-only">Toggle theme</span>
            </Button>
          )}

          {/* CTA Section - Responsive visibility and sizing */}
          <div className="hidden sm:flex items-center gap-2 ml-1 md:ml-2">
            <Badge variant="secondary" className="hidden md:flex text-xs whitespace-nowrap">
              Beta
            </Badge>
            <Button
              size="sm"
              className="go-gradient text-white text-xs sm:text-sm px-3 sm:px-4 h-9 whitespace-nowrap"
            >
              Get Started
            </Button>
          </div>

          {/* Mobile menu button - Shows on mobile and small tablets */}
          <Button
            variant="ghost"
            size="icon"
            className="flex md:hidden h-8 w-8 sm:h-9 sm:w-9 ml-1"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-expanded={isMobileMenuOpen}
            aria-controls="mobile-menu"
            aria-label={isMobileMenuOpen ? "Close menu" : "Open menu"}
          >
            <Menu className="h-4 w-4 sm:h-5 sm:w-5" />
            <span className="sr-only">Toggle menu</span>
          </Button>
        </div>
      </div>

      {/* Mobile Menu Overlay */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/20 backdrop-blur-sm md:hidden animate-in fade-in duration-200"
          onClick={() => setIsMobileMenuOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Mobile Navigation - Shows only on mobile and small tablets */}
      {isMobileMenuOpen && (
        <nav
          id="mobile-menu"
          className="relative z-50 block md:hidden border-t border-border bg-background/98 backdrop-blur-md shadow-xl animate-in slide-in-from-top-2 duration-200"
          role="navigation"
          aria-label="Mobile navigation"
        >
          <div className="container px-4 py-4 space-y-2 max-w-screen-2xl">
            {navigationItems.map((item) => (
              <Link
                key={item.title}
                href={item.href}
                className="flex items-start space-x-3 rounded-lg p-3 text-sm hover:bg-accent hover:text-accent-foreground transition-colors active:scale-98 active:bg-accent/80"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                <item.icon className="h-5 w-5 flex-shrink-0 mt-0.5" />
                <div className="flex-1 min-w-0 space-y-0.5">
                  <div className="font-semibold">{item.title}</div>
                  <div className="text-xs text-muted-foreground leading-relaxed">
                    {item.description}
                  </div>
                  {item.badge && (
                    <Badge variant="secondary" className="text-xs mt-1">
                      {item.badge}
                    </Badge>
                  )}
                </div>
              </Link>
            ))}

            {/* Mobile CTA Section */}
            <div className="pt-3 mt-2 border-t border-border space-y-2">
              <Button size="default" className="w-full go-gradient text-white font-semibold">
                Get Started
              </Button>
              <div className="flex items-center justify-center">
                <Badge variant="secondary" className="text-xs">
                  Beta
                </Badge>
              </div>
            </div>
          </div>
        </nav>
      )}
    </header>
  );
};

export default Header;
