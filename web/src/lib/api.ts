const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export function apiUrl(path: string): string {
  return `${API_BASE}${path}`;
}
