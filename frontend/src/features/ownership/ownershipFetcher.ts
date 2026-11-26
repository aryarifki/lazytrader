export async function ownershipFetcher<T = unknown>(url: string): Promise<T> {
  const response = await fetch(url, { credentials: "include" });
  if (!response.ok) {
    let detail = `HTTP ${response.status}`;
    try {
      const body = await response.json();
      detail = body?.detail ?? detail;
    } catch {
      // keep the HTTP status
    }
    throw new Error(detail);
  }
  return response.json() as Promise<T>;
}
