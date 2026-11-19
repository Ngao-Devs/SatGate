import Link from "next/link";
import { ArrowRight, ShieldCheck, Zap, Mail } from "lucide-react";
import { PrimaryButton } from "@/components/PrimaryButton";

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col bg-white text-satBlack">
      <header className="mx-auto flex w-full max-w-7xl items-center justify-between px-4 py-6 sm:px-6 lg:px-8">
        <div className="flex items-center gap-2">
          {/* eslint-disable-next-line @next/next/no-html-link-for-pages */}
          <a
            href="/"
            className="brand-font text-3xl font-extrabold tracking-tight text-satBlack hover:text-satBlue transition"
          >
            SatGo
          </a>
        </div>
        <nav className="hidden md:block">
          <ul className="brand-font flex gap-8 text-lg font-bold">
            <li>
              <a href="#features" className="hover:text-satBlue transition">
                Features
              </a>
            </li>
            <li>
              <a href="#how-it-works" className="hover:text-satBlue transition">
                How it works
              </a>
            </li>
            <li>
              <Link href="/demo" className="hover:text-satBlue transition">
                Demo
              </Link>
            </li>
          </ul>
        </nav>
        <div className="flex items-center gap-4">
          <Link href="/dashboard">
            <PrimaryButton className="brand-font">
              Get Started
              <ArrowRight size={16} />
            </PrimaryButton>
          </Link>
        </div>
      </header>

      <main className="flex-grow">
        {/* Hero Section */}
        <section className="relative overflow-hidden py-24 sm:py-32 lg:py-40">
          {/* Antigravity-inspired background */}
          <div className="absolute inset-0 -z-10 h-full w-full bg-[linear-gradient(to_right,#e2e8f0_1px,transparent_1px),linear-gradient(to_bottom,#e2e8f0_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)] opacity-70">
            {/* Floating gradient blobs */}
            <div className="absolute -top-40 -right-40 h-[600px] w-[600px] rounded-full bg-blue-500/10 blur-[130px] md:animate-pulse duration-5000" />
            <div className="absolute top-1/4 left-1/4 h-[400px] w-[400px] rounded-full bg-indigo-500/10 blur-[110px] md:animate-pulse duration-7000" />
          </div>

          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center">
            <h1 className="brand-font text-5xl font-extrabold tracking-tight sm:text-7xl leading-[1.2] sm:leading-[1.1]">
              Stop{" "}
              <span className="inline-block md:animate-bounce-custom underline decoration-[8px] decoration-satRed/80 underline-offset-[10px]">
                spam
              </span>{" "}
              with{" "}
              <span className="inline-block md:animate-bounce-custom delay-200 text-satBlue underline decoration-[8px] decoration-satBlue/80 underline-offset-[10px]">
                Sats
              </span>
              .
            </h1>
            <p className="mx-auto mt-8 max-w-2xl text-lg sm:text-xl leading-8 text-slate-600">
              SatGo adds a tiny cost to your contact form. Real leads won&apos;t mind paying 5 sats,
              but automated spammers will find it prohibitively expensive.
            </p>
            <div className="mt-12 flex items-center justify-center gap-x-6">
              <Link href="/dashboard">
                <PrimaryButton className="h-14 px-8 text-lg rounded-xl shadow-lg shadow-satBlue/15 hover:shadow-satBlue/25 transition">
                  Protect your inbox
                  <ArrowRight size={20} />
                </PrimaryButton>
              </Link>
              <Link
                href="/demo"
                className="text-sm font-bold leading-6 text-satBlack hover:text-satBlue transition flex items-center gap-1.5"
              >
                Live demo <span aria-hidden="true">→</span>
              </Link>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section id="features" className="bg-slate-50 py-24 sm:py-32">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="brand-font text-3xl font-bold tracking-tight sm:text-4xl">
                Everything you need to secure your forms
              </h2>
              <p className="mt-4 text-lg text-slate-600">
                Built for independent founders and creators who are tired of sorting through junk
                mail.
              </p>
            </div>
            <div className="mx-auto mt-16 max-w-7xl sm:mt-20 lg:mt-24">
              <dl className="grid max-w-xl grid-cols-1 gap-x-8 gap-y-10 lg:max-w-none lg:grid-cols-3">
                <Feature
                  icon={<ShieldCheck className="text-satBlue" size={24} />}
                  title="Proof of Intent"
                  description="A paid message is a verified message. Require an economic commitment before anyone reaches your inbox."
                />
                <Feature
                  icon={<Zap className="text-satBlue" size={24} />}
                  title="Lightning Fast"
                  description="Powered by the Bitcoin Lightning Network. Payments are instant, global, and have near-zero fees."
                />
                <Feature
                  icon={<Mail className="text-satBlue" size={24} />}
                  title="Easy Integration"
                  description="Drop our iframe widget into any site—WordPress, Webflow, or custom React apps—in seconds."
                />
              </dl>
            </div>
          </div>
        </section>

        {/* How it works */}
        <section id="how-it-works" className="py-24 sm:py-32">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="grid grid-cols-1 gap-x-12 gap-y-16 lg:grid-cols-2 lg:items-center">
              <div>
                <h2 className="brand-font text-3xl font-bold tracking-tight sm:text-4xl">
                  Simple for you, simple for your leads.
                </h2>
                <div className="mt-10 space-y-8">
                  <Step
                    number="1"
                    title="Create a protected form"
                    description="Set your price (e.g., 5-10 sats) and generate your unique embed code."
                  />
                  <Step
                    number="2"
                    title="Embed on your site"
                    description="Paste the iframe snippet into your contact page. It works everywhere."
                  />
                  <Step
                    number="3"
                    title="Receive verified messages"
                    description="Paid messages appear instantly in your SatGo dashboard. No more spam."
                  />
                </div>
              </div>
              <div className="rounded-2xl border border-slate-200 bg-white p-2 shadow-2xl">
                <div className="overflow-hidden rounded-xl border border-slate-100 bg-slate-50 p-4">
                  {/* Mock Dashboard UI */}
                  <div className="flex flex-col gap-4">
                    {/* Top Bar */}
                    <div className="flex items-center justify-between border-b border-slate-200 pb-3">
                      <div className="flex items-center gap-2">
                        <div className="h-3 w-3 rounded-full bg-red-400"></div>
                        <div className="h-3 w-3 rounded-full bg-yellow-400"></div>
                        <div className="h-3 w-3 rounded-full bg-green-400"></div>
                      </div>
                      <div className="flex items-center gap-2">
                        <div className="h-4 w-4 rounded-full bg-slate-200"></div>
                        <div className="h-3 w-24 rounded bg-slate-200"></div>
                      </div>
                    </div>

                    <div className="flex gap-4">
                      {/* Sidebar Mock */}
                      <div className="flex w-24 sm:w-32 flex-col gap-3 shrink-0">
                        <div className="h-3 w-16 rounded bg-slate-300 mb-1"></div>
                        <div className="h-10 w-full rounded-md bg-satBlue/10 border border-satBlue/30 p-2">
                          <div className="h-2 w-12 rounded bg-satBlue/40 mb-1"></div>
                          <div className="h-1.5 w-8 rounded bg-satBlue/20"></div>
                        </div>
                        <div className="h-10 w-full rounded-md bg-white border border-slate-200 p-2">
                          <div className="h-2 w-10 rounded bg-slate-200 mb-1"></div>
                          <div className="h-1.5 w-6 rounded bg-slate-100"></div>
                        </div>
                        <div className="mt-2 h-3 w-20 rounded bg-slate-300 mb-1"></div>
                        <div className="h-20 w-full rounded-md border border-dashed border-slate-300"></div>
                      </div>

                      {/* Main Content Mock */}
                      <div className="flex-1 flex flex-col gap-4">
                        {/* Metrics */}
                        <div className="grid grid-cols-3 gap-2">
                          <div className="h-14 rounded-lg bg-white border border-slate-200 p-2 shadow-sm">
                            <div className="h-1.5 w-8 rounded bg-slate-100 mb-2"></div>
                            <div className="h-3 w-6 rounded bg-satBlue/30"></div>
                          </div>
                          <div className="h-14 rounded-lg bg-white border border-slate-200 p-2 shadow-sm">
                            <div className="h-1.5 w-10 rounded bg-slate-100 mb-2"></div>
                            <div className="h-3 w-8 rounded bg-satGreen/30"></div>
                          </div>
                          <div className="h-14 rounded-lg bg-white border border-slate-200 p-2 shadow-sm">
                            <div className="h-1.5 w-8 rounded bg-slate-100 mb-2"></div>
                            <div className="h-3 w-4 rounded bg-satBlue/30"></div>
                          </div>
                        </div>

                        {/* Inbox Content Overview */}
                        <div className="flex flex-col gap-2">
                          <div className="h-3 w-28 rounded bg-slate-300 mb-1"></div>
                          <div className="rounded-lg bg-white border border-slate-200 p-3 shadow-sm">
                            <div className="flex items-center justify-between mb-2">
                              <div className="h-3 w-16 rounded bg-slate-200"></div>
                              <div className="h-3 w-10 rounded-full bg-satGreen/10"></div>
                            </div>
                            <div className="space-y-1.5">
                              <div className="h-2 w-full rounded bg-slate-50"></div>
                              <div className="h-2 w-5/6 rounded bg-slate-50"></div>
                            </div>
                          </div>
                          <div className="rounded-lg bg-white border border-slate-200 p-3 shadow-sm">
                            <div className="flex items-center justify-between mb-2">
                              <div className="h-3 w-20 rounded bg-slate-200"></div>
                              <div className="h-3 w-10 rounded-full bg-satGreen/10"></div>
                            </div>
                            <div className="space-y-1.5">
                              <div className="h-2 w-full rounded bg-slate-50"></div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-slate-200 bg-white px-4 py-12 sm:px-6 lg:px-8">
        <div className="mx-auto max-w-7xl flex flex-col md:flex-row justify-between items-center gap-6">
          <div className="flex items-center gap-2">
            <span className="brand-font text-xl font-bold">SatGo</span>
          </div>
          <p className="text-sm text-slate-500">© 2026 SatGo. Built for the future of the web.</p>
        </div>
      </footer>
    </div>
  );
}

function Feature({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-50">{icon}</div>
      <dt className="brand-font text-xl font-bold text-satBlack">{title}</dt>
      <dd className="text-base leading-7 text-slate-600">{description}</dd>
    </div>
  );
}

function Step({
  number,
  title,
  description,
}: {
  number: string;
  title: string;
  description: string;
}) {
  return (
    <div className="flex gap-4">
      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-satBlue text-sm font-bold text-white">
        {number}
      </div>
      <div>
        <h3 className="brand-font text-lg font-bold text-satBlack">{title}</h3>
        <p className="mt-1 text-slate-600">{description}</p>
      </div>
    </div>
  );
}
