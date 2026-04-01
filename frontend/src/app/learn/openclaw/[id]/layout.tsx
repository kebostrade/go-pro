const lessons = [
  { id: 'OC-01' }, { id: 'OC-02' }, { id: 'OC-03' }, { id: 'OC-04' },
  { id: 'OC-05' }, { id: 'OC-06' }, { id: 'OC-07' }, { id: 'OC-08' },
  { id: 'OC-09' }, { id: 'OC-10' }, { id: 'OC-11' },
];

export async function generateStaticParams() {
  return lessons.map((lesson) => ({
    id: lesson.id,
  }));
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}
