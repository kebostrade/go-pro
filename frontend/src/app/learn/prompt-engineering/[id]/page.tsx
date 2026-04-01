import PromptEngineeringLessonClient from './PromptEngineeringLessonClient';

export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return <PromptEngineeringLessonClient params={params} />;
}
