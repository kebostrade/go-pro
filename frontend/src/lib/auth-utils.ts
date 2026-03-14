let firebaseErrorClass: (new (code: string, message: string) => { code: string; message: string }) | null = null;

function getFirebaseErrorClass() {
  if (!firebaseErrorClass) {
    try {
      const { FirebaseError: FE } = require('firebase/app');
      firebaseErrorClass = FE;
    } catch {
      firebaseErrorClass = class FirebaseError {
        code: string;
        message: string;
        constructor(code: string, message: string) {
          this.code = code;
          this.message = message;
        }
      };
    }
  }
  return firebaseErrorClass;
}

/**
 * Enhanced Authentication Utilities
 * Provides user-friendly error messages, rate limiting, and phone auth support
 */

// Error message mapping for better UX
export const AUTH_ERROR_MESSAGES: Record<string, string> = {
  // Email/Password errors
  'auth/email-already-in-use': 'This email is already registered. Please sign in instead.',
  'auth/invalid-email': 'Please enter a valid email address.',
  'auth/user-disabled': 'This account has been disabled. Please contact support.',
  'auth/user-not-found': 'No account found with this email. Please sign up first.',
  'auth/wrong-password': 'Incorrect password. Please try again.',
  'auth/weak-password': 'Password should be at least 6 characters long.',

  // OAuth errors
  'auth/account-exists-with-different-credential': 'An account already exists with this email using a different sign-in method.',
  'auth/popup-blocked': 'Sign-in popup was blocked. Please allow popups for this site.',
  'auth/popup-closed-by-user': 'Sign-in was cancelled. Please try again.',
  'auth/cancelled-popup-request': 'Only one sign-in popup allowed at a time.',

  // Network errors
  'auth/network-request-failed': 'Network error. Please check your internet connection.',
  'auth/timeout': 'Request timed out. Please try again.',

  // Security errors
  'auth/too-many-requests': 'Too many failed attempts. Please try again later.',
  'auth/requires-recent-login': 'This action requires you to sign in again for security.',
  'auth/invalid-verification-code': 'Invalid verification code. Please try again.',
  'auth/invalid-verification-id': 'Verification failed. Please request a new code.',

  // Phone auth errors
  'auth/invalid-phone-number': 'Please enter a valid phone number.',
  'auth/missing-phone-number': 'Phone number is required.',
  'auth/quota-exceeded': 'SMS quota exceeded. Please try again later.',
  'auth/captcha-check-failed': 'reCAPTCHA verification failed. Please try again.',

  // MFA errors
  'auth/multi-factor-auth-required': 'Additional verification required.',
  'auth/maximum-second-factor-count-exceeded': 'Maximum number of second factors reached.',

  // General errors
  'auth/operation-not-allowed': 'This sign-in method is not enabled. Please contact support.',
  'auth/internal-error': 'An internal error occurred. Please try again.',
};

/**
 * Get user-friendly error message from Firebase error
 */
