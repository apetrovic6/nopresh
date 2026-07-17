import { toast, type ExternalToast } from "sonner"

type ToastVariant = "success" | "info" | "warning" | "error"

const baseOptions: ExternalToast = {
  richColors: true,
  dismissible: true,
  closeButton: true,
}

function show(variant: ToastVariant, title: string, description?: string, options?: ExternalToast) {
  return toast[variant](title, { ...baseOptions, description, ...options })
}

// `toast` is a module-level singleton, so this object is stable — no state, no
// memo needed. Exported directly too, for use outside React components.
export const notify = {
  success: (title: string, description?: string, options?: ExternalToast) => show("success", title, description, options),
  info: (title: string, description?: string, options?: ExternalToast) => show("info", title, description, options),
  warning: (title: string, description?: string, options?: ExternalToast) => show("warning", title, description, options),
  error: (title: string, description?: string, options?: ExternalToast) => show("error", title, description, options),
}

export function useToast() {
  return notify
}
