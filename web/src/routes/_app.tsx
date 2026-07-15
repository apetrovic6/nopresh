import { AppSidebar } from '#/components/Sidebar';
import { SidebarProvider, SidebarTrigger } from '#/components/ui/sidebar';
import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import Footer from '#/components/Footer';
import { callUnaryMethod, createQueryOptions } from '@connectrpc/connect-query';
import { transport } from '#/integrations/connect';
import { me, logout } from '#/gen/proto/auth/v1/auth-AuthService_connectquery';
import { checkMe } from '#/lib/serverFns/fetchMe';
import { Toaster } from '#/components/ui/sonner';
import { getSettings } from '#/gen/proto/settings/v1/settings-SettingsService_connectquery';


const meQueryOptions = createQueryOptions(me, {}, { transport });
const settingsQuerytOptions = createQueryOptions(getSettings, {}, { transport });

export const Route = createFileRoute('/_app')({
  component: RouteComponent,

  beforeLoad: async ({ context }) => {
    if (import.meta.env.SSR || !context.auth.state.user?.isAuthenticated) {
      try {
        const data = import.meta.env.SSR
          ? await checkMe()
          : await callUnaryMethod(transport, me, {});

        context.auth.actions.addUser({ email: data.email, name: data.name, isAuthenticated: true });

        // Return the raw proto response so loader can prime React Query cache
        // without a second network call.
        return { _meResponse: data };
      } catch {
        if (!import.meta.env.SSR) {
          try {
            await callUnaryMethod(transport, logout, {});
          }
          catch { }
        }

        throw redirect({ to: "/auth/login" });
      }
    }
  },

  loader: ({ context }) => {
    if (import.meta.env.SSR) {
      // On SSR, beforeLoad already called checkMe() and stored the response in context.
      // Write it into the React Query cache so useSuspenseQuery(me, {}) reads from it
      // on the client after hydration (router.tsx dehydrate/hydrate transfers the cache).
      const meResponse = (context as any)._meResponse;
      if (meResponse) {
        context.queryClient.setQueryData(meQueryOptions.queryKey, meResponse);
      }
    } else {
      // On the client, fetch via the regular transport.
      // ensureQueryData is a no-op if the data is already in the cache.
      return Promise.all([
        context.queryClient.ensureQueryData(settingsQuerytOptions),
        context.queryClient.ensureQueryData(meQueryOptions)]);
    }
  },
})

function RouteComponent() {
  return (
    <SidebarProvider className="h-svh">
      <AppSidebar />
      <div className="flex flex-col flex-1 min-h-0">
        <SidebarTrigger />
        <Outlet />
        <Footer />
      </div>
      <Toaster />
    </SidebarProvider>
  )
}
