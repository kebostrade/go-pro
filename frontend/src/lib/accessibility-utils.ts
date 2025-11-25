// Accessibility utilities for GO-PRO platform

// ARIA attributes and roles
export const ARIA = {
  // Common ARIA attributes
  label: (label: string) => ({ 'aria-label': label }),
  labelledBy: (id: string) => ({ 'aria-labelledby': id }),
  describedBy: (id: string) => ({ 'aria-describedby': id }),
  expanded: (expanded: boolean) => ({ 'aria-expanded': expanded }),
  selected: (selected: boolean) => ({ 'aria-selected': selected }),
  checked: (checked: boolean) => ({ 'aria-checked': checked }),
  disabled: (disabled: boolean) => ({ 'aria-disabled': disabled }),
  hidden: (hidden: boolean) => ({ 'aria-hidden': hidden }),
  live: (politeness: 'polite' | 'assertive' | 'off') => ({ 'aria-live': politeness }),
  atomic: (atomic: boolean) => ({ 'aria-atomic': atomic }),
  busy: (busy: boolean) => ({ 'aria-busy': busy }),
  current: (current: 'page' | 'step' | 'location' | 'date' | 'time' | 'true' | 'false') => ({ 'aria-current': current }),
  level: (level: number) => ({ 'aria-level': level }),
  setSize: (size: number) => ({ 'aria-setsize': size }),
  posInSet: (position: number) => ({ 'aria-posinset': position }),
  
  // Form-specific ARIA attributes
  required: (required: boolean) => ({ 'aria-required': required }),
  invalid: (invalid: boolean) => ({ 'aria-invalid': invalid }),
  errorMessage: (id: string) => ({ 'aria-errormessage': id }),
  
  // Navigation-specific ARIA attributes
  hasPopup: (popup: 'menu' | 'listbox' | 'tree' | 'grid' | 'dialog' | 'true' | 'false') => ({ 'aria-haspopup': popup }),
  controls: (id: string) => ({ 'aria-controls': id }),
  owns: (id: string) => ({ 'aria-owns': id }),
  
  // Progress and status
  valueNow: (value: number) => ({ 'aria-valuenow': value }),
  valueMin: (min: number) => ({ 'aria-valuemin': min }),
  valueMax: (max: number) => ({ 'aria-valuemax': max }),
  valueText: (text: string) => ({ 'aria-valuetext': text }),
} as const;

// Common ARIA roles
export const ROLES = {
  // Landmark roles
  banner: 'banner',
  navigation: 'navigation',
  main: 'main',
  complementary: 'complementary',
  contentinfo: 'contentinfo',
  search: 'search',
  form: 'form',
  region: 'region',
  
  // Widget roles
  button: 'button',
  link: 'link',
  menuitem: 'menuitem',
  tab: 'tab',
  tabpanel: 'tabpanel',
  option: 'option',
  checkbox: 'checkbox',
  radio: 'radio',
  slider: 'slider',
  spinbutton: 'spinbutton',
  textbox: 'textbox',
  combobox: 'combobox',
  listbox: 'listbox',
  tree: 'tree',
  grid: 'grid',
  dialog: 'dialog',
  alertdialog: 'alertdialog',
  tooltip: 'tooltip',
  status: 'status',
  alert: 'alert',
  log: 'log',
  marquee: 'marquee',
  timer: 'timer',
  
  // Structure roles
  article: 'article',
  document: 'document',
  application: 'application',
  group: 'group',
  heading: 'heading',
  img: 'img',
  list: 'list',
  listitem: 'listitem',
  table: 'table',
  row: 'row',
  cell: 'cell',
  columnheader: 'columnheader',
  rowheader: 'rowheader',
  separator: 'separator',
  toolbar: 'toolbar',
  menu: 'menu',
  menubar: 'menubar',
  tablist: 'tablist',
  presentation: 'presentation',
  none: 'none',
} as const;

