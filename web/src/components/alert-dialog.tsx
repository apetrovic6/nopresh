import { AlertDialog as ShadcnAlertDialog, AlertDialogHeader, AlertDialogContent, AlertDialogTitle, AlertDialogDescription, AlertDialogCancel, AlertDialogFooter, AlertDialogAction } from "./ui/alert-dialog"

interface AlertDialogProps {
  opened: boolean,
  onOpenChange: (value: boolean) => void,
  title: string,
  description: string,
  ok: string,
  cancel: string,
  onCancel: () => void,
  onConfirm: () => void,
}


export default function AlertDialog({ opened, onOpenChange, onConfirm, onCancel, title, description, ok, cancel }: AlertDialogProps) {
  return (
    <ShadcnAlertDialog open={opened} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>
            {description}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onCancel} variant="outline">{cancel}</AlertDialogCancel>
          <AlertDialogAction onClick={onConfirm} variant="destructive">{ok}</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </ShadcnAlertDialog>
  )
}
