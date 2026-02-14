import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "structify — Go structs to PostgreSQL, instantly",
  description:
    "Paste a Go struct and generate PostgreSQL CREATE TABLE statements or database/sql CRUD code in one click.",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="antialiased">{children}</body>
    </html>
  );
}
