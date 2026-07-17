import { createStore } from "@tanstack/store";


type AuthStoreUpdate = (updater: (prev: AuthStore) => AuthStore) => void;

export type AuthStore = {
  user?: AuthenticatedUser
}
export type AuthActions = {
  addUser: (user: AuthenticatedUser) => void
  logoutUser: () => void
}

export type AuthenticatedUser = {
  name: string,
  email: string,
  isAuthenticated: boolean
}


function addUserAction(user: AuthenticatedUser, setState: (updater: (prev: AuthStore) => AuthStore) => void) {
  setState((prev) => ({
    ...prev,
    user
  }))
}


function logoutUserAction(setState: AuthStoreUpdate) {
  setState(prev => ({
    ...prev, user: undefined
  }))
}

export const authStore = createStore<AuthStore, AuthActions>(
  {
    user: undefined
  },
  ({ setState }) =>
  ({
    addUser: (user: AuthenticatedUser) => addUserAction(user, setState),
    logoutUser: () => logoutUserAction(setState),
  }),
)


