import { cn } from "#/lib/utils.ts"
import { Button } from "#/components/ui/button.tsx"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "#/components/ui/card.tsx"
import {
  Field,
  FieldDescription,
  FieldGroup,
} from "#/components/ui/field.tsx"
import { login } from "#/gen/proto/auth/v1/auth-AuthService_connectquery"
import { useMutation } from "@connectrpc/connect-query"
import { Link, useNavigate } from "@tanstack/react-router"
import { useAppForm } from "#/hooks/demo.form"
import Error from "./Error"

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const { error, isError, mutateAsync } = useMutation(login);
  const navigate = useNavigate();

  const form = useAppForm({
    defaultValues: {
      email: "",
      password: "",
    },
    onSubmit: async ({ value }) => {
      const { email, password } = value;
      await mutateAsync({ email, password }, {
        onError: (r) => console.log("error", r),
        onSuccess: (_) => {
          navigate({ to: "/" })
        }
      });
    }
  })


  function onRegisterSubmit(e: React.SubmitEvent): void {
    e.preventDefault();
    e.stopPropagation();
    form.handleSubmit();
  }

  return (
    <div className={cn("flex flex-col gap-6 w-full max-w-md", className)} {...props}>
      <Card>
        <CardHeader className="text-center">
          {isError && <Error error={`${error.rawMessage}`} />}
          <CardTitle className="text-xl">Login</CardTitle>
          <CardDescription>
            Enter your email below to login your account
          </CardDescription>
        </CardHeader>
        <CardContent className="">
          <form onSubmit={onRegisterSubmit}>
            <FieldGroup className="">
              <form.AppField name="email">
                {field => <field.TextField type="email" label="Email" placeholder="john.doe@example.com" />}
              </form.AppField>

              <Field>
                  <form.AppField name="password">
                    {field => <field.TextField type="password" label="Password" />}
                  </form.AppField>
              </Field>
              <Field>
                <Button type="submit">Login</Button>
                <FieldDescription className="text-center">
                  Don't have an account? <Link to="/auth/register">Register</Link>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
