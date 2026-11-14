import Link from "next/link";
import { ArrowLeft, ExternalLink } from "lucide-react";

export default function DemoPage() {
  return (
    <main className="min-h-screen bg-slate-50/30 text-satBlack flex flex-col">
      {/* Header / Navigation */}
      <header className="border-b border-slate-200/80 bg-white/80 backdrop-blur-md sticky top-0 z-10">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <Link
            className="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-3.5 py-2 text-sm font-semibold text-slate-700 shadow-sm transition hover:bg-slate-50 hover:text-satBlack hover:border-slate-300"
            href="/"
          >
            <ArrowLeft size={16} aria-hidden="true" />
            Back to dashboard
          </Link>
          <div className="flex items-center gap-2 rounded-full bg-blue-50 px-3 py-1 text-xs font-bold text-satBlue">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-satBlue"></span>
            </span>
            Demo Environment
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-br from-slate-950 via-[#0b0f14] to-slate-900 text-white">
        {/* Subtle decorative glow */}
        <div className="absolute -left-1/4 -top-1/4 h-[500px] w-[500px] rounded-full bg-blue-500/10 blur-[100px] pointer-events-none" />
        <div className="absolute -right-1/4 -bottom-1/4 h-[500px] w-[500px] rounded-full bg-blue-500/10 blur-[100px] pointer-events-none" />

        <div className="relative mx-auto grid w-full max-w-7xl gap-12 lg:gap-16 px-6 py-16 sm:py-24 lg:grid-cols-[1fr_450px]">
          <div className="flex flex-col justify-center">
            <p className="text-xs font-bold uppercase tracking-wider text-blue-400">
              Demo portfolio site
            </p>
            <h1 className="brand-font mt-4 max-w-2xl text-4xl font-bold leading-tight text-white sm:text-5xl lg:text-6xl">
              Sarah builds sharp websites for independent founders.
            </h1>
            <p className="mt-6 max-w-xl text-lg leading-relaxed text-slate-300">
              This page shows how SatGo fits into a normal website. The contact form below is
              embedded as an iframe and only delivers paid messages.
            </p>
          </div>

          <div className="rounded-2xl border border-white/10 bg-white/[0.03] backdrop-blur-md p-8 text-white shadow-2xl flex flex-col justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-wider text-blue-400">Services</p>
              <ul className="mt-6 space-y-4 text-sm leading-6 text-slate-300">
                <li className="flex items-center gap-2.5">
                  <span className="h-1.5 w-1.5 rounded-full bg-blue-400" />
                  Landing pages for new products
                </li>
                <li className="flex items-center gap-2.5">
                  <span className="h-1.5 w-1.5 rounded-full bg-blue-400" />
                  Responsive portfolio sites
                </li>
                <li className="flex items-center gap-2.5">
                  <span className="h-1.5 w-1.5 rounded-full bg-blue-400" />
                  Conversion-focused redesigns
                </li>
              </ul>
            </div>
            <a
              className="mt-8 inline-flex items-center justify-center gap-2 rounded-xl bg-satBlue px-5 py-3.5 text-sm font-semibold text-white transition hover:bg-blue-700 shadow-lg shadow-satBlue/25 hover:shadow-satBlue/35"
              href="#contact"
            >
              <ExternalLink size={16} aria-hidden="true" />
              Contact Sarah
            </a>
          </div>
        </div>
      </section>

      {/* Embedded Iframe Form Section */}
      <section
        className="mx-auto grid w-full max-w-7xl gap-12 lg:gap-16 px-6 py-20 lg:py-28 lg:grid-cols-[1fr_480px]"
        id="contact"
      >
        <div className="flex flex-col justify-center">
          <p className="text-xs font-bold uppercase tracking-wider text-satBlue">
            Protected contact
          </p>
          <h2 className="brand-font mt-3 text-3xl font-bold text-satBlack sm:text-4xl leading-tight">
            Spam-resistant by default
          </h2>
          <p className="mt-4 max-w-xl text-lg leading-relaxed text-slate-600">
            A visitor pays a few sats before the message reaches Sarah. Real leads keep a simple
            path, while automated bulk spam gets a direct cost.
          </p>
          <div className="mt-10 grid gap-4 sm:grid-cols-3">
            <Fact label="Widget" value="iframe" />
            <Fact label="Payment" value="WebLN or QR" />
            <Fact label="Demo price" value="5 sats" />
          </div>
        </div>

        <div className="flex justify-center lg:justify-end">
          <div className="w-full max-w-[460px] overflow-hidden rounded-2xl border border-slate-200/80 bg-white shadow-2xl shadow-slate-200 transition duration-300 hover:shadow-3xl">
            <iframe
              className="h-[620px] w-full"
              src="/embed/demo-form"
              title="SatGo contact form"
            />
          </div>
        </div>
      </section>
    </main>
  );
}

function Fact({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-xl border border-slate-200/60 bg-white p-5 transition duration-200 hover:border-slate-300 hover:shadow-sm">
      <p className="text-[11px] font-bold uppercase tracking-wider text-slate-400">{label}</p>
      <p className="brand-font mt-2 text-lg font-bold text-satBlack">{value}</p>
    </div>
  );
}
