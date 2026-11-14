import type { ButtonHTMLAttributes, ReactNode } from "react";

type PrimaryButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  children: ReactNode;
  tone?: "primary" | "success" | "neutral";
};

export function PrimaryButton({
  children,
  className = "",
  tone = "primary",
  ...props
}: PrimaryButtonProps) {
  const toneClass =
    tone === "success"
      ? "bg-satGreen text-white hover:bg-green-700"
      : tone === "neutral"
        ? "bg-satBlack text-white hover:bg-black"
        : "bg-satBlue text-white hover:bg-blue-700";

  return (
    <button
      className={`inline-flex h-11 items-center justify-center gap-2 rounded-md px-4 text-sm font-semibold transition disabled:opacity-60 ${toneClass} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
