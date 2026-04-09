import type { Metadata } from "next";
import { DocsCategoryPage } from "@/app/docs/_components/DocsCategoryPage";

export const metadata: Metadata = {
  title: "Project Guide Docs - structify",
  description: "Structify project guide for local commands, structure, and troubleshooting.",
};

export default function DocsProjectGuidePage() {
  return <DocsCategoryPage categoryId="project-guide" />;
}