// Focus management utilities
export const focusUtils = {
  // Trap focus within an element
  trapFocus: (element: HTMLElement) => {
    const focusableElements = element.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );
    const firstElement = focusableElements[0] as HTMLElement;
    const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement;

    const handleTabKey = (e: KeyboardEvent) => {
      if (e.key === 'Tab') {
        if (e.shiftKey) {
          if (document.activeElement === firstElement) {
            lastElement.focus();
            e.preventDefault();
          }
        } else {
          if (document.activeElement === lastElement) {
            firstElement.focus();
            e.preventDefault();
          }
        }
      }
    };

    element.addEventListener('keydown', handleTabKey);
    firstElement?.focus();

    return () => {
      element.removeEventListener('keydown', handleTabKey);
    };
  },

  // Focus first focusable element
  focusFirst: (element: HTMLElement) => {
    const focusable = element.querySelector(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    ) as HTMLElement;
    focusable?.focus();
  },

  // Check if element is focusable
  isFocusable: (element: HTMLElement): boolean => {
    const focusableSelectors = [
      'button:not([disabled])',
      '[href]',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      '[tabindex]:not([tabindex="-1"])'
    ];
    
    return focusableSelectors.some(selector => element.matches(selector));
  },

  // Get all focusable elements within a container
  getFocusableElements: (container: HTMLElement): HTMLElement[] => {
    const focusableSelectors = [
      'button:not([disabled])',
      '[href]',
      'input:not([disabled])',
      'select:not([disabled])',
      'textarea:not([disabled])',
      '[tabindex]:not([tabindex="-1"])'
    ].join(', ');
    
    return Array.from(container.querySelectorAll(focusableSelectors));
  }
};

// Keyboard navigation utilities
export const keyboardUtils = {
  // Common keyboard event handlers
  onEnterOrSpace: (callback: () => void) => (e: KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      callback();
    }
  },

  onEscape: (callback: () => void) => (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      e.preventDefault();
      callback();
    }
  },

  onArrowKeys: (callbacks: {
    up?: () => void;
    down?: () => void;
    left?: () => void;
    right?: () => void;
  }) => (e: KeyboardEvent) => {
    switch (e.key) {
      case 'ArrowUp':
        e.preventDefault();
        callbacks.up?.();
        break;
      case 'ArrowDown':
        e.preventDefault();
        callbacks.down?.();
        break;
      case 'ArrowLeft':
        e.preventDefault();
        callbacks.left?.();
        break;
      case 'ArrowRight':
        e.preventDefault();
        callbacks.right?.();
        break;
    }
  },

  // Navigation for lists and menus
  handleListNavigation: (
    items: HTMLElement[],
    currentIndex: number,
    setCurrentIndex: (index: number) => void,
    onSelect?: (index: number) => void
  ) => (e: KeyboardEvent) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        const nextIndex = currentIndex < items.length - 1 ? currentIndex + 1 : 0;
        setCurrentIndex(nextIndex);
        items[nextIndex]?.focus();
        break;
      case 'ArrowUp':
        e.preventDefault();
        const prevIndex = currentIndex > 0 ? currentIndex - 1 : items.length - 1;
        setCurrentIndex(prevIndex);
        items[prevIndex]?.focus();
        break;
      case 'Home':
        e.preventDefault();
        setCurrentIndex(0);
        items[0]?.focus();
        break;
      case 'End':
        e.preventDefault();
        const lastIndex = items.length - 1;
        setCurrentIndex(lastIndex);
        items[lastIndex]?.focus();
        break;
      case 'Enter':
      case ' ':
        e.preventDefault();
        onSelect?.(currentIndex);
        break;
    }
  }
};

// Screen reader utilities
export const screenReaderUtils = {
  // Announce message to screen readers
  announce: (message: string, priority: 'polite' | 'assertive' = 'polite') => {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', priority);
    announcement.setAttribute('aria-atomic', 'true');
    announcement.className = 'sr-only';
    announcement.textContent = message;
    
    document.body.appendChild(announcement);
    
    // Remove after announcement
    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  },

  // Create visually hidden text for screen readers
  createSROnlyText: (text: string): HTMLSpanElement => {
    const span = document.createElement('span');
    span.className = 'sr-only';
    span.textContent = text;
    return span;
  },

  // Update screen reader text
  updateSRText: (element: HTMLElement, text: string) => {
    let srElement = element.querySelector('.sr-only');
    if (!srElement) {
      srElement = screenReaderUtils.createSROnlyText(text);
      element.appendChild(srElement);
    } else {
      srElement.textContent = text;
    }
  }
};

