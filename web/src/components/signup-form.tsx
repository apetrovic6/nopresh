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

import { z } from "zod";
import { useAppForm } from "#/hooks/demo.form"
import { useMutation } from "@connectrpc/connect-query";
import { register } from "#/gen/proto/auth/v1/auth-AuthService_connectquery";
import Error from "./Error";
import { Link, useNavigate } from "@tanstack/react-router";

const PASSWORD_MIN_LENGTH = 5;

const schema = z.object({
  name: z.string().trim().min(1, "Name is required"),
  email: z.email().trim(),
  password: z.string().min(PASSWORD_MIN_LENGTH, "Password must be at least 5 characters long"),
  confirmedPassword: z.string()
})
  .refine(data => data.password === data.confirmedPassword, {
  path: ["confirmedPassword"],
  error: "Passwords don't match",
})


export function SignupForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const { error, isError, mutateAsync } = useMutation(register);
  const navigate = useNavigate();

  const form = useAppForm({
    defaultValues: {
      name: "",
      email: "",
      password: "",
      confirmedPassword: "",
    },
    validators: {
      onBlur: schema,
    },
    onSubmit: async ({ value }) => {
      const { name, email, password } = value;
      await mutateAsync({ email, name, password }, {
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
          <CardTitle className="text-xl">Create your account</CardTitle>
          <CardDescription>
            Enter your email below to create your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={onRegisterSubmit}>
            <FieldGroup>
              <form.AppField name="name">
                {(field) => <field.TextField label="Name" placeholder="John Doe" />}
              </form.AppField>


              <form.AppField name="email">
                {field => <field.TextField type="email" label="Email" placeholder="john.doe@example.com" />}
              </form.AppField>


              <Field>
                <Field className="grid grid-cols-2 gap-4">
                  <form.AppField name="password">
                    {field => <field.TextField type="password" label="Password" />}
                  </form.AppField>

                  <form.AppField name="confirmedPassword" validators={
                    {

                      onChange: ({ value, fieldApi }) =>
                        fieldApi.form.getFieldValue('password') !== value
                          ? "Passwords don't match"
                          : undefined,

                      // onBlur: ({ value, fieldApi }) =>
                      //   fieldApi.form.getFieldValue('password') !== value
                      //     ? "Passwords don't match"
                      //     : undefined,

                      // onSubmit: ({ value, fieldApi }) =>
                      //   fieldApi.form.getFieldValue('password') !== value
                      //     ? "Passwords don't match"
                      //     : undefined,
                    }
                  }>
                    {field => <field.TextField type="password" label="Confirm Password" />}
                  </form.AppField>
                </Field>
                <FieldDescription>
                  Must be at least {PASSWORD_MIN_LENGTH} characters long.
                </FieldDescription>
              </Field>
              <Field>
                <Button type="submit">Create Account</Button>
                <FieldDescription className="text-center">
                  Already have an account? <Link to="/auth/login">Sign in</Link>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
