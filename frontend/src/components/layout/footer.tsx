import Link from "next/link";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import {
  Github,
  Twitter,
  Linkedin,
  Mail,
  Heart,
  ExternalLink,
  Sparkles,
} from "lucide-react";

const Footer = () => {
  const currentYear = new Date().getFullYear();

  const footerLinks = {
    learn: [
      { name: "Getting Started", href: "/learn/getting-started" },
      { name: "Go Basics", href: "/learn/basics" },
      { name: "Advanced Topics", href: "/learn/advanced" },
      { name: "Best Practices", href: "/learn/best-practices" },
      { name: "Prompt Engineering", href: "/learn/prompt-engineering", badge: "New" },
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
    community: [
      { name: "Discord", href: "https://discord.gg/golang", external: true },
      { name: "Forum", href: "/community/forum" },
      { name: "Contributors", href: "/community/contributors" },
      { name: "Events", href: "/community/events" },
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
      <div className="container max-w-screen-3xl px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-16 3xl:px-24 py-8 sm:py-10 lg:py-12 xl:py-14 2xl:py-16 3xl:py-20">
        {/* Main footer content - Enhanced grid for big screens */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-6 3xl:grid-cols-7 gap-6 sm:gap-8 lg:gap-10 xl:gap-12 2xl:gap-14 3xl:gap-16">
          {/* Brand section - Spans more columns on big screens */}
          <div className="lg:col-span-2 xl:col-span-2 2xl:col-span-2 3xl:col-span-3 space-y-4 xl:space-y-5 2xl:space-y-6">
            <Link href="/" className="inline-flex items-center space-x-2 xl:space-x-3 2xl:space-x-4 group">
              <div className="flex h-8 w-8 sm:h-9 sm:w-9 xl:h-10 xl:w-10 2xl:h-12 2xl:w-12 3xl:h-14 3xl:w-14 items-center justify-center rounded-lg bg-primary group-hover:shadow-lg group-hover:shadow-primary/30 transition-shadow">
                <span className="text-base sm:text-lg xl:text-xl 2xl:text-2xl 3xl:text-3xl font-bold text-primary-foreground">G</span>
              </div>
              <div className="flex flex-col">
                <span className="text-base sm:text-lg xl:text-xl 2xl:text-2xl 3xl:text-3xl font-bold go-gradient-text">GO-PRO</span>
                <span className="text-[10px] sm:text-xs xl:text-sm 2xl:text-base 3xl:text-lg text-muted-foreground -mt-0.5">Learn Go Programming</span>
              </div>
            </Link>

            <p className="text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground leading-relaxed max-w-md xl:max-w-lg 2xl:max-w-xl 3xl:max-w-2xl">
              Master Go programming through interactive lessons, hands-on exercises, and real-world projects.
              From basics to microservices, we've got you covered.
            </p>

            <div className="flex flex-wrap items-center gap-2 xl:gap-3">
              <Badge variant="secondary" className="text-[10px] sm:text-xs xl:text-sm 2xl:text-base">
                Open Source
              </Badge>
              <Badge variant="outline" className="text-[10px] sm:text-xs xl:text-sm 2xl:text-base">
                MIT License
              </Badge>
            </div>

            {/* Social links with better mobile touch targets */}
            <div className="flex items-center gap-2 sm:gap-3 xl:gap-4 2xl:gap-5">
              {socialLinks.map((social) => (
                <Link
                  key={social.name}
                  href={social.href}
                  className="inline-flex items-center justify-center h-9 w-9 sm:h-10 sm:w-10 xl:h-11 xl:w-11 2xl:h-12 2xl:w-12 3xl:h-14 3xl:w-14 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-all"
                  target={social.href.startsWith('http') ? '_blank' : undefined}
                  rel={social.href.startsWith('http') ? 'noopener noreferrer' : undefined}
                  aria-label={social.name}
                >
                  <social.icon className="h-4 w-4 sm:h-[18px] sm:w-[18px] xl:h-5 xl:w-5 2xl:h-6 2xl:w-6 3xl:h-7 3xl:w-7" />
                  <span className="sr-only">{social.name}</span>
                </Link>
              ))}
            </div>
          </div>

          {/* Links sections - Enhanced for big screens */}
          <div className="space-y-3 xl:space-y-4">
            <h3 className="font-semibold text-sm sm:text-base xl:text-lg 2xl:text-xl 3xl:text-2xl mb-3 xl:mb-4">Learn</h3>
            <ul className="space-y-2 sm:space-y-2.5 xl:space-y-3 2xl:space-y-4">
              {footerLinks.learn.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="inline-flex items-center gap-1.5 text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground hover:text-foreground transition-colors py-1"
                  >
                    {link.name}
                    {link.badge && (
                      <Badge variant="secondary" className="text-[10px] sm:text-xs xl:text-sm 2xl:text-base px-1.5 py-0 xl:px-2 xl:py-0.5 2xl:px-2.5 2xl:py-1 bg-primary/10 text-primary">
                        {link.badge}
                      </Badge>
                    )}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3 xl:space-y-4">
            <h3 className="font-semibold text-sm sm:text-base xl:text-lg 2xl:text-xl 3xl:text-2xl mb-3 xl:mb-4">Practice</h3>
            <ul className="space-y-2 sm:space-y-2.5 xl:space-y-3 2xl:space-y-4">
              {footerLinks.practice.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="inline-block text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground hover:text-foreground transition-colors py-1"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3 xl:space-y-4">
            <h3 className="font-semibold text-sm sm:text-base xl:text-lg 2xl:text-xl 3xl:text-2xl mb-3 xl:mb-4">Resources</h3>
            <ul className="space-y-2 sm:space-y-2.5 xl:space-y-3 2xl:space-y-4">
              {footerLinks.resources.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="inline-block text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground hover:text-foreground transition-colors py-1"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>

            <div className="pt-3 xl:pt-4">
              <h3 className="font-semibold text-sm sm:text-base xl:text-lg 2xl:text-xl 3xl:text-2xl mb-3 xl:mb-4">Community</h3>
              <ul className="space-y-2 sm:space-y-2.5 xl:space-y-3 2xl:space-y-4">
                {footerLinks.community.map((link) => (
                  <li key={link.name}>
                    <Link
                      href={link.href}
                      className="inline-flex items-center text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground hover:text-foreground transition-colors py-1"
                      target={link.external ? '_blank' : undefined}
                      rel={link.external ? 'noopener noreferrer' : undefined}
                    >
                      {link.name}
                      {link.external && <ExternalLink className="ml-1 h-2.5 w-2.5 sm:h-3 sm:w-3 xl:h-3.5 xl:w-3.5 2xl:h-4 2xl:w-4" />}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>

        <Separator className="my-6 sm:my-8 xl:my-10 2xl:my-12 3xl:my-16" />

        {/* Bottom section - Fully responsive for big screens */}
        <div className="flex flex-col sm:flex-row justify-between items-center gap-4 sm:gap-6 xl:gap-8 2xl:gap-10 3xl:gap-12">
          <div className="flex flex-col sm:flex-row items-center gap-2 sm:gap-4 xl:gap-6 2xl:gap-8 text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground text-center sm:text-left">
            <span className="whitespace-nowrap">© {currentYear} GO-PRO. All rights reserved.</span>
            <div className="flex items-center gap-3 sm:gap-4 xl:gap-6 2xl:gap-8">
              <Link
                href="/privacy"
                className="hover:text-foreground transition-colors py-1 whitespace-nowrap"
              >
                Privacy Policy
              </Link>
              <span className="hidden sm:inline text-muted-foreground/50">•</span>
              <Link
                href="/terms"
                className="hover:text-foreground transition-colors py-1 whitespace-nowrap"
              >
                Terms of Service
              </Link>
            </div>
          </div>

          <div className="flex items-center gap-1 text-xs sm:text-sm xl:text-base 2xl:text-lg 3xl:text-xl text-muted-foreground">
            <span>Made with</span>
            <Heart className="h-3 w-3 sm:h-3.5 sm:w-3.5 xl:h-4 xl:w-4 2xl:h-5 2xl:w-5 3xl:h-6 3xl:w-6 text-red-500 fill-current mx-0.5" />
            <span>for the Go community</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
