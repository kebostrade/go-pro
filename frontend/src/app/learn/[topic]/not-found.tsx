import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { BookOpen } from 'lucide-react';

export default function NotFound() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800 flex items-center justify-center">
      <div className="text-center max-w-md mx-auto px-4">
        <BookOpen className="w-16 h-16 mx-auto text-gray-400 mb-4" />
        <h1 className="text-2xl font-bold mb-2 text-gray-900 dark:text-white">Topic Not Found</h1>
        <p className="text-muted-foreground mb-6">
          The topic you&apos;re looking for doesn&apos;t exist or hasn&apos;t been added yet.
        </p>
        <Link href="/learn">
          <Button>Back to Learning Hub</Button>
        </Link>
      </div>
    </div>
  );
}
