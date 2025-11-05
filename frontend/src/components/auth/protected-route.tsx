'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/auth-context';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireEmailVerification?: boolean;
  requiredRole?: 'student' | 'instructor' | 'admin';
  fallbackPath?: string;
}

export function ProtectedRoute({
  children,
  requireEmailVerification = false,
  requiredRole,
  fallbackPath = '/auth/signin',
}: ProtectedRouteProps) {
  const { user, userProfile, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading) {
      // Not authenticated
      if (!user) {
        router.push(fallbackPath);
        return;
      }

      // Email verification required
      if (requireEmailVerification && !user.emailVerified) {
        router.push('/auth/verify-email');
        return;
      }

      // Role check
      if (requiredRole && userProfile) {
        const hasRequiredRole = checkRole(userProfile.role, requiredRole);
        if (!hasRequiredRole) {
          router.push('/unauthorized');
          return;
        }
      }
    }
  }, [user, userProfile, loading, requireEmailVerification, requiredRole, fallbackPath, router]);

  // Show loading state
  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  // Not authenticated
  if (!user) {
    return null;
  }

  // Email not verified
  if (requireEmailVerification && !user.emailVerified) {
    return null;
  }

  // Insufficient role
  if (requiredRole && userProfile && !checkRole(userProfile.role, requiredRole)) {
    return null;
  }

  return <>{children}</>;
}

function checkRole(userRole: string, requiredRole: string): boolean {
  const roleHierarchy = {
    student: 1,
    instructor: 2,
    admin: 3,
  };

  const userLevel = roleHierarchy[userRole as keyof typeof roleHierarchy] || 0;
  const requiredLevel = roleHierarchy[requiredRole as keyof typeof roleHierarchy] || 0;

  return userLevel >= requiredLevel;
}
