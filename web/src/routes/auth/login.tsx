import { LoginForm } from '#/components/login-form'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/auth/login')({
  component: Login,
})

function Login() {
  return <LoginForm />
}
