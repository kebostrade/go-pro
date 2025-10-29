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
      <div className="container max-w-screen-2xl px-4 py-12">
        {/* Main footer content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8">
          {/* Brand section */}
          <div className="lg:col-span-2">
            <Link href="/" className="flex items-center space-x-2 mb-4">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <span className="text-lg font-bold text-primary-foreground">G</span>
              </div>
              <div className="flex flex-col">
                <span className="text-lg font-bold go-gradient-text">GO-PRO</span>
                <span className="text-xs text-muted-foreground -mt-1">Learn Go Programming</span>
              </div>
            </Link>
            <p className="text-sm text-muted-foreground mb-4 max-w-md">
              Master Go programming through interactive lessons, hands-on exercises, and real-world projects. 
              From basics to microservices, we've got you covered.
            </p>
            <div className="flex items-center space-x-2 mb-4">
              <Badge variant="secondary" className="text-xs">
                Open Source
              </Badge>
              <Badge variant="outline" className="text-xs">
                MIT License
              </Badge>
            </div>
            <div className="flex items-center space-x-3">
              {socialLinks.map((social) => (
                <Link
                  key={social.name}
                  href={social.href}
                  className="text-muted-foreground hover:text-foreground transition-colors"
                  target={social.href.startsWith('http') ? '_blank' : undefined}
                  rel={social.href.startsWith('http') ? 'noopener noreferrer' : undefined}
                >
                  <social.icon className="h-4 w-4" />
                  <span className="sr-only">{social.name}</span>
                </Link>
              ))}
            </div>
          </div>

          {/* Links sections */}
          <div>
            <h3 className="font-semibold text-sm mb-3">Learn</h3>
            <ul className="space-y-2">
              {footerLinks.learn.map((link) => (
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

          <div>
            <h3 className="font-semibold text-sm mb-3">Practice</h3>
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

          <div>
            <h3 className="font-semibold text-sm mb-3">Resources</h3>
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
            
            <h3 className="font-semibold text-sm mb-3 mt-6">Community</h3>
            <ul className="space-y-2">
              {footerLinks.community.map((link) => (
                <li key={link.name}>
                  <Link
                    href={link.href}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors flex items-center"
                    target={link.external ? '_blank' : undefined}
                    rel={link.external ? 'noopener noreferrer' : undefined}
                  >
                    {link.name}
                    {link.external && <ExternalLink className="ml-1 h-3 w-3" />}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <Separator className="my-8" />

        {/* Bottom section */}
        <div className="flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0">
          <div className="flex items-center space-x-4 text-sm text-muted-foreground">
            <span>© {currentYear} GO-PRO. All rights reserved.</span>
            <Link href="/privacy" className="hover:text-foreground transition-colors">
              Privacy Policy
            </Link>
            <Link href="/terms" className="hover:text-foreground transition-colors">
              Terms of Service
            </Link>
          </div>
          
          <div className="flex items-center space-x-1 text-sm text-muted-foreground">
            <span>Made with</span>
            <Heart className="h-3 w-3 text-red-500 fill-current" />
            <span>for the Go community</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
