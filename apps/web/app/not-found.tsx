import Link from "next/link";
import { ArrowLeft, Ghost } from "lucide-react";

export default function NotFound() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-[#f8fafc] px-6 text-center text-satBlack">
      <div className="relative">
        <div className="absolute -inset-8 rounded-full bg-blue-50 opacity-50 blur-2xl" />
        <Ghost className="relative text-satBlue" size={80} />
      </div>
      <h1 className="brand-font mt-8 text-4xl font-bold sm:text-6xl">404: Lost in Orbit</h1>
      <p className="mt-4 max-w-md text-lg text-slate-600">
        We couldn&apos;t find the page you&apos;re looking for. It might have been moved or never
        existed.
      </p>
      <Link
        href="/"
        className="mt-10 inline-flex h-11 items-center justify-center gap-2 rounded-md bg-satBlack px-6 text-sm font-semibold text-white transition hover:bg-black"
      >
        <ArrowLeft size={18} />
        Back to safety
      </Link>
    </main>
  );
}
