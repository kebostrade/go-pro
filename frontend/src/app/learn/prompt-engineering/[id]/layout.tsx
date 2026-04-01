const lessons = [
  { id: "PE-01" }, { id: "PE-02" }, { id: "PE-03" }, { id: "PE-04" },
  { id: "PE-05" }, { id: "PE-06" }, { id: "PE-07" }, { id: "PE-08" },
  { id: "PE-09" }, { id: "PE-10" }, { id: "PE-11" }, { id: "PE-12" },
  { id: "PE-13" }, { id: "PE-14" }, { id: "PE-15" }, { id: "PE-16" },
  { id: "PE-17" }, { id: "PE-18" }, { id: "PE-19" }, { id: "PE-20" },
];

export async function generateStaticParams() {
  return lessons.map((lesson) => ({
    id: lesson.id,
  }));
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}
