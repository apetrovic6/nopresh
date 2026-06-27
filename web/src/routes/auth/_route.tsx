import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/auth/_route')({
  component: RouteComponent,
  beforeLoad: async ({ context }) => {
    if (context.auth.state.user?.isAuthenticated) {
      throw redirect({ to: "/" });
    }
  }
})

function RouteComponent() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <Outlet />
    </div>
  )
}
