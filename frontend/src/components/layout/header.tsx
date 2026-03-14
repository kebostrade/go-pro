"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
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
  Brain,
  Sparkles,
  Bot,
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
import { LucideIcon } from "lucide-react";

interface NavItem {
  title: string;
  href: string;
  description: string;
  icon: LucideIcon;
}

const Header = () => {
  const { user, signOut, loading } = useAuth();
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

  const navigationItems: NavItem[] = [
    { title: "Learn", href: "/learn", description: "Interactive Go lessons and tutorials", icon: BookOpen },
    { title: "Practice", href: "/practice", description: "Coding exercises and challenges", icon: Code2 },
    { title: "Prompt Engineering", href: "/learn/prompt-engineering", description: "Master AI prompts and LLM techniques", icon: Sparkles },
    { title: "OpenClaw", href: "/learn/openclaw", description: "Build self-hosted AI agents", icon: Bot },
    { title: "Interviews", href: "/interviews", description: "Master coding interviews", icon: Brain },
    { title: "Tutorials", href: "/tutorials", description: "Comprehensive tutorials", icon: GraduationCap },
    { title: "Projects", href: "/projects", description: "Real-world Go applications", icon: Trophy },
    { title: "Community", href: "/community", description: "Connect with Go developers", icon: Users },
  ];

  return (
    <header
      className={cn(
        "sticky top-0 z-50 w-full border-b transition-all duration-300",
        scrolled
          ? "border-border/50 bg-background/95 backdrop-blur-xl shadow-sm"
          : "border-border/30 bg-background/80 backdrop-blur-lg"
      )}
      role="banner"
    >
      <div className="mx-auto flex h-14 sm:h-16 items-center justify-between px-3 sm:px-4 md:px-6 lg:px-8 xl:px-12 2xl:px-16 max-w-[1800px] w-full">
        <Link href="/" className="flex items-center gap-2 sm:gap-3 group shrink-0">
          <div className="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 shadow-md shadow-blue-500/20 group-hover:shadow-blue-500/40 transition-all duration-300 group-hover:scale-105">
            <span className="text-base sm:text-lg font-bold text-white">G</span>
          </div>
          <div className="flex flex-col">
            <span className="text-base sm:text-lg font-bold bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent whitespace-nowrap">
              GO-PRO
            </span>
            <span className="hidden sm:block text-[10px] sm:text-xs text-muted-foreground -mt-0.5 whitespace-nowrap">
              Learn Go Programming
            </span>
          </div>
        </Link>

        <nav className="hidden lg:flex items-center gap-0.5" aria-label="Main navigation">
          {navigationItems.map((item) => (
            <Link
              key={item.title}
              href={item.href}
              className={cn(
                "flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm font-medium",
                "text-muted-foreground hover:text-foreground hover:bg-accent/50",
                "transition-all duration-200"
              )}
            >
              <item.icon className="h-4 w-4" />
              <span>{item.title}</span>
            </Link>
          ))}
        </nav>

        <div className="flex items-center gap-1 sm:gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleTheme}
            className="h-8 w-8 sm:h-9 sm:w-9 rounded-lg"
            aria-label="Toggle theme"
            suppressHydrationWarning
          >
            {mounted && isDark ? (
              <Sun className="h-4 w-4 sm:h-5 sm:w-5 text-yellow-500" />
            ) : (
              <Moon className="h-4 w-4 sm:h-5 sm:w-5 text-blue-600" />
            )}
          </Button>

          {mounted && !loading && (
            <div className="flex items-center">
              {user ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="relative h-8 w-8 sm:h-9 sm:w-9 rounded-full">
                      <Avatar className="h-7 w-7 sm:h-8 sm:w-8">
                        <AvatarImage src={user.photoURL || undefined} alt={user.displayName || 'User'} />
                        <AvatarFallback className="bg-gradient-to-br from-blue-600 to-cyan-500 text-white text-xs sm:text-sm">
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
                <Button variant="ghost" size="sm" className="h-8 sm:h-9 text-xs sm:text-sm hidden sm:flex" asChild>
                  <Link href="/signin">
                    <LogIn className="mr-1.5 sm:mr-2 h-3.5 w-3.5 sm:h-4 sm:w-4" />
                    Sign in
                  </Link>
                </Button>
              )}
            </div>
          )}

          <Button
            variant="ghost"
            size="icon"
            className="lg:hidden h-8 w-8 sm:h-9 sm:w-9 rounded-lg"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-expanded={isMobileMenuOpen}
            aria-controls="mobile-menu"
            aria-label={isMobileMenuOpen ? "Close menu" : "Open menu"}
          >
            {isMobileMenuOpen ? <X className="h-4 w-4 sm:h-5 sm:w-5" /> : <Menu className="h-4 w-4 sm:h-5 sm:w-5" />}
          </Button>
        </div>
      </div>

      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm lg:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
          aria-hidden="true"
        />
      )}

      {isMobileMenuOpen && (
        <nav
          id="mobile-menu"
          className="fixed inset-y-0 right-0 z-50 w-full max-w-xs bg-background border-l shadow-2xl lg:hidden"
          role="navigation"
          aria-label="Mobile navigation"
        >
          <div className="flex flex-col h-full">
            <div className="flex items-center justify-between p-3 sm:p-4 border-b">
              <div className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500">
                  <span className="text-sm font-bold text-white">G</span>
                </div>
                <span className="text-lg font-bold bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent">
                  GO-PRO
                </span>
              </div>
              <Button variant="ghost" size="icon" onClick={() => setIsMobileMenuOpen(false)} className="h-8 w-8">
                <X className="h-5 w-5" />
              </Button>
            </div>

            <div className="flex-1 overflow-y-auto p-3 space-y-1">
              {navigationItems.map((item) => (
                <Link
                  key={item.title}
                  href={item.href}
                  className="group flex items-center gap-3 p-3 rounded-xl hover:bg-accent transition-colors"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <div className="flex h-9 w-9 sm:h-10 sm:w-10 items-center justify-center rounded-lg bg-blue-500/10 group-hover:bg-blue-500/20 transition-colors">
                    <item.icon className="h-4 w-4 sm:h-5 sm:w-5 text-blue-600" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <span className="text-sm sm:text-base font-medium">{item.title}</span>
                    <p className="text-xs text-muted-foreground line-clamp-1">{item.description}</p>
                  </div>
                  <ChevronRight className="h-4 w-4 text-muted-foreground group-hover:translate-x-0.5 transition-transform" />
                </Link>
              ))}
            </div>

            <div className="p-3 sm:p-4 border-t space-y-2 sm:space-y-3">
              {mounted && !loading && (
                user ? (
                  <>
                    <div className="flex items-center gap-3 p-2.5 sm:p-3 rounded-lg bg-muted/50">
                      <Avatar className="h-9 w-9 sm:h-10 sm:w-10">
                        <AvatarImage src={user.photoURL || undefined} alt={user.displayName || 'User'} />
                        <AvatarFallback className="bg-gradient-to-br from-blue-600 to-cyan-500 text-white text-sm">
                          {user.displayName?.charAt(0)?.toUpperCase() || user.email?.charAt(0)?.toUpperCase() || 'U'}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate">{user.displayName || 'User'}</p>
                        <p className="text-xs text-muted-foreground truncate">{user.email}</p>
                      </div>
                    </div>
                    <Button variant="outline" size="sm" className="w-full justify-start h-9 sm:h-10" asChild onClick={() => setIsMobileMenuOpen(false)}>
                      <Link href="/profile"><UserCircle className="mr-2 h-4 w-4" />Profile</Link>
                    </Button>
                    <Button variant="outline" size="sm" className="w-full justify-start h-9 sm:h-10" asChild onClick={() => setIsMobileMenuOpen(false)}>
                      <Link href="/settings"><Settings className="mr-2 h-4 w-4" />Settings</Link>
                    </Button>
                    <Button variant="destructive" size="sm" className="w-full justify-start h-9 sm:h-10" onClick={handleSignOut}>
                      <LogOut className="mr-2 h-4 w-4" />Sign out
                    </Button>
                  </>
                ) : (
                  <>
                    <Button variant="outline" size="sm" className="w-full h-9 sm:h-11" asChild onClick={() => setIsMobileMenuOpen(false)}>
                      <Link href="/signin"><LogIn className="mr-2 h-4 w-4" />Sign in</Link>
                    </Button>
                    <Button size="sm" className="w-full h-9 sm:h-11 bg-gradient-to-r from-blue-600 to-cyan-500 hover:from-blue-700 hover:to-cyan-600" asChild onClick={() => setIsMobileMenuOpen(false)}>
                      <Link href="/signup">Get Started Free</Link>
                    </Button>
                  </>
                )
              )}
            </div>
          </div>
        </nav>
      )}
    </header>
  );
};

export default Header;