// Color contrast and visual accessibility
export const visualUtils = {
  // Check if user prefers reduced motion
  prefersReducedMotion: (): boolean => {
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  },

  // Check if user prefers high contrast
  prefersHighContrast: (): boolean => {
    return window.matchMedia('(prefers-contrast: high)').matches;
  },

  // Check if user prefers dark mode
  prefersDarkMode: (): boolean => {
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  },

  // Apply reduced motion styles conditionally
  withReducedMotion: (normalClass: string, reducedClass: string = ''): string => {
    if (typeof window !== 'undefined' && visualUtils.prefersReducedMotion()) {
      return reducedClass;
    }
    return normalClass;
  }
};

// Form accessibility helpers
export const formUtils = {
  // Generate accessible form field props
  getFieldProps: (
    id: string,
    label: string,
    error?: string,
    description?: string,
    required: boolean = false
  ) => {
    const props: Record<string, any> = {
      id,
      'aria-label': label,
      'aria-required': required,
    };

    if (error) {
      props['aria-invalid'] = true;
      props['aria-errormessage'] = `${id}-error`;
    }

    if (description) {
      props['aria-describedby'] = `${id}-description`;
    }

    return props;
  },

  // Generate error message props
  getErrorProps: (fieldId: string) => ({
    id: `${fieldId}-error`,
    role: 'alert',
    'aria-live': 'polite'
  }),

  // Generate description props
  getDescriptionProps: (fieldId: string) => ({
    id: `${fieldId}-description`
  })
};

// Common accessibility patterns
export const a11yPatterns = {
  // Skip link
  skipLink: {
    className: "absolute left-[-10000px] top-auto w-1 h-1 overflow-hidden focus:left-6 focus:top-7 focus:w-auto focus:h-auto focus:overflow-visible focus:z-50 focus:bg-primary focus:text-primary-foreground focus:px-4 focus:py-2 focus:rounded",
    href: "#main-content",
    children: "Skip to main content"
  },

  // Loading state
  loading: {
    role: ROLES.status,
    'aria-live': 'polite',
    'aria-label': 'Loading content'
  },

  // Error state
  error: {
    role: ROLES.alert,
    'aria-live': 'assertive'
  },

  // Success state
  success: {
    role: ROLES.status,
    'aria-live': 'polite'
  },

  // Modal dialog
  modal: {
    role: ROLES.dialog,
    'aria-modal': true,
    'aria-labelledby': 'modal-title',
    'aria-describedby': 'modal-description'
  },

  // Tooltip
  tooltip: {
    role: ROLES.tooltip,
    'aria-hidden': true
  },

  // Tab panel
  tabPanel: (id: string, tabId: string) => ({
    role: ROLES.tabpanel,
    id,
    'aria-labelledby': tabId,
    tabIndex: 0
  }),

  // Tab
  tab: (id: string, panelId: string, selected: boolean) => ({
    role: ROLES.tab,
    id,
    'aria-controls': panelId,
    'aria-selected': selected,
    tabIndex: selected ? 0 : -1
  })
};

// Utility to combine accessibility props
export const combineA11yProps = (...propObjects: Record<string, any>[]): Record<string, any> => {
  return propObjects.reduce((combined, props) => ({ ...combined, ...props }), {});
};

// Export commonly used combinations
export const COMMON_A11Y = {
  button: (label: string) => ({
    type: 'button',
    'aria-label': label,
    className: 'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2'
  }),
  
  link: (label: string) => ({
    'aria-label': label,
    className: 'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded'
  }),
  
  heading: (level: number) => ({
    role: ROLES.heading,
    'aria-level': level
  }),
  
  list: () => ({
    role: ROLES.list
  }),
  
  listItem: () => ({
    role: ROLES.listitem
  })
} as const;
