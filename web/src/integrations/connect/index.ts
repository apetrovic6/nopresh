import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient } from "@tanstack/react-query";
import type { Interceptor } from "@connectrpc/connect";

const baseUrl = "http://localhost:5000/api";

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
