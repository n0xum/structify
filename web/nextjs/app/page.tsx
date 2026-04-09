import { Suspense } from "react";
import { StructifyApp } from "@/components/StructifyApp";

export default function Page() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-[var(--color-bg)]">
          <span className="text-sm text-[var(--color-text-muted)]">Loading...</span>
        </div>
      }
    >
      <StructifyApp />
    </Suspense>
  );
}
