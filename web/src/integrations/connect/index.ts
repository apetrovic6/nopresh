import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient } from "@tanstack/react-query";
import type { Interceptor } from "@connectrpc/connect";
import { env } from "#/env";

// Resolve the API base URL at *runtime* so a single built image works in any
// environment. Set `API_URL` at runtime (and optionally `SSR_API_URL` when the
// server needs a different, container-internal address than the browser). Falls
// back to the build-time `VITE_BASE_URL` (used in dev).
function resolveBaseUrl(): string {
  if (import.meta.env.SSR) {
    return (
      process.env.SSR_API_URL ?? process.env.API_URL ?? env.VITE_BASE_URL
    );
  }
  // Injected by the server into `window.__API_URL__` (see routes/__root.tsx).
  return (globalThis as { __API_URL__?: string }).__API_URL__ || env.VITE_BASE_URL;
}

export const baseUrl: string = resolveBaseUrl();

// Interceptor that sets the cookies with jwt access and refresh tokens during SSR
const forwardCookies: Interceptor = (next) => async (req) => {
  const { getCookie } = await import("@tanstack/react-start/server");
  const jwt = getCookie("jwt") ?? "";
  const refresh = getCookie("refresh") ?? "";
  if (jwt || refresh) {
    req.header.set("Cookie", `jwt=${jwt}; refresh=${refresh}`);
  }
  return next(req);
};

export const transport = createConnectTransport({
  baseUrl,
  interceptors: import.meta.env.SSR ? [forwardCookies] : [],
  // On the client the browser attaches the httpOnly cookies automatically.
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
});

export const queryClient = new QueryClient();
