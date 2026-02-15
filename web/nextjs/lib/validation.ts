export function validateInput(source: string): string[] {
  const warnings: string[] = [];

  if (!source.trim()) {
    warnings.push("Input is empty.");
    return warnings;
  }

  if (source.length > 100 * 1024) {
    warnings.push("Input exceeds 100 KB limit.");
  }

  // Check for at least one exported struct declaration
  if (!/\b[A-Z][a-zA-Z0-9]*\s+struct\s*\{/.test(source)) {
    warnings.push("No exported struct found. Struct names must start with an uppercase letter.");
  }

  return warnings;
}

export function validatePackageName(name: string): string | null {
  if (!name) return "Package name is required.";
  if (!/^[a-z][a-z0-9_]*$/.test(name)) {
    return "Package name must contain only lowercase letters, digits and underscores.";
  }
  return null;
}
