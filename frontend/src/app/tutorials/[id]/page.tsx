import React from 'react';
import { notFound } from 'next/navigation';
import { getTutorialById, tutorials } from '@/lib/tutorials-data';
import { TutorialViewer } from '@/components/learning/tutorial-viewer';

export async function generateStaticParams() {
  return tutorials.map((tutorial) => ({
    id: tutorial.id,
  }));
}

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = await params;
  const tutorial = getTutorialById((resolvedParams?.id as string) || 'intro');

  if (!tutorial) {
    return {
      title: 'Tutorial Not Found',
    };
  }

  return {
    title: `${tutorial.title} | Go Pro Learning Platform`,
    description: tutorial.description,
  };
}

export default async function TutorialPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = await params;
  const tutorial = getTutorialById((resolvedParams?.id as string) || 'intro');

  if (!tutorial) {
    notFound();
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4 py-12">
        <TutorialViewer tutorial={tutorial} />
      </div>
    </div>
  );
}

