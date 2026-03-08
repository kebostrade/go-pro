'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import {
  BookOpen,
  FileEdit,
  GraduationCap,
  BarChart3,
  GitBranch,
  Users,
  MessageSquare,
  UserSearch,
  Settings,
  Menu,
  X,
  ChevronLeft,
} from 'lucide-react';

const navigation = [
  { name: 'Dashboard', href: '/cms', icon: BarChart3 },
  {
    name: 'Content',
    icon: FileEdit,
    children: [
      { name: 'Lessons', href: '/cms/content/lessons' },
      { name: 'Assessments', href: '/cms/content/assessments' },
      { name: 'Media Library', href: '/cms/content/media' },
    ],
  },
  {
    name: 'Grading',
    icon: GraduationCap,
    children: [
      { name: 'Gradebook', href: '/cms/grading' },
      { name: 'Submissions', href: '/cms/grading/submissions' },
    ],
  },
  {
    name: 'Analytics',
    icon: BarChart3,
    children: [
      { name: 'Overview', href: '/cms/analytics' },
      { name: 'Students', href: '/cms/analytics/students' },
      { name: 'Performance', href: '/cms/analytics/performance' },
    ],
  },
  {
    name: 'Learning Paths',
    icon: GitBranch,
    children: [
      { name: 'All Paths', href: '/cms/paths' },
      { name: 'Create Path', href: '/cms/paths/new' },
    ],
  },
  {
    name: 'Collaboration',
    icon: Users,
    children: [
      { name: 'Peer Reviews', href: '/cms/collaboration/peer-reviews' },
      { name: 'Study Groups', href: '/cms/collaboration/study-groups' },
      { name: 'Mentors', href: '/cms/collaboration/mentors' },
    ],
  },
];

export default function CMSLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [collapsedSections, setCollapsedSections] = useState<Set<string>>(new Set());
  const pathname = usePathname();

  const toggleSection = (sectionName: string) => {
    setCollapsedSections((prev) => {
      const next = new Set(prev);
      if (next.has(sectionName)) {
        next.delete(sectionName);
      } else {
        next.add(sectionName);
      }
      return next;
    });
  };

  const isActive = (href: string) => {
    if (href === '/cms') {
      return pathname === '/cms';
    }
    return pathname.startsWith(href);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-gray-900/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-50 w-64 transform bg-white border-r border-gray-200 transition-transform duration-200 ease-in-out lg:translate-x-0 lg:static lg:inset-0',
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex h-full flex-col">
          {/* Logo */}
          <div className="flex h-16 items-center justify-between px-6 border-b border-gray-200">
            <Link href="/cms" className="flex items-center space-x-2">
              <BookOpen className="h-6 w-6 text-blue-600" />
              <span className="text-xl font-bold text-gray-900">CMS</span>
            </Link>
            <button
              onClick={() => setSidebarOpen(false)}
              className="lg:hidden text-gray-500 hover:text-gray-700"
            >
              <X className="h-6 w-6" />
            </button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 overflow-y-auto px-4 py-6 space-y-1">
            {navigation.map((item) => {
              const Icon = item.icon;
              const hasChildren = 'children' in item && item.children;
              const isCollapsed = collapsedSections.has(item.name);

              if (hasChildren) {
                return (
                  <div key={item.name}>
                    <button
                      onClick={() => toggleSection(item.name)}
                      className="flex w-full items-center justify-between rounded-lg px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 transition-colors"
                    >
                      <div className="flex items-center space-x-3">
                        <Icon className="h-5 w-5 text-gray-400" />
                        <span>{item.name}</span>
                      </div>
                      <ChevronLeft
                        className={cn(
                          'h-4 w-4 transition-transform',
                          isCollapsed ? '-rotate-90' : ''
                        )}
                      />
                    </button>
                    {!isCollapsed && (
                      <div className="mt-1 ml-8 space-y-1">
                        {item.children.map((child) => (
                          <Link
                            key={child.href}
                            href={child.href}
                            className={cn(
                              'flex rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                              isActive(child.href)
                                ? 'bg-blue-50 text-blue-700'
                                : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
                            )}
                          >
                            {child.name}
                          </Link>
                        ))}
                      </div>
                    )}
                  </div>
                );
              }

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    'flex rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                    isActive(item.href)
                      ? 'bg-blue-50 text-blue-700'
                      : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                  )}
                >
                  <div className="flex items-center space-x-3">
                    <Icon className="h-5 w-5 text-gray-400" />
                    <span>{item.name}</span>
                  </div>
                </Link>
              );
            })}
          </nav>

          {/* User section */}
          <div className="border-t border-gray-200 p-4">
            <Link
              href="/settings"
              className="flex items-center space-x-3 rounded-lg px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 transition-colors"
            >
              <Settings className="h-5 w-5 text-gray-400" />
              <span>Settings</span>
            </Link>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Top bar */}
        <header className="flex h-16 items-center justify-between border-b border-gray-200 bg-white px-4 sm:px-6 lg:px-8">
          <button
            onClick={() => setSidebarOpen(true)}
            className="text-gray-500 hover:text-gray-700 lg:hidden"
          >
            <Menu className="h-6 w-6" />
          </button>

          <div className="flex flex-1 items-center justify-end space-x-4">
            {/* Breadcrumb */}
            <nav className="hidden sm:flex" aria-label="Breadcrumb">
              <ol className="flex items-center space-x-2">
                <li>
                  <Link href="/cms" className="text-sm text-gray-500 hover:text-gray-700">
                    CMS
                  </Link>
                </li>
                {pathname !== '/cms' && (
                  <>
                    <span className="text-gray-300">/</span>
                    <li className="text-sm text-gray-700">
                      {pathname.split('/').pop()?.replace(/-/g, ' ').replace(/^\w/, c => c.toUpperCase())}
                    </li>
                  </>
                )}
              </ol>
            </nav>
          </div>
        </header>

        {/* Page content */}
        <main className="p-4 sm:p-6 lg:p-8">
          <div className="mx-auto max-w-7xl">{children}</div>
        </main>
      </div>
    </div>
  );
}
