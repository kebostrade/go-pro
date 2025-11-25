// Responsive design utilities and constants for GO-PRO platform

// Breakpoint constants (matching Tailwind CSS defaults)
export const BREAKPOINTS = {
  sm: 640,   // Small devices (landscape phones, 640px and up)
  md: 768,   // Medium devices (tablets, 768px and up)
  lg: 1024,  // Large devices (desktops, 1024px and up)
  xl: 1280,  // Extra large devices (large desktops, 1280px and up)
  '2xl': 1536 // 2X large devices (larger desktops, 1536px and up)
} as const;

// Common responsive class patterns
export const RESPONSIVE_CLASSES = {
  // Container classes
  container: "container max-w-screen-2xl px-4 py-8",
  containerSmall: "container max-w-4xl px-4 py-6",
  containerLarge: "container max-w-screen-2xl px-4 sm:px-6 lg:px-8 py-8",

  // Grid layouts
  gridCols1: "grid grid-cols-1",
  gridCols2: "grid grid-cols-1 md:grid-cols-2",
  gridCols3: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3",
  gridCols4: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4",
  gridAutoFit: "grid grid-cols-[repeat(auto-fit,minmax(300px,1fr))]",

  // Flex layouts
  flexCol: "flex flex-col",
  flexRow: "flex flex-col md:flex-row",
  flexBetween: "flex items-center justify-between",
  flexCenter: "flex items-center justify-center",
  flexStart: "flex items-center justify-start",
  flexEnd: "flex items-center justify-end",

  // Spacing
  gap2: "gap-2",
  gap4: "gap-4",
  gap6: "gap-6",
  gap8: "gap-8",
  gapResponsive: "gap-4 md:gap-6 lg:gap-8",

  // Text sizes
  textSm: "text-sm",
  textBase: "text-sm md:text-base",
  textLg: "text-base md:text-lg",
  textXl: "text-lg md:text-xl",
  text2xl: "text-xl md:text-2xl",
  text3xl: "text-2xl md:text-3xl",

  // Padding
  p4: "p-4",
  p6: "p-6",
  pResponsive: "p-4 md:p-6",
  px4: "px-4",
  px6: "px-6",
  pxResponsive: "px-4 md:px-6",
  py4: "py-4",
  py6: "py-6",
  pyResponsive: "py-4 md:py-6",

  // Margins
  mb4: "mb-4",
  mb6: "mb-6",
  mb8: "mb-8",
  mbResponsive: "mb-4 md:mb-6 lg:mb-8",
  mt4: "mt-4",
  mt6: "mt-6",
  mt8: "mt-8",
  mtResponsive: "mt-4 md:mt-6 lg:mt-8",

  // Width and height
  wFull: "w-full",
  hFull: "h-full",
  wScreen: "w-screen",
  hScreen: "h-screen",
  maxWScreen: "max-w-screen-2xl",

  // Visibility
  hiddenSm: "hidden sm:block",
  hiddenMd: "hidden md:block",
  hiddenLg: "hidden lg:block",
  showSm: "block sm:hidden",
  showMd: "block md:hidden",
  showLg: "block lg:hidden",
} as const;

// Component-specific responsive patterns
export const COMPONENT_PATTERNS = {
  // Card layouts
  cardGrid: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6",
  cardGridLarge: "grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6",
  cardList: "space-y-4 md:space-y-6",

  // Navigation
  navMobile: "md:hidden",
  navDesktop: "hidden md:flex",
  navTablet: "hidden md:flex lg:hidden",
  navFull: "hidden lg:flex",

  // Sidebar layouts
  sidebarLayout: "grid grid-cols-1 lg:grid-cols-4 gap-6",
  sidebarMain: "lg:col-span-3",
  sidebarAside: "lg:col-span-1",

  // Two column layouts
  twoColLayout: "grid grid-cols-1 lg:grid-cols-2 gap-6",
  twoColLayoutReverse: "grid grid-cols-1 lg:grid-cols-2 gap-6 lg:grid-flow-col-dense",

  // Three column layouts
  threeColLayout: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6",
  threeColMain: "md:col-span-2 lg:col-span-2",
  threeColSide: "md:col-span-1 lg:col-span-1",

  // Header patterns
  headerMobile: "flex h-14 sm:h-16 items-center justify-between px-4",
  headerDesktop: "flex h-16 lg:h-18 items-center justify-between px-6 lg:px-8",
  headerFull: "flex h-14 sm:h-16 lg:h-18 items-center justify-between px-4 sm:px-6 lg:px-8",

  // Button groups
  buttonGroup: "flex flex-col sm:flex-row gap-2",
  buttonGroupReverse: "flex flex-col-reverse sm:flex-row gap-2",

  // Form layouts
  formGrid: "grid grid-cols-1 md:grid-cols-2 gap-4",
  formStack: "space-y-4 md:space-y-6",

  // Stats grids
  statsGrid2: "grid grid-cols-2 gap-4",
  statsGrid4: "grid grid-cols-2 md:grid-cols-4 gap-4",
  statsGridAuto: "grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4",
} as const;

