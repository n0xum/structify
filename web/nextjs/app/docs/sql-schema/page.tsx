import type { Metadata } from "next";
import { DocsCategoryPage } from "@/app/docs/_components/DocsCategoryPage";

export const metadata: Metadata = {
  title: "SQL Schema Docs - structify",
  description: "Structify SQL schema documentation with tags, constraints, indexes, and foreign keys.",
};

export default function DocsSqlSchemaPage() {
  return <DocsCategoryPage categoryId="sql-schema" />;
}
