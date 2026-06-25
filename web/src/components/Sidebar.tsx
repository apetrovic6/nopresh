import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenuButton,
} from "@/components/ui/sidebar"
import { Link, linkOptions, useMatchRoute } from "@tanstack/react-router"

const options = linkOptions([
  {
    to: "/",
    label: "Home",
  },
  {
    to: "/medication",
    label: "Medications",
  },
  {
    to: "/auth/login",
    label: "Login",
  },
  {
    to: "/auth/register",
    label: "Register",
  },
  {
    to: "/settings",
    label: "Settings",
  }
]);


export function AppSidebar() {
  const matchRoute = useMatchRoute();

  return (
    <Sidebar>
      <SidebarHeader >
        ugala bugala
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          {options.map(opt => {
            const isActive = matchRoute({ to: opt.to });
            return <SidebarMenuButton key={opt.to} asChild isActive={!!isActive}>
              <Link {...opt} activeProps={{ className: `font-bold` }} className="p-2">
                {opt.label}
              </Link>
            </SidebarMenuButton>
          })
          }
        </SidebarGroup>
        <SidebarGroup />
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  )
}
