import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'

// Side-effect import: installs Date.prototype.toTimestamp. Must run on both
// SSR and client, which this shared entry does.
import './lib/extensions/date'

import { setupRouterSsrQueryIntegration } from '@tanstack/react-router-ssr-query'
import { getContext } from './integrations/tanstack-query/root-provider'
import { dehydrate, hydrate, QueryClientProvider } from '@tanstack/react-query'

export function getRouter() {
  const context = getContext()

  const router = createTanStackRouter({
    routeTree,
    context,
    scrollRestoration: true,
    defaultPreload: 'intent',
    defaultPreloadStaleTime: 0,

    dehydrate: () => {
      return {
        queryClientState: dehydrate(context.queryClient),
      }
    },

    hydrate: (dehydrated) => {
      hydrate(context.queryClient, dehydrated.queryClientState);
    },
    Wrap: ({ children }) => {
      return (
        <QueryClientProvider client={context.queryClient}>
          {children}
        </QueryClientProvider>
      )
    }
  })

  setupRouterSsrQueryIntegration({ router, queryClient: context.queryClient })

  return router
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
