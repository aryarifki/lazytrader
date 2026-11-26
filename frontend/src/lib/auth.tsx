import React, { createContext, useContext, useCallback } from "react";
import useSWR from "swr";

export interface SessionUser {
  id?: number;
  email?: string | null;
  name?: string | null;
  image?: string | null;
  role?: string;
  is_approved?: boolean;
}

export interface SessionData {
  user?: SessionUser;
  expires?: string;
}

interface AuthContextValue {
  data: SessionData | null;
  status: "loading" | "authenticated" | "unauthenticated";
  refresh: () => Promise<SessionData | undefined>;
}

const AuthContext = createContext<AuthContextValue>({
  data: null,
  status: "loading",
  refresh: async () => undefined,
});

const sessionFetcher = (url: string) =>
  fetch(url, { credentials: "include", cache: "no-store" }).then((r) => (r.ok ? r.json() : {}));

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data, isLoading, mutate } = useSWR<SessionData>("/api/auth/session", sessionFetcher, {
    revalidateOnFocus: true,
    shouldRetryOnError: false,
  });

  const hasUser = Boolean(data?.user?.email);
  const status: AuthContextValue["status"] = isLoading
    ? "loading"
    : hasUser
    ? "authenticated"
    : "unauthenticated";

  const refresh = useCallback(() => mutate(), [mutate]);

  return React.createElement(
    AuthContext.Provider,
    { value: { data: hasUser ? data ?? null : null, status, refresh } },
    children
  );
}

export function useSession() {
  return useContext(AuthContext);
}

export function signIn(provider: string = "google", opts?: { callbackUrl?: string }) {
  const cb = opts?.callbackUrl || "/";
  if (provider === "google") {
    window.location.href = `/api/auth/login?callbackUrl=${encodeURIComponent(cb)}`;
    return;
  }
  window.location.href = "/login";
}

export async function signOut(opts?: { callbackUrl?: string }) {
  try {
    await fetch("/api/auth/logout", { method: "POST", credentials: "include" });
  } catch {
    // ignore network error and redirect
  }
  window.location.href = opts?.callbackUrl || "/";
}
