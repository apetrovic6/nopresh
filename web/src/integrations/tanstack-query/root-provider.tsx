import { authStore } from '#/store/auth-store'
import { queryClient } from '../connect'

export function getContext() {
  return {
    queryClient,
    auth: authStore
  }
}
export default function TanstackQueryProvider() {}
