import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { Link, linkOptions, useMatchRoute, useNavigate } from "@tanstack/react-router"
import { DropdownMenu, DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem, DropdownMenuTrigger } from "./ui/dropdown-menu";
import { useQuery, useSuspenseQuery } from "@connectrpc/connect-query";
import { logout, me } from "#/gen/proto/auth/v1/auth-AuthService_connectquery";
import { ChevronRightIcon } from "lucide-react";
import { Item, ItemActions, ItemDescription, ItemGroup, ItemTitle } from "./ui/item";
import { authStore } from "#/store/auth-store";

const options = linkOptions([
  {
    to: "/",
    label: "Home",
  },
  {
    to: "/medication",
    label: "Medication",
  },
  {
    to: "/settings",
    label: "Settings",
  }
]);


export function AppSidebar() {
  const matchRoute = useMatchRoute();
  const navigate = useNavigate()
  const { refetch: logoutQuery } = useQuery(logout, {}, { enabled: false });
  const { data: user } = useSuspenseQuery(me, {});

  async function onLogout(_: React.MouseEvent) {
    await logoutQuery();
    authStore.actions.logoutUser();
    navigate({ to: "/auth/login" });
  }

  return (
    <Sidebar>
      <SidebarHeader >
        <h1>
          NoPresh
        </h1>
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
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild>
              <DropdownMenu>
                <DropdownMenuTrigger className="w-full">
                  <User email={user.email} name={user.name} />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuGroup>
                    <DropdownMenuItem>Profile</DropdownMenuItem>
                    <DropdownMenuItem asChild>
                      <Link to="/settings"> Settings</Link>
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={onLogout}>Logout</DropdownMenuItem>
                  </DropdownMenuGroup>
                </DropdownMenuContent>
              </DropdownMenu>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}

interface UserProps {
  name: string,
  email: string,
}

function User({ email, name }: UserProps) {

  const title = name ? name[0].toUpperCase() + name.slice(1) : "";

  return (

    <Item variant={"outline"} className="h-min m-0 px-3 py-2 flex justify-between">
      <ItemGroup>
        <ItemTitle>
          {title}
        </ItemTitle>
        <ItemDescription>
          {email}
        </ItemDescription>
      </ItemGroup>
      <ItemActions>
        <ChevronRightIcon className="size-5" />
      </ItemActions>
    </Item>
  )
}
