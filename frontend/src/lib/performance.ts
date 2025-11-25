// Performance optimization utilities for the lesson pages

export interface PerformanceMetrics {
  loadTime: number;
  renderTime: number;
  interactionTime: number;
  memoryUsage?: number;
  bundleSize?: number;
}

export class PerformanceMonitor {
  private startTime: number;
  private metrics: PerformanceMetrics;

  constructor() {
    this.startTime = performance.now();
    this.metrics = {
      loadTime: 0,
      renderTime: 0,
      interactionTime: 0
    };
  }

  // Measure page load time
  measureLoadTime(): number {
    const loadTime = performance.now() - this.startTime;
    this.metrics.loadTime = loadTime;
    return loadTime;
  }

  // Measure component render time
  measureRenderTime(componentName: string, renderFn: () => void): number {
    const startRender = performance.now();
    renderFn();
    const renderTime = performance.now() - startRender;
    
    console.log(`${componentName} render time: ${renderTime.toFixed(2)}ms`);
    return renderTime;
  }

  // Measure interaction response time
  measureInteraction(interactionName: string, interactionFn: () => void): number {
    const startInteraction = performance.now();
    interactionFn();
    const interactionTime = performance.now() - startInteraction;
    
    console.log(`${interactionName} interaction time: ${interactionTime.toFixed(2)}ms`);
    this.metrics.interactionTime = Math.max(this.metrics.interactionTime, interactionTime);
    return interactionTime;
  }

  // Get memory usage (if available)
  getMemoryUsage(): number | undefined {
    if ('memory' in performance) {
      const memory = (performance as any).memory;
      return memory.usedJSHeapSize / 1024 / 1024; // MB
    }
    return undefined;
  }

  // Get all metrics
  getMetrics(): PerformanceMetrics {
    return {
      ...this.metrics,
      memoryUsage: this.getMemoryUsage()
    };
  }

  // Log performance summary
  logSummary(): void {
    const metrics = this.getMetrics();
    console.group('Performance Metrics');
    console.log(`Load Time: ${metrics.loadTime.toFixed(2)}ms`);
    console.log(`Render Time: ${metrics.renderTime.toFixed(2)}ms`);
    console.log(`Interaction Time: ${metrics.interactionTime.toFixed(2)}ms`);
    if (metrics.memoryUsage) {
      console.log(`Memory Usage: ${metrics.memoryUsage.toFixed(2)}MB`);
    }
    console.groupEnd();
  }
}

// Debounce utility for performance optimization
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

// Throttle utility for performance optimization
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}

// Lazy loading utility
export function createIntersectionObserver(
  callback: (entries: IntersectionObserverEntry[]) => void,
  options?: IntersectionObserverInit
): IntersectionObserver {
  const defaultOptions: IntersectionObserverInit = {
    root: null,
    rootMargin: '50px',
    threshold: 0.1,
    ...options
  };

  return new IntersectionObserver(callback, defaultOptions);
}

// Image optimization utility
export function optimizeImage(src: string, width?: number, height?: number): string {
  // In a real implementation, this would integrate with an image optimization service
  // For now, we'll just return the original src
  if (width && height) {
    return `${src}?w=${width}&h=${height}&fit=crop&auto=format`;
  }
  return src;
}

// Bundle size analyzer (for development)
export function analyzeBundleSize(): void {
  if (process.env.NODE_ENV === 'development') {
    // This would integrate with webpack-bundle-analyzer or similar
    console.log('Bundle analysis would run here in development mode');
  }
}

// Performance budget checker
export interface PerformanceBudget {
  maxLoadTime: number; // ms
  maxRenderTime: number; // ms
  maxInteractionTime: number; // ms
  maxMemoryUsage: number; // MB
}

export function checkPerformanceBudget(
  metrics: PerformanceMetrics,
  budget: PerformanceBudget
): { passed: boolean; violations: string[] } {
  const violations: string[] = [];

  if (metrics.loadTime > budget.maxLoadTime) {
    violations.push(`Load time exceeded: ${metrics.loadTime.toFixed(2)}ms > ${budget.maxLoadTime}ms`);
  }

  if (metrics.renderTime > budget.maxRenderTime) {
    violations.push(`Render time exceeded: ${metrics.renderTime.toFixed(2)}ms > ${budget.maxRenderTime}ms`);
  }

  if (metrics.interactionTime > budget.maxInteractionTime) {
    violations.push(`Interaction time exceeded: ${metrics.interactionTime.toFixed(2)}ms > ${budget.maxInteractionTime}ms`);
  }

  if (metrics.memoryUsage && metrics.memoryUsage > budget.maxMemoryUsage) {
    violations.push(`Memory usage exceeded: ${metrics.memoryUsage.toFixed(2)}MB > ${budget.maxMemoryUsage}MB`);
  }

  return {
    passed: violations.length === 0,
    violations
  };
}

// Accessibility utilities
export function checkAccessibility(): void {
  // Check for common accessibility issues
  const issues: string[] = [];

  // Check for images without alt text
  const images = document.querySelectorAll('img:not([alt])');
  if (images.length > 0) {
    issues.push(`${images.length} images missing alt text`);
  }

  // Check for buttons without accessible names
  const buttons = document.querySelectorAll('button:not([aria-label]):not([aria-labelledby])');
  const buttonsWithoutText = Array.from(buttons).filter(btn => !btn.textContent?.trim());
  if (buttonsWithoutText.length > 0) {
    issues.push(`${buttonsWithoutText.length} buttons missing accessible names`);
  }

  // Check for form inputs without labels
  const inputs = document.querySelectorAll('input:not([aria-label]):not([aria-labelledby])');
  const inputsWithoutLabels = Array.from(inputs).filter(input => {
    const id = input.getAttribute('id');
    return !id || !document.querySelector(`label[for="${id}"]`);
  });
  if (inputsWithoutLabels.length > 0) {
    issues.push(`${inputsWithoutLabels.length} form inputs missing labels`);
  }

  // Check for proper heading hierarchy
  const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
  let previousLevel = 0;
  let hierarchyIssues = 0;
  
  headings.forEach(heading => {
    const level = parseInt(heading.tagName.charAt(1));
    if (level > previousLevel + 1) {
      hierarchyIssues++;
    }
    previousLevel = level;
  });
  
  if (hierarchyIssues > 0) {
    issues.push(`${hierarchyIssues} heading hierarchy violations`);
  }

  if (issues.length > 0) {
    console.group('Accessibility Issues');
    issues.forEach(issue => console.warn(issue));
    console.groupEnd();
  } else {
    console.log('✅ No accessibility issues found');
  }
}

// Responsive design checker
export function checkResponsiveDesign(): void {
  const breakpoints = [
    { name: 'Mobile', width: 375 },
    { name: 'Tablet', width: 768 },
    { name: 'Desktop', width: 1024 },
    { name: 'Large Desktop', width: 1440 }
  ];

  console.group('Responsive Design Check');
  
  breakpoints.forEach(breakpoint => {
    // This would ideally test at different viewport sizes
    console.log(`${breakpoint.name} (${breakpoint.width}px): Layout would be tested here`);
  });
  
  console.groupEnd();
}

// Global performance monitor instance
export const performanceMonitor = new PerformanceMonitor();

// Default performance budget for lesson pages
export const defaultPerformanceBudget: PerformanceBudget = {
  maxLoadTime: 3000, // 3 seconds
  maxRenderTime: 100, // 100ms
  maxInteractionTime: 50, // 50ms
  maxMemoryUsage: 50 // 50MB
};
