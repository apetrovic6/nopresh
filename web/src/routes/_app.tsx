import { AppSidebar } from '#/components/Sidebar';
import { SidebarProvider, SidebarTrigger } from '#/components/ui/sidebar';
import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import Footer from '#/components/Footer';
import { callUnaryMethod, createQueryOptions } from '@connectrpc/connect-query';
import { transport } from '#/integrations/connect';
import { me, logout } from '#/gen/proto/auth/v1/auth-AuthService_connectquery';
import { Toaster } from '#/components/ui/sonner';
import { getSettings } from '#/gen/proto/settings/v1/settings-SettingsService_connectquery';


const meQueryOptions = createQueryOptions(me, {}, { transport });
const settingsQuerytOptions = createQueryOptions(getSettings, {}, { transport });

export const Route = createFileRoute('/_app')({
  component: RouteComponent,


  beforeLoad: async ({ context }) => {
    if (!context.auth.state.user?.isAuthenticated) {
      try {
        const data = await callUnaryMethod(transport, me, {});
        context.auth.actions.addUser({ email: data.email, name: data.name, isAuthenticated: true });
        // Prime the cache so the loader's ensureQueryData(me) is a no-op.
        context.queryClient.setQueryData(meQueryOptions.queryKey, data);
      } catch {
        try {
          await callUnaryMethod(transport, logout, {});
        } catch { }
        throw redirect({ to: "/auth/login" });
      }
    }
  },

  loader: ({ context }) => {
    return Promise.all([
      context.queryClient.ensureQueryData(settingsQuerytOptions),
      context.queryClient.ensureQueryData(meQueryOptions),
    ]);
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