// Utility functions for responsive behavior
export const useResponsive = () => {
  // Check if we're on the client side
  const isClient = typeof window !== 'undefined';
  
  if (!isClient) {
    return {
      isMobile: false,
      isTablet: false,
      isDesktop: true,
      screenWidth: 1024,
      breakpoint: 'lg' as keyof typeof BREAKPOINTS
    };
  }

  const screenWidth = window.innerWidth;
  
  const isMobile = screenWidth < BREAKPOINTS.md;
  const isTablet = screenWidth >= BREAKPOINTS.md && screenWidth < BREAKPOINTS.lg;
  const isDesktop = screenWidth >= BREAKPOINTS.lg;
  
  let breakpoint: keyof typeof BREAKPOINTS = 'sm';
  if (screenWidth >= BREAKPOINTS['2xl']) breakpoint = '2xl';
  else if (screenWidth >= BREAKPOINTS.xl) breakpoint = 'xl';
  else if (screenWidth >= BREAKPOINTS.lg) breakpoint = 'lg';
  else if (screenWidth >= BREAKPOINTS.md) breakpoint = 'md';
  else if (screenWidth >= BREAKPOINTS.sm) breakpoint = 'sm';

  return {
    isMobile,
    isTablet,
    isDesktop,
    screenWidth,
    breakpoint
  };
};

// Helper function to combine responsive classes
export const cn = (...classes: (string | undefined | null | false)[]): string => {
  return classes.filter(Boolean).join(' ');
};

// Responsive image sizes for different breakpoints
export const IMAGE_SIZES = {
  avatar: {
    sm: "w-8 h-8",
    md: "w-10 h-10",
    lg: "w-12 h-12",
    xl: "w-16 h-16"
  },
  icon: {
    sm: "w-4 h-4",
    md: "w-5 h-5",
    lg: "w-6 h-6",
    xl: "w-8 h-8"
  },
  logo: {
    sm: "h-6 w-6",
    md: "h-8 w-8",
    lg: "h-10 w-10",
    xl: "h-12 w-12"
  }
} as const;

// Responsive typography scale
export const TYPOGRAPHY = {
  heading: {
    h1: "text-2xl md:text-3xl lg:text-4xl font-bold",
    h2: "text-xl md:text-2xl lg:text-3xl font-bold",
    h3: "text-lg md:text-xl lg:text-2xl font-semibold",
    h4: "text-base md:text-lg lg:text-xl font-semibold",
    h5: "text-sm md:text-base lg:text-lg font-medium",
    h6: "text-xs md:text-sm lg:text-base font-medium"
  },
  body: {
    large: "text-base md:text-lg",
    normal: "text-sm md:text-base",
    small: "text-xs md:text-sm",
    tiny: "text-xs"
  },
  display: {
    large: "text-3xl md:text-4xl lg:text-5xl xl:text-6xl font-bold",
    medium: "text-2xl md:text-3xl lg:text-4xl xl:text-5xl font-bold",
    small: "text-xl md:text-2xl lg:text-3xl xl:text-4xl font-bold"
  }
} as const;

// Animation and transition classes
export const ANIMATIONS = {
  fadeIn: "animate-in fade-in duration-200",
  slideIn: "animate-in slide-in-from-bottom-4 duration-300",
  scaleIn: "animate-in zoom-in-95 duration-200",
  hover: "transition-all duration-200 hover:scale-105",
  hoverSoft: "transition-colors duration-200",
  loading: "animate-pulse"
} as const;

// Accessibility helpers
export const A11Y = {
  srOnly: "sr-only",
  focusVisible: "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary",
  skipLink: "absolute left-[-10000px] top-auto w-1 h-1 overflow-hidden focus:left-6 focus:top-7 focus:w-auto focus:h-auto focus:overflow-visible",
  highContrast: "contrast-more:border-black contrast-more:text-black"
} as const;

// Common responsive patterns for specific components
export const getResponsiveCardGrid = (itemsCount: number): string => {
  if (itemsCount <= 2) return COMPONENT_PATTERNS.twoColLayout;
  if (itemsCount <= 4) return COMPONENT_PATTERNS.cardGrid;
  return COMPONENT_PATTERNS.cardGridLarge;
};

export const getResponsiveTextSize = (level: 'small' | 'medium' | 'large' = 'medium'): string => {
  switch (level) {
    case 'small': return TYPOGRAPHY.body.small;
    case 'large': return TYPOGRAPHY.body.large;
    default: return TYPOGRAPHY.body.normal;
  }
};

export const getResponsiveSpacing = (size: 'small' | 'medium' | 'large' = 'medium'): string => {
  switch (size) {
    case 'small': return RESPONSIVE_CLASSES.gap4;
    case 'large': return RESPONSIVE_CLASSES.gap8;
    default: return RESPONSIVE_CLASSES.gap6;
  }
};

// Export commonly used combinations
export const COMMON_LAYOUTS = {
  pageContainer: cn(RESPONSIVE_CLASSES.container),
  cardContainer: cn(COMPONENT_PATTERNS.cardGrid, RESPONSIVE_CLASSES.gapResponsive),
  flexContainer: cn(RESPONSIVE_CLASSES.flexBetween, RESPONSIVE_CLASSES.mbResponsive),
  sectionHeader: cn(TYPOGRAPHY.heading.h2, RESPONSIVE_CLASSES.mbResponsive),
  buttonGroup: cn(COMPONENT_PATTERNS.buttonGroup),
  statsGrid: cn(COMPONENT_PATTERNS.statsGrid4, RESPONSIVE_CLASSES.gap4)
} as const;
