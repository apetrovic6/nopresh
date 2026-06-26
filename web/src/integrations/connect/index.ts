import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient } from "@tanstack/react-query";

export const transport = createConnectTransport({
  baseUrl: "http://localhost:5000/api",
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
});

export const queryClient = new QueryClient();