export function getAuthErrorMessage(error: unknown): string {
  const FirebaseError = getFirebaseErrorClass();
  
  // Check if it's a Firebase-like error by checking for 'code' property
  // Using type guard instead of instanceof for better type narrowing
  if (error !== null && typeof error === 'object' && 'code' in error) {
    const firebaseError = error as { code: string; message?: string };
    return AUTH_ERROR_MESSAGES[firebaseError.code] || firebaseError.message || 'An unexpected error occurred. Please try again.';
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'An unexpected error occurred. Please try again.';
}

/**
 * Rate limiting for authentication attempts
 */
class RateLimiter {
  private attempts: Map<string, { count: number; firstAttempt: number }> = new Map();
  private readonly maxAttempts: number = 5;
  private readonly windowMs: number = 15 * 60 * 1000; // 15 minutes

  check(identifier: string): { allowed: boolean; retryAfter?: number } {
    const now = Date.now();
    const record = this.attempts.get(identifier);

    if (!record) {
      this.attempts.set(identifier, { count: 1, firstAttempt: now });
      return { allowed: true };
    }

    // Reset if window has passed
    if (now - record.firstAttempt > this.windowMs) {
      this.attempts.set(identifier, { count: 1, firstAttempt: now });
      return { allowed: true };
    }

    // Check if limit exceeded
    if (record.count >= this.maxAttempts) {
      const retryAfter = Math.ceil((record.firstAttempt + this.windowMs - now) / 1000);
      return { allowed: false, retryAfter };
    }

    // Increment counter
    record.count++;
    return { allowed: true };
  }

  reset(identifier: string): void {
    this.attempts.delete(identifier);
  }
}

export const authRateLimiter = new RateLimiter();

/**
 * Session persistence options
 */
export enum SessionPersistence {
  LOCAL = 'local', // Persists even when browser is closed
  SESSION = 'session', // Persists only while tab is open
  NONE = 'none', // No persistence
}

/**
 * Phone number formatting and validation
 */
export class PhoneValidator {
  private static readonly PHONE_REGEX = /^\+[1-9]\d{1,14}$/;

  /**
   * Validate phone number format (E.164)
   */
  static isValid(phone: string): boolean {
    return this.PHONE_REGEX.test(phone);
  }

  /**
   * Format phone number to E.164 format
   * Example: (555) 123-4567 → +15551234567
   */
  static format(phone: string, countryCode: string = '+1'): string {
    // Remove all non-digit characters
    const digits = phone.replace(/\D/g, '');

    // Add country code if not present
    if (!phone.startsWith('+')) {
      return `${countryCode}${digits}`;
    }

    return `+${digits}`;
  }

  /**
   * Mask phone number for display
   * Example: +15551234567 → +1 (555) ***-4567
   */
  static mask(phone: string): string {
    if (!this.isValid(phone)) return phone;

    const cleaned = phone.replace(/\D/g, '');
    const countryCode = phone.match(/^\+(\d{1,3})/)?.[1] || '1';
    const localNumber = cleaned.slice(countryCode.length);

    if (localNumber.length === 10) {
      return `+${countryCode} (${localNumber.slice(0, 3)}) ***-${localNumber.slice(6)}`;
    }

    return phone;
  }
}

/**
 * Password strength validator
 */
export class PasswordValidator {
  private static readonly MIN_LENGTH = 8;
  private static readonly STRENGTH_REGEX = {
    weak: /^.{6,7}$/,
    medium: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/,
    strong: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/,
  };

  /**
   * Check password strength
   */
  static getStrength(password: string): 'weak' | 'medium' | 'strong' {
    if (this.STRENGTH_REGEX.strong.test(password)) return 'strong';
    if (this.STRENGTH_REGEX.medium.test(password)) return 'medium';
    return 'weak';
  }

  /**
   * Validate password meets minimum requirements
   */
  static isValid(password: string): { valid: boolean; message?: string } {
    if (password.length < this.MIN_LENGTH) {
      return { valid: false, message: `Password must be at least ${this.MIN_LENGTH} characters` };
    }

    if (!/[a-z]/.test(password)) {
      return { valid: false, message: 'Password must contain a lowercase letter' };
    }

    if (!/[A-Z]/.test(password)) {
      return { valid: false, message: 'Password must contain an uppercase letter' };
    }

    if (!/\d/.test(password)) {
      return { valid: false, message: 'Password must contain a number' };
    }

    return { valid: true };
  }

  /**
   * Generate strength indicator UI data
   */
  static getStrengthIndicator(password: string): {
    strength: 'weak' | 'medium' | 'strong';
    percentage: number;
    color: string;
    label: string;
  } {
    const strength = this.getStrength(password);

    const indicators = {
      weak: { percentage: 33, color: 'red', label: 'Weak' },
      medium: { percentage: 66, color: 'yellow', label: 'Medium' },
      strong: { percentage: 100, color: 'green', label: 'Strong' },
    };

    return { strength, ...indicators[strength] };
  }
}

/**
 * Email validation
 */
export class EmailValidator {
  private static readonly EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

  static isValid(email: string): boolean {
    return this.EMAIL_REGEX.test(email);
  }

  /**
   * Check if email is from a disposable email provider
   */
  static isDisposable(email: string): boolean {
    const disposableDomains = [
      'tempmail.com',
      'guerrillamail.com',
      'mailinator.com',
      '10minutemail.com',
      'throwaway.email',
    ];

    const domain = email.split('@')[1]?.toLowerCase();
    return disposableDomains.includes(domain);
  }
}

/**
 * Authentication analytics
 */
export class AuthAnalytics {
  /**
   * Track authentication event
   */
  static trackEvent(event: string, properties?: Record<string, any>): void {
    if (typeof window !== 'undefined' && (window as any).gtag) {
      (window as any).gtag('event', event, properties);
    }

    console.log(`[Auth Analytics] ${event}`, properties);
  }

  /**
   * Track successful signin
   */
  static trackSignIn(method: 'email' | 'google' | 'github' | 'phone'): void {
    this.trackEvent('sign_in', { method });
  }

  /**
   * Track signup
   */
  static trackSignUp(method: 'email' | 'google' | 'github' | 'phone'): void {
    this.trackEvent('sign_up', { method });
  }

  /**
   * Track authentication error
   */
  static trackError(error: string, context: string): void {
    this.trackEvent('auth_error', { error, context });
  }
}
