import OpenClawLessonClient from './OpenClawLessonClient';

export default function Page({ params }: { params: Promise<{ id: string }> }) {
  return <OpenClawLessonClient params={params} />;
}
