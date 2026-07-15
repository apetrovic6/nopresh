import { createContext, useContext, useMemo, useState, type ComponentProps, type PropsWithChildren, type ReactNode } from "react";
import { EllipsisVertical } from "lucide-react";
import { Button } from "./ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "./ui/dropdown-menu";
import Dialog from "./dialog";
import AlertDialog from "./alert-dialog";

type ActionMenuContextValue = {
  active: string | null;
  open: (id: string) => void;
  close: () => void;
};

const ActionMenuContext = createContext<ActionMenuContextValue | null>(null);

function useActionMenu() {
  const ctx = useContext(ActionMenuContext);
  if (!ctx) {
    throw new Error("ActionMenu.* components must be used inside <ActionMenu>");
  }
  return ctx;
}

/**
 * Coordinates a dropdown of actions with any number of dialogs/alerts.
 *
 * A single "active action" id decides which overlay is open, so adding an
 * overlay never adds state. An <ActionMenu.Item action="x"> opens the overlay
 * registered as <ActionMenu.Dialog action="x"> / <ActionMenu.Alert action="x">.
 */
export function ActionMenu({ children }: PropsWithChildren) {
  const [active, setActive] = useState<string | null>(null);

  const value = useMemo<ActionMenuContextValue>(
    () => ({ active, open: setActive, close: () => setActive(null) }),
    [active],
  );

  return <ActionMenuContext.Provider value={value}>{children}</ActionMenuContext.Provider>;
}

function Menu({ children }: PropsWithChildren) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost">
          <EllipsisVertical />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">{children}</DropdownMenuContent>
    </DropdownMenu>
  );
}

interface ItemProps {
  /** Opens the overlay registered under this id. Omit for a plain action. */
  action?: string;
  /** Runs an inline handler (with or without an overlay). */
  onSelect?: () => void;
  icon?: ReactNode;
  children: ReactNode;
}

function Item({ action, onSelect, icon, children }: ItemProps) {
  const { open } = useActionMenu();
  return (
    <DropdownMenuItem
      onSelect={() => {
        if (action) open(action);
        onSelect?.();
      }}
    >
      {icon}
      {children}
    </DropdownMenuItem>
  );
}

type ActionDialogProps = { action: string; children: (close: () => void) => ReactNode } & Omit<
  ComponentProps<typeof Dialog>,
  "open" | "onOpenChange" | "children"
>;

function ActionDialog({ action, children, ...props }: ActionDialogProps) {
  const { active, close } = useActionMenu();
  return (
    <Dialog
      open={active === action}
      onOpenChange={(o) => {
        if (!o) close();
      }}
      {...props}
    >
      {children(close)}
    </Dialog>
  );
}

type ActionAlertProps = { action: string; onConfirm: () => void } & Omit<
  ComponentProps<typeof AlertDialog>,
  "opened" | "onOpenChange" | "onCancel" | "onConfirm"
>;

function ActionAlert({ action, onConfirm, ...props }: ActionAlertProps) {
  const { active, close } = useActionMenu();
  return (
    <AlertDialog
      opened={active === action}
      onOpenChange={(o) => {
        if (!o) close();
      }}
      onCancel={close}
      onConfirm={() => {
        onConfirm();
        close();
      }}
      {...props}
    />
  );
}

ActionMenu.Menu = Menu;
ActionMenu.Item = Item;
ActionMenu.Dialog = ActionDialog;
ActionMenu.Alert = ActionAlert;
