import { Suspense } from "react";
import { StructifyApp } from "@/components/StructifyApp";

export default function Page() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-zinc-950">
          <span className="text-zinc-400 text-sm">Loading…</span>
        </div>
      }
    >
      <StructifyApp />
    </Suspense>
  );
}
