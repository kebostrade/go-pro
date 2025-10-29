"use client";

import { useState, useEffect } from "react";

// Hook to detect screen size
export const useScreenSize = () => {
  const [screenSize, setScreenSize] = useState<'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'>('md');

  useEffect(() => {
    const checkScreenSize = () => {
      const width = window.innerWidth;
      if (width < 480) setScreenSize('xs');
      else if (width < 640) setScreenSize('sm');
      else if (width < 768) setScreenSize('md');
      else if (width < 1024) setScreenSize('lg');
      else if (width < 1280) setScreenSize('xl');
      else setScreenSize('2xl');
    };

    checkScreenSize();
    window.addEventListener('resize', checkScreenSize);
    return () => window.removeEventListener('resize', checkScreenSize);
  }, []);

  return screenSize;
};

// Responsive header configuration
export const getHeaderConfig = (screenSize: string) => {
  const configs = {
    xs: {
      height: 'h-12',
      logoSize: 'h-6 w-6',
      logoTextSize: 'text-sm',
      iconSize: 'h-3 w-3',
      buttonSize: 'h-7 w-7',
      padding: 'px-2',
      spacing: 'space-x-1',
      showSubtitle: false,
      showGitHub: false,
      showBeta: false,
      navType: 'mobile'
    },
    sm: {
      height: 'h-14',
      logoSize: 'h-7 w-7',
      logoTextSize: 'text-base',
      iconSize: 'h-3 w-3',
      buttonSize: 'h-8 w-8',
      padding: 'px-3',
      spacing: 'space-x-1',
      showSubtitle: true,
      showGitHub: false,
      showBeta: true,
      navType: 'mobile'
    },
    md: {
      height: 'h-16',
      logoSize: 'h-8 w-8',
      logoTextSize: 'text-lg',
      iconSize: 'h-4 w-4',
      buttonSize: 'h-9 w-9',
      padding: 'px-4',
      spacing: 'space-x-2',
      showSubtitle: true,
      showGitHub: true,
      showBeta: true,
      navType: 'simplified'
    },
    lg: {
      height: 'h-16',
      logoSize: 'h-8 w-8',
      logoTextSize: 'text-lg',
      iconSize: 'h-4 w-4',
      buttonSize: 'h-9 w-9',
      padding: 'px-6',
      spacing: 'space-x-2',
      showSubtitle: true,
      showGitHub: true,
      showBeta: true,
      navType: 'full'
    },
    xl: {
      height: 'h-18',
      logoSize: 'h-9 w-9',
      logoTextSize: 'text-xl',
      iconSize: 'h-4 w-4',
      buttonSize: 'h-10 w-10',
      padding: 'px-8',
      spacing: 'space-x-3',
      showSubtitle: true,
      showGitHub: true,
      showBeta: true,
      navType: 'full'
    },
    '2xl': {
      height: 'h-20',
      logoSize: 'h-10 w-10',
      logoTextSize: 'text-2xl',
      iconSize: 'h-5 w-5',
      buttonSize: 'h-11 w-11',
      padding: 'px-12',
      spacing: 'space-x-4',
      showSubtitle: true,
      showGitHub: true,
      showBeta: true,
      navType: 'full'
    }
  };

  return configs[screenSize as keyof typeof configs] || configs.md;
};

// Mobile menu animation hook
export const useMobileMenu = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [isAnimating, setIsAnimating] = useState(false);

  const toggle = () => {
    if (isOpen) {
      setIsAnimating(true);
      setTimeout(() => {
        setIsOpen(false);
        setIsAnimating(false);
      }, 200);
    } else {
      setIsOpen(true);
    }
  };

  const close = () => {
    setIsAnimating(true);
    setTimeout(() => {
      setIsOpen(false);
      setIsAnimating(false);
    }, 200);
  };

  return { isOpen, isAnimating, toggle, close };
};

// Touch-friendly button component
interface TouchButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
  variant?: 'ghost' | 'outline' | 'default';
}

export const TouchButton = ({ 
  children, 
  onClick, 
  className = '', 
  size = 'md',
  variant = 'ghost'
}: TouchButtonProps) => {
  const sizeClasses = {
    sm: 'min-h-[40px] min-w-[40px] p-2',
    md: 'min-h-[44px] min-w-[44px] p-2',
    lg: 'min-h-[48px] min-w-[48px] p-3'
  };

  const variantClasses = {
    ghost: 'hover:bg-accent hover:text-accent-foreground',
    outline: 'border border-input hover:bg-accent hover:text-accent-foreground',
    default: 'bg-primary text-primary-foreground hover:bg-primary/90'
  };

  return (
    <button
      onClick={onClick}
      className={`
        inline-flex items-center justify-center rounded-md text-sm font-medium 
        transition-colors focus-visible:outline-none focus-visible:ring-2 
        focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 
        disabled:pointer-events-none ring-offset-background
        ${sizeClasses[size]}
        ${variantClasses[variant]}
        ${className}
      `}
    >
      {children}
    </button>
  );
};

// Responsive navigation item
interface ResponsiveNavItemProps {
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  title: string;
  description?: string;
  onClick?: () => void;
  isMobile?: boolean;
}

export const ResponsiveNavItem = ({ 
  href, 
  icon: Icon, 
  title, 
  description, 
  onClick,
  isMobile = false 
}: ResponsiveNavItemProps) => {
  if (isMobile) {
    return (
      <a
        href={href}
        onClick={onClick}
        className="flex items-center space-x-3 rounded-lg p-3 min-h-[48px] text-sm hover:bg-accent hover:text-accent-foreground transition-colors active:bg-accent/80"
      >
        <Icon className="h-5 w-5 flex-shrink-0" />
        <div className="flex-1 min-w-0">
          <div className="font-medium truncate">{title}</div>
          {description && (
            <div className="text-xs text-muted-foreground line-clamp-1">{description}</div>
          )}
        </div>
      </a>
    );
  }

  return (
    <a
      href={href}
      onClick={onClick}
      className="flex items-center space-x-2 px-3 py-2 text-sm font-medium rounded-md hover:bg-accent hover:text-accent-foreground transition-colors header-nav-item"
    >
      <Icon className="h-4 w-4" />
      <span>{title}</span>
    </a>
  );
};
