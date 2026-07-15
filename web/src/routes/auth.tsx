import Header from '#/components/Header'
import { me } from '#/gen/proto/auth/v1/auth-AuthService_connectquery';
import { transport } from '#/integrations/connect';
import { checkMe } from '#/lib/serverFns/fetchMe';
import { callUnaryMethod } from '@connectrpc/connect-query';
import { createFileRoute, isRedirect, Outlet, redirect } from '@tanstack/react-router'


export const Route = createFileRoute('/auth')({
  component: RouteComponent,
  // Check if the user is already logged in, in that case, redirect to / 
  beforeLoad: async () => {
      try {
        const data = import.meta.env.SSR
          ? await checkMe()
          : await callUnaryMethod(transport, me, {});

        if (data) {
          throw redirect({ to: "/" });
        }
      } catch (e){
        if (isRedirect(e)) throw e;
      }
    }
  },
)

function RouteComponent() {
  return (
    <div className='min-h-screen flex flex-col'>
      <Header />
      <div className='flex-1 flex items-center justify-center'>
        <Outlet />
      </div>
    </div>
  )
}
