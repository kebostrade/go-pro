import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Header from "@/components/layout/header";
import Footer from "@/components/layout/footer";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "GO-PRO | Master Go Programming",
  description: "A comprehensive Go programming learning platform with interactive lessons, hands-on exercises, and real-world projects. Master Go from basics to microservices.",
  keywords: ["Go", "Golang", "Programming", "Learning", "Tutorial", "Exercises", "Backend", "API", "Microservices"],
  authors: [{ name: "GO-PRO Team" }],
  creator: "GO-PRO Learning Platform",
  openGraph: {
    title: "GO-PRO | Master Go Programming",
    description: "Learn Go programming through interactive lessons and hands-on projects",
    type: "website",
    locale: "en_US",
  },
  twitter: {
    card: "summary_large_image",
    title: "GO-PRO | Master Go Programming",
    description: "Learn Go programming through interactive lessons and hands-on projects",
  },
};

export function generateViewport() {
  return {
    width: 'device-width',
    initialScale: 1,
    themeColor: [
      { media: "(prefers-color-scheme: light)", color: "#00ADD8" },
      { media: "(prefers-color-scheme: dark)", color: "#00ADD8" },
    ],
  }
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased min-h-screen flex flex-col bg-background`}
        suppressHydrationWarning
      >
        {/* Skip link for accessibility */}
        <a
          href="#main-content"
          className="absolute left-[-10000px] top-auto w-1 h-1 overflow-hidden focus:left-6 focus:top-7 focus:w-auto focus:h-auto focus:overflow-visible focus:z-50 focus:bg-primary focus:text-primary-foreground focus:px-4 focus:py-2 focus:rounded transition-all"
          tabIndex={0}
        >
          Skip to main content
        </a>

        <Header />
        <main id="main-content" className="flex-1 relative" tabIndex={-1}>
          <div className="min-h-full">
            {children}
          </div>
        </main>
        <Footer />
      </body>
    </html>
  );
}
