/**
 * Password validation utilities
 * Validates that passwords meet security requirements
 */

export function validatePassword(pwd: string): string | null {
  if (pwd.length < 8) {
    return 'Password must be at least 8 characters long.';
  }
  if (!/[A-Z]/.test(pwd)) {
    return 'Password must contain at least one uppercase letter.';
  }
  if (!/[a-z]/.test(pwd)) {
    return 'Password must contain at least one lowercase letter.';
  }
  if (!/[0-9]/.test(pwd)) {
    return 'Password must contain at least one number.';
  }
  return null;
}
