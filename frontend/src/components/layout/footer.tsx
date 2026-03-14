import Link from "next/link";
import { Separator } from "@/components/ui/separator";
import {
  Github,
  Twitter,
  Linkedin,
  Mail,
  Heart,
} from "lucide-react";

const Footer = () => {
  const currentYear = new Date().getFullYear();

  const footerLinks = {
    learn: [
      { name: "Getting Started", href: "/learn/getting-started" },
      { name: "Go Basics", href: "/learn/basics" },
      { name: "Advanced Topics", href: "/learn/advanced" },
      { name: "Prompt Engineering", href: "/learn/prompt-engineering" },
    ],
    practice: [
      { name: "Exercises", href: "/practice/exercises" },
      { name: "Challenges", href: "/practice/challenges" },
      { name: "Code Review", href: "/practice/code-review" },
      { name: "Leaderboard", href: "/practice/leaderboard" },
    ],
    resources: [
      { name: "Documentation", href: "/docs" },
      { name: "API Reference", href: "/api" },
      { name: "Examples", href: "/examples" },
      { name: "Blog", href: "/blog" },
    ],
  };

  const socialLinks = [
    { name: "GitHub", href: "https://github.com/go-pro", icon: Github },
    { name: "Twitter", href: "https://twitter.com/gopro", icon: Twitter },
    { name: "LinkedIn", href: "https://linkedin.com/company/gopro", icon: Linkedin },
    { name: "Email", href: "mailto:hello@gopro.dev", icon: Mail },
  ];

  return (
    <footer className="border-t border-border bg-background">
      <div className="mx-auto px-4 sm:px-6 md:px-8 lg:px-12 xl:px-16 2xl:px-20 py-8 sm:py-10 lg:py-12 max-w-[1800px] w-full">
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-6 lg:gap-8">
          <div className="col-span-2 sm:col-span-3 lg:col-span-2 space-y-4">
            <Link href="/" className="inline-flex items-center gap-2 group">
              <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-gradient-to-br from-blue-600 to-cyan-500 shadow-md">
                <span className="text-lg font-bold text-white">G</span>
              </div>
              <div className="flex flex-col">
                <span className="text-lg font-bold bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent">GO-PRO</span>
                <span className="text-xs text-muted-foreground -mt-0.5">Learn Go Programming</span>
              </div>
            </Link>

            <p className="text-sm text-muted-foreground leading-relaxed max-w-md">
              Master Go programming through interactive lessons, hands-on exercises, and real-world projects.
            </p>

            <div className="flex items-center gap-3">
              {socialLinks.map((social) => (
                <Link
                  key={social.name}
                  href={social.href}
                  className="inline-flex items-center justify-center h-9 w-9 rounded-lg text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
                  target={social.href.startsWith('http') ? '_blank' : undefined}
                  rel={social.href.startsWith('http') ? 'noopener noreferrer' : undefined}
                  aria-label={social.name}
                >
                  <social.icon className="h-4 w-4" />
                </Link>
              ))}
            </div>
          </div>

          <div className="space-y-3">
            <h3 className="font-semibold text-sm">Learn</h3>
            <ul className="space-y-2">
              {footerLinks.learn.slice(0, 4).map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <h3 className="font-semibold text-sm">Practice</h3>
            <ul className="space-y-2">
              {footerLinks.practice.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <h3 className="font-semibold text-sm">Resources</h3>
            <ul className="space-y-2">
              {footerLinks.resources.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <Separator className="my-8" />

        <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
          <div className="flex flex-wrap items-center justify-center gap-x-4 gap-y-2 text-sm text-muted-foreground">
            <span>© {currentYear} GO-PRO. All rights reserved.</span>
            <Link href="/privacy" className="hover:text-foreground transition-colors">
              Privacy Policy
            </Link>
            <Link href="/terms" className="hover:text-foreground transition-colors">
              Terms of Service
            </Link>
          </div>

          <div className="flex items-center gap-1 text-sm text-muted-foreground">
            <span>Made with</span>
            <Heart className="h-4 w-4 text-red-500 fill-current" />
            <span>for the Go community</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
