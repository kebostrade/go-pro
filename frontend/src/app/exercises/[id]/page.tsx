import ExerciseClient from './ExerciseClient';

// Generate static params for static export
export function generateStaticParams() {
  return [
    { id: '1' },
    { id: '2' },
    { id: '3' },
  ];
}

// Server component wrapper for static export
// Next.js 15: params is now a Promise
export default async function ExercisePage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <ExerciseClient exerciseId={id} />;
}
