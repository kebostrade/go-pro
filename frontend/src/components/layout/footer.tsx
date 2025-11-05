import Link from "next/link";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { 
  Github, 
  Twitter, 
  Linkedin, 
  Mail,
  Heart,
  ExternalLink
} from "lucide-react";

const Footer = () => {
  const currentYear = new Date().getFullYear();

  const footerLinks = {
    learn: [
      { name: "Getting Started", href: "/learn/getting-started" },
      { name: "Go Basics", href: "/learn/basics" },
      { name: "Advanced Topics", href: "/learn/advanced" },
      { name: "Best Practices", href: "/learn/best-practices" },
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
      <div className="container max-w-screen-2xl px-4 sm:px-6 lg:px-8 py-8 sm:py-10 lg:py-12">
        {/* Main footer content */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-6 sm:gap-8 lg:gap-10">
          {/* Brand section */}
          <div className="lg:col-span-2 space-y-4">
            <Link href="/" className="inline-flex items-center space-x-2 group">
              <div className="flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-lg bg-primary group-hover:shadow-lg group-hover:shadow-primary/30 transition-shadow">
                <span className="text-base sm:text-lg font-bold text-primary-foreground">G</span>
              </div>
              <div className="flex flex-col">
                <span className="text-base sm:text-lg font-bold go-gradient-text">GO-PRO</span>
                <span className="text-[10px] sm:text-xs text-muted-foreground -mt-0.5">Learn Go Programming</span>
              </div>
            </Link>

            <p className="text-xs sm:text-sm text-muted-foreground leading-relaxed max-w-md">
              Master Go programming through interactive lessons, hands-on exercises, and real-world projects.
              From basics to microservices, we've got you covered.
            </p>

            <div className="flex flex-wrap items-center gap-2">
              <Badge variant="secondary" className="text-[10px] sm:text-xs">
                Open Source
              </Badge>
              <Badge variant="outline" className="text-[10px] sm:text-xs">
                MIT License
              </Badge>
            </div>

            {/* Social links with better mobile touch targets */}
            <div className="flex items-center gap-2 sm:gap-3">
              {socialLinks.map((social) => (
                <Link
                  key={social.name}
                  href={social.href}
                  className="inline-flex items-center justify-center h-9 w-9 sm:h-10 sm:w-10 rounded-md text-muted-foreground hover:text-foreground hover:bg-accent transition-all"
                  target={social.href.startsWith('http') ? '_blank' : undefined}
                  rel={social.href.startsWith('http') ? 'noopener noreferrer' : undefined}
                  aria-label={social.name}
                >
                  <social.icon className="h-4 w-4 sm:h-[18px] sm:w-[18px]" />
                  <span className="sr-only">{social.name}</span>
                </Link>
              ))}
            </div>
          </div>

          {/* Links sections - Better mobile layout */}
          <div className="space-y-3">
            <h3 className="font-semibold text-sm sm:text-base mb-3">Learn</h3>
            <ul className="space-y-2 sm:space-y-2.5">
              {footerLinks.learn.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="inline-block text-xs sm:text-sm text-muted-foreground hover:text-foreground transition-colors py-1"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <h3 className="font-semibold text-sm sm:text-base mb-3">Practice</h3>
            <ul className="space-y-2 sm:space-y-2.5">
              {footerLinks.practice.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="inline-block text-xs sm:text-sm text-muted-foreground hover:text-foreground transition-colors py-1"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <h3 className="font-semibold text-sm sm:text-base mb-3">Resources</h3>
            <ul className="space-y-2 sm:space-y-2.5">
              {footerLinks.resources.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="inline-block text-xs sm:text-sm text-muted-foreground hover:text-foreground transition-colors py-1"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>

            <div className="pt-3">
              <h3 className="font-semibold text-sm sm:text-base mb-3">Community</h3>
              <ul className="space-y-2 sm:space-y-2.5">
                {footerLinks.community.map((link) => (
                  <li key={link.name}>
                    <Link
                      href={link.href}
                      className="inline-flex items-center text-xs sm:text-sm text-muted-foreground hover:text-foreground transition-colors py-1"
                      target={link.external ? '_blank' : undefined}
                      rel={link.external ? 'noopener noreferrer' : undefined}
                    >
                      {link.name}
                      {link.external && <ExternalLink className="ml-1 h-2.5 w-2.5 sm:h-3 sm:w-3" />}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>

        <Separator className="my-6 sm:my-8" />

        {/* Bottom section - Fully responsive */}
        <div className="flex flex-col sm:flex-row justify-between items-center gap-4 sm:gap-6">
          <div className="flex flex-col sm:flex-row items-center gap-2 sm:gap-4 text-xs sm:text-sm text-muted-foreground text-center sm:text-left">
            <span className="whitespace-nowrap">© {currentYear} GO-PRO. All rights reserved.</span>
            <div className="flex items-center gap-3 sm:gap-4">
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

          <div className="flex items-center gap-1 text-xs sm:text-sm text-muted-foreground">
            <span>Made with</span>
            <Heart className="h-3 w-3 sm:h-3.5 sm:w-3.5 text-red-500 fill-current mx-0.5" />
            <span>for the Go community</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
