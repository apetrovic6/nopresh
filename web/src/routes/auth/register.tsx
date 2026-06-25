import { SignupForm } from '#/components/signup-form'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/auth/register')({
  component: Register,
})

function Register() {
  return (
      <SignupForm />
  )
}
