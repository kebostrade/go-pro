"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
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
  X,
  ChevronRight,
  LogIn,
  LogOut,
  Settings,
  UserCircle,
} from "lucide-react";
import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/auth-context";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { UserCheck } from "lucide-react";

const Header = () => {
  const { user, backendUser, signOut, loading } = useAuth();
  const [isDark, setIsDark] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  const handleSignOut = async () => {
    try {
      await signOut();
      setIsMobileMenuOpen(false);
    } catch (error) {
      console.error('Sign out error:', error);
    }
  };

  useEffect(() => {
    setMounted(true);
    const savedTheme = localStorage.getItem('theme');
    const systemPrefersDark = globalThis.matchMedia('(prefers-color-scheme: dark)').matches;

    if (savedTheme === 'dark' || (!savedTheme && systemPrefersDark)) {
      setIsDark(true);
      document.documentElement.classList.add('dark');
    }
  }, []);

  // Handle scroll for dynamic header
  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Close mobile menu on Escape key
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
      title: "Practice",
      href: "/practice",
      description: "Coding exercises and challenges",
      icon: Code2,
    },
    {
      title: "Tutorials",
      href: "/tutorials",
      description: "19 comprehensive tutorials from basics to advanced",
      icon: GraduationCap,
      badge: "19 Tutorials",
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
    {
      title: "Customers",
      href: "/customers",
      description: "Manage your customers",
      icon: UserCheck,
    },
  ];

  return (
    <header
      className={cn(
        "sticky top-0 z-50 w-full border-b transition-all duration-300",
        scrolled
          ? "border-border bg-background/80 backdrop-blur-lg shadow-sm"
          : "border-border/40 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60"
      )}
      role="banner"
    >
      <div className="container flex h-14 sm:h-16 lg:h-18 max-w-screen-2xl items-center justify-between px-3 sm:px-4 md:px-6 lg:px-8">
        {/* Enhanced Logo - Fully Responsive */}
        <Link href="/" className="flex items-center space-x-2 sm:space-x-3 group shrink-0 min-w-0">
          <div className="flex h-7 w-7 sm:h-8 sm:w-8 lg:h-9 lg:w-9 items-center justify-center rounded-lg bg-linear-to-br from-primary to-primary/80 shadow-lg group-hover:shadow-xl group-hover:shadow-primary/30 transition-all duration-300 group-hover:scale-110 group-hover:rotate-3 shrink-0">
            <span className="text-sm sm:text-lg lg:text-xl font-bold text-primary-foreground">G</span>
          </div>
          <div className="flex flex-col">
            <span className="text-base sm:text-lg lg:text-xl font-bold go-gradient-text group-hover:scale-105 transition-transform duration-300 whitespace-nowrap">GO-PRO</span>
            <span className="hidden sm:block text-xs lg:text-sm text-muted-foreground -mt-1 group-hover:text-foreground/80 transition-colors whitespace-nowrap">Learn Go Programming</span>
          </div>
        </Link>

        {/* Desktop Navigation - Clean and modern */}
        <nav className="hidden lg:flex items-center gap-1" aria-label="Main navigation">
          {navigationItems.map((item) => (
            <Link
              key={item.title}
              href={item.href}
              className={cn(
                "group relative flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200",
                "hover:bg-accent hover:text-accent-foreground",
                "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              )}
            >
              <item.icon className="h-4 w-4 transition-transform group-hover:scale-110" />
              <span>{item.title}</span>
              {item.badge && (
                <Badge variant="secondary" className="text-xs ml-1 px-1.5 py-0">
                  New
                </Badge>
              )}
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-linear-to-r from-blue-600 to-cyan-600 scale-x-0 group-hover:scale-x-100 transition-transform origin-left" />
            </Link>
          ))}
        </nav>

        {/* Tablet Navigation - Simplified */}
        <nav className="hidden md:flex lg:hidden items-center gap-0.5" aria-label="Tablet navigation">
          {navigationItems.slice(0, 3).map((item) => (
            <Link
              key={item.title}
              href={item.href}
              className="flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-md hover:bg-accent hover:text-accent-foreground transition-colors"
              title={item.description}
            >
              <item.icon className="h-4 w-4" />
              <span className="hidden md:inline">{item.title}</span>
            </Link>
          ))}
        </nav>

        {/* Right Actions */}
        <div className="flex items-center gap-2 md:gap-3 ml-auto">
          {/* Theme Toggle - Enhanced */}
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleTheme}
            className={cn(
              "h-9 w-9 md:h-10 md:w-10 rounded-lg relative overflow-hidden",
              "hover:bg-accent transition-all duration-300"
            )}
            aria-label="Toggle theme"
            suppressHydrationWarning
          >
            <div className="absolute inset-0 bg-linear-to-br from-yellow-400/20 to-orange-400/20 dark:from-blue-400/20 dark:to-purple-400/20 opacity-0 hover:opacity-100 transition-opacity" />
            {mounted && isDark ? (
              <Sun className="h-5 w-5 text-yellow-500 hover:rotate-90 transition-transform duration-500 relative z-10" />
            ) : (
              <Moon className="h-5 w-5 text-blue-600 dark:text-blue-400 hover:-rotate-12 transition-transform duration-500 relative z-10" />
            )}
          </Button>

          {/* Auth Section - Responsive */}
          {mounted && !loading && (
            <div className="flex items-center gap-2">
              {user ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      className="relative h-8 w-8 sm:h-9 sm:w-9 md:h-10 md:w-10 rounded-full"
                    >
                      <Avatar className="h-7 w-7 sm:h-8 sm:w-8 md:h-9 md:w-9">
                        <AvatarImage src={user.photoURL || undefined} alt={user.displayName || 'User'} />
                        <AvatarFallback className="bg-linear-to-br from-primary to-primary/80 text-xs sm:text-sm">
                          {user.displayName?.charAt(0)?.toUpperCase() || user.email?.charAt(0)?.toUpperCase() || 'U'}
                        </AvatarFallback>
                      </Avatar>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="w-52 sm:w-56" align="end" forceMount>
                    <DropdownMenuLabel className="font-normal">
                      <div className="flex flex-col space-y-1">
                        <p className="text-sm font-medium leading-none truncate">
                          {user.displayName || 'User'}
                        </p>
                        <p className="text-xs leading-none text-muted-foreground truncate">
                          {user.email}
                        </p>
                        {backendUser?.role && (
                          <Badge variant="secondary" className="w-fit mt-1 text-xs">
                            {backendUser.role}
                          </Badge>
                        )}
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem asChild>
                      <Link href="/profile" className="cursor-pointer">
                        <UserCircle className="mr-2 h-4 w-4" />
                        Profile
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link href="/settings" className="cursor-pointer">
                        <Settings className="mr-2 h-4 w-4" />
                        Settings
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      onClick={handleSignOut}
                      className="cursor-pointer text-destructive focus:text-destructive"
                    >
                      <LogOut className="mr-2 h-4 w-4" />
                      Sign out
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-8 px-3 sm:h-9 sm:px-4 md:h-10 md:px-6 text-xs sm:text-sm lg:text-base"
                    >
                      <LogIn className="mr-1.5 sm:mr-2 h-3.5 w-3.5 sm:h-4 sm:w-4 lg:h-5 lg:w-5" />
                      <span>Sign in</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="w-44 sm:w-48" align="end">
                    <DropdownMenuItem asChild>
                      <Link href="/signin" className="cursor-pointer">
                        <LogIn className="mr-2 h-4 w-4" />
                        Sign in
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem asChild>
                      <Link href="/profile" className="cursor-pointer">
                        <UserCircle className="mr-2 h-4 w-4" />
                        Profile
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link href="/settings" className="cursor-pointer">
                        <Settings className="mr-2 h-4 w-4" />
                        Settings
                      </Link>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </div>
          )}

          {/* Mobile Menu Button */}
          <Button
            variant="ghost"
            size="icon"
            className="flex md:hidden h-9 w-9 rounded-lg"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-expanded={isMobileMenuOpen}
            aria-controls="mobile-menu"
            aria-label={isMobileMenuOpen ? "Close menu" : "Open menu"}
          >
            {isMobileMenuOpen ? (
              <X className="h-5 w-5" />
            ) : (
              <Menu className="h-5 w-5" />
            )}
          </Button>
        </div>
      </div>

      {/* Mobile Menu Overlay */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm md:hidden animate-in fade-in duration-200"
          onClick={() => setIsMobileMenuOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Mobile Navigation - Full Screen Drawer */}
      {isMobileMenuOpen && (
        <nav
          id="mobile-menu"
          className="fixed inset-y-0 right-0 z-50 w-full max-w-sm bg-background border-l border-border shadow-2xl md:hidden animate-in slide-in-from-right duration-300"
          role="navigation"
          aria-label="Mobile navigation"
        >
          <div className="flex flex-col h-full">
            {/* Mobile Menu Header */}
            <div className="flex items-center justify-between p-4 border-b border-border">
              <div className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-linear-to-br from-blue-600 to-cyan-600">
                  <span className="text-sm font-bold text-white">G</span>
                </div>
                <span className="text-lg font-bold bg-linear-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent">
                  GO-PRO
                </span>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setIsMobileMenuOpen(false)}
                className="h-8 w-8"
              >
                <X className="h-5 w-5" />
              </Button>
            </div>

            {/* Mobile Menu Content */}
            <div className="flex-1 overflow-y-auto p-4 space-y-1">
              {navigationItems.map((item) => (
                <Link
                  key={item.title}
                  href={item.href}
                  className={cn(
                    "group flex items-center gap-3 p-4 rounded-xl",
                    "hover:bg-accent transition-all duration-200",
                    "border border-transparent hover:border-border"
                  )}
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors">
                    <item.icon className="h-5 w-5 text-primary" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-semibold text-base">{item.title}</span>
                      {item.badge && (
                        <Badge variant="secondary" className="text-xs">
                          New
                        </Badge>
                      )}
                    </div>
                    <p className="text-sm text-muted-foreground line-clamp-1">
                      {item.description}
                    </p>
                  </div>
                  <ChevronRight className="h-5 w-5 text-muted-foreground group-hover:text-foreground group-hover:translate-x-1 transition-all" />
                </Link>
              ))}
            </div>

            {/* Mobile Menu Footer */}
            <div className="p-4 border-t border-border space-y-3">
              {mounted && !loading && (
                user ? (
                  <>
                    {/* User Info */}
                    <div className="flex items-center gap-3 p-3 rounded-lg bg-accent/50">
                      <Avatar className="h-12 w-12">
                        <AvatarImage src={user.photoURL || undefined} alt={user.displayName || 'User'} />
                        <AvatarFallback className="bg-linear-to-br from-primary to-primary/80">
                          {user.displayName?.charAt(0)?.toUpperCase() || user.email?.charAt(0)?.toUpperCase() || 'U'}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-semibold truncate">
                          {user.displayName || 'User'}
                        </p>
                        <p className="text-xs text-muted-foreground truncate">
                          {user.email}
                        </p>
                        {backendUser?.role && (
                          <Badge variant="secondary" className="mt-1 text-xs">
                            {backendUser.role}
                          </Badge>
                        )}
                      </div>
                    </div>

                    {/* User Actions */}
                    <div className="space-y-2">
                      <Button
                        variant="outline"
                        size="lg"
                        className="w-full justify-start"
                        asChild
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        <Link href="/profile">
                          <UserCircle className="mr-2 h-5 w-5" />
                          View Profile
                        </Link>
                      </Button>
                      <Button
                        variant="outline"
                        size="lg"
                        className="w-full justify-start"
                        asChild
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        <Link href="/settings">
                          <Settings className="mr-2 h-5 w-5" />
                          Settings
                        </Link>
                      </Button>
                      <Button
                        variant="destructive"
                        size="lg"
                        className="w-full justify-start"
                        onClick={handleSignOut}
                      >
                        <LogOut className="mr-2 h-5 w-5" />
                        Sign out
                      </Button>
                    </div>
                  </>
                ) : (
                  <>
                    <Button
                      size="lg"
                      variant="outline"
                      className="w-full h-12 font-semibold"
                      asChild
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      <Link href="/signin">
                        <LogIn className="mr-2 h-5 w-5" />
                        Sign in
                      </Link>
                    </Button>
                    <Button
                      size="lg"
                      className="w-full bg-linear-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-white font-semibold shadow-lg h-12"
                      asChild
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      <Link href="/signup">Get Started Free</Link>
                    </Button>
                  </>
                )
              )}
              <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground">
                <Badge variant="secondary" className="text-xs">
                  Beta Version
                </Badge>
                <span>•</span>
                <span>Free to use</span>
              </div>
            </div>
          </div>
        </nav>
      )}
    </header>
  );
};

export default Header;
