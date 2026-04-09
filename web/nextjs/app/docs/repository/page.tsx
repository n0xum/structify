import type { Metadata } from "next";
import { DocsCategoryPage } from "@/app/docs/_components/DocsCategoryPage";

export const metadata: Metadata = {
  title: "Repository Docs - structify",
  description: "Structify repository generation documentation, conventions, and workflows.",
};

export default function DocsRepositoryPage() {
  return <DocsCategoryPage categoryId="repository" />;
}
