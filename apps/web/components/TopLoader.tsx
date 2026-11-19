"use client";

import NextTopLoader from "nextjs-toploader";

export default function TopLoader() {
  return (
    <NextTopLoader
      color="#2563eb"
      height={3}
      showSpinner={false}
      shadow="0 0 10px #2563eb,0 0 5px #2563eb"
    />
  );
}
