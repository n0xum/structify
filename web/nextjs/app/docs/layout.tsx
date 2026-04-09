import type { ReactNode } from "react";
import { DocsSidebar } from "@/app/docs/_components/DocsSidebar";

export default function DocsLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-[var(--color-bg)] text-[var(--color-text-primary)] selection:bg-[var(--color-accent-soft)]">
      <div className="pointer-events-none fixed inset-0 overflow-hidden">
        <div className="absolute -left-20 top-0 h-[35rem] w-[35rem] rounded-full bg-[var(--color-accent-soft)] blur-[140px]" />
        <div className="absolute -bottom-20 right-0 h-[30rem] w-[30rem] rounded-full bg-[var(--color-accent-soft)] blur-[140px]" />
      </div>

      <div className="relative z-10 mx-auto flex min-h-screen w-full max-w-[1600px]">
        <DocsSidebar />
        <main className="min-w-0 flex-1 px-6 py-10 lg:px-14 lg:py-16 xl:px-20">{children}</main>
      </div>
    </div>
  );
}
