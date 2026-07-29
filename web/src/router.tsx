import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'

// Side-effect import: installs Date.prototype.toTimestamp. Must run on both
import './lib/extensions/date'

import { getContext } from './integrations/tanstack-query/root-provider'
import { QueryClientProvider } from '@tanstack/react-query'

export function getRouter() {
  const context = getContext()

  const router = createTanStackRouter({
    routeTree,
    context,
    scrollRestoration: true,
    defaultPreload: 'intent',
    defaultPreloadStaleTime: 0,

    Wrap: ({ children }) => {
      return (
        <QueryClientProvider client={context.queryClient}>
          {children}
        </QueryClientProvider>
      )
    }
  })


  return router
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
