import type { PropsWithChildren, ReactNode } from "react";
import { Dialog as ShadcnDialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "./ui/dialog";

interface DialogProps extends PropsWithChildren {
  open: boolean,
  onOpenChange: (value: boolean) => void,
  title: string,
  description?: string,
  submit: ReactNode
  cancel: ReactNode
}

export default function Dialog({ open, onOpenChange, title, description, cancel, submit, children }: DialogProps) {
  return (
    <ShadcnDialog open={open} onOpenChange={onOpenChange}>
      <form>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
            <DialogDescription>
              {description}
            </DialogDescription>
          </DialogHeader>
          {children}
          <DialogFooter>
            <DialogClose asChild>
              {cancel}
            </DialogClose>
            {submit}
          </DialogFooter>
        </DialogContent>
      </form>
    </ShadcnDialog>
  )
}
