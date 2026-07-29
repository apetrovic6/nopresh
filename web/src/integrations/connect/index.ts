import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient } from "@tanstack/react-query";


export const baseUrl: string = "/api";

export const transport = createConnectTransport({
  baseUrl,
  // On the client the browser attaches the httpOnly cookies automatically.
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
});

export const queryClient = new QueryClient();
