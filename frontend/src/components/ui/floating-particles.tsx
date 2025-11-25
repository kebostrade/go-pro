"use client";

import React from 'react';

interface FloatingParticlesProps {
  count?: number;
  className?: string;
}

export function FloatingParticles({ count = 20, className = "" }: FloatingParticlesProps) {
  const particles = Array.from({ length: count }, (_, i) => ({
    id: i,
    size: Math.random() * 4 + 2,
    left: Math.random() * 100,
    animationDuration: Math.random() * 20 + 10,
    animationDelay: Math.random() * 5,
    opacity: Math.random() * 0.5 + 0.1,
  }));

  return (
    <div className={`absolute inset-0 overflow-hidden pointer-events-none ${className}`}>
      {particles.map((particle) => (
        <div
          key={particle.id}
          className="absolute rounded-full bg-gradient-to-br from-primary/30 to-blue-500/30 blur-sm"
          style={{
            width: `${particle.size}px`,
            height: `${particle.size}px`,
            left: `${particle.left}%`,
            bottom: '-10%',
            opacity: particle.opacity,
            animation: `float-up ${particle.animationDuration}s linear infinite`,
            animationDelay: `${particle.animationDelay}s`,
          }}
        />
      ))}
      <style jsx>{`
        @keyframes float-up {
          0% {
            transform: translateY(0) translateX(0);
            opacity: 0;
          }
          10% {
            opacity: ${Math.random() * 0.5 + 0.3};
          }
          90% {
            opacity: ${Math.random() * 0.5 + 0.3};
          }
          100% {
            transform: translateY(-100vh) translateX(${Math.random() * 100 - 50}px);
            opacity: 0;
          }
        }
      `}</style>
    </div>
  );
}

interface AnimatedBackgroundProps {
  variant?: 'gradient' | 'mesh' | 'dots';
  className?: string;
}

export function AnimatedBackground({ variant = 'gradient', className = "" }: AnimatedBackgroundProps) {
  if (variant === 'gradient') {
    return (
      <div className={`absolute inset-0 overflow-hidden pointer-events-none ${className}`}>
        <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-blue-500/5 animate-gradient-shift" />
        <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-br from-cyan-500/10 to-transparent rounded-full blur-3xl float-animation" />
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-gradient-to-tr from-blue-500/10 to-transparent rounded-full blur-3xl float-animation" style={{ animationDelay: '2s' }} />
      </div>
    );
  }

  if (variant === 'mesh') {
    return (
      <div className={`absolute inset-0 overflow-hidden pointer-events-none ${className}`}>
        <div className="absolute inset-0 opacity-30">
          <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-primary/20 via-transparent to-blue-500/20" 
               style={{ 
                 backgroundImage: `
                   linear-gradient(to right, rgba(var(--primary), 0.1) 1px, transparent 1px),
                   linear-gradient(to bottom, rgba(var(--primary), 0.1) 1px, transparent 1px)
                 `,
                 backgroundSize: '50px 50px'
               }} 
          />
        </div>
      </div>
    );
  }

  if (variant === 'dots') {
    return (
      <div className={`absolute inset-0 overflow-hidden pointer-events-none ${className}`}>
        <div className="absolute inset-0 opacity-20"
             style={{
               backgroundImage: 'radial-gradient(circle, rgba(var(--primary), 0.3) 1px, transparent 1px)',
               backgroundSize: '30px 30px'
             }}
        />
      </div>
    );
  }

  return null;
}

interface GlowingOrbProps {
  color?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  position?: { top?: string; bottom?: string; left?: string; right?: string };
  className?: string;
}

export function GlowingOrb({ 
  color = 'primary', 
  size = 'md', 
  position = { top: '50%', left: '50%' },
  className = "" 
}: GlowingOrbProps) {
  const sizes = {
    sm: 'w-32 h-32',
    md: 'w-64 h-64',
    lg: 'w-96 h-96',
    xl: 'w-[500px] h-[500px]'
  };

  const colorClasses = {
    primary: 'from-primary/30 to-primary/10',
    blue: 'from-blue-500/30 to-blue-500/10',
    cyan: 'from-cyan-500/30 to-cyan-500/10',
    green: 'from-green-500/30 to-green-500/10',
    yellow: 'from-yellow-500/30 to-yellow-500/10',
    orange: 'from-orange-500/30 to-orange-500/10',
    purple: 'from-purple-500/30 to-purple-500/10',
  };

  return (
    <div 
      className={`absolute ${sizes[size]} bg-gradient-to-br ${colorClasses[color as keyof typeof colorClasses] || colorClasses.primary} rounded-full blur-3xl float-animation pointer-events-none ${className}`}
      style={position}
    />
  );
}

export default FloatingParticles;

