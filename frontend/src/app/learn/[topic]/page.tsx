import { notFound } from 'next/navigation';
import { topics, getTopicById } from '@/lib/topics-data';
import TopicViewer from '@/components/learning/topic-viewer';

export async function generateStaticParams() {
  return topics.map((topic) => ({
    topic: topic.id,
  }));
}

interface PageProps {
  params: Promise<{ topic: string }>;
}

export default async function TopicPage({ params }: PageProps) {
  const { topic: topicId } = await params;
  const topic = getTopicById(topicId);
  
  if (!topic) {
    notFound();
  }
  
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4 py-8">
        <TopicViewer topic={topic} />
      </div>
    </div>
  );
}

export async function generateMetadata({ params }: PageProps) {
  const { topic: topicId } = await params;
  const topic = getTopicById(topicId);
  
  if (!topic) return { title: 'Topic Not Found' };
  
  return {
    title: `${topic.title} | GO-PRO`,
    description: topic.description,
  };
}
