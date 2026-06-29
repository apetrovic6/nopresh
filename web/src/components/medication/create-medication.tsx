import { createMedication } from "#/gen/proto/medication/v1/medication-MedicationService_connectquery";
import { MEDICATIONMEAUSERMENT, MEDICATIONMEAUSERMENTSchema } from "#/gen/proto/medication/v1/medication_pb"
import { useAppForm } from "#/hooks/demo.form"
import { useMutation } from "@connectrpc/connect-query";
import { Button } from "../ui/button";
import { Field, FieldGroup } from "../ui/field";
import z from "zod";

const measurmentLabels: Partial<Record<string, string>> = {
  MEDICATIONMEAUSERMENT_MG: "mg",
  MEDICATIONMEAUSERMENT_G: "g",
}

const validMeasurments = MEDICATIONMEAUSERMENTSchema.values
  .filter(x => x.number !== 0)
  .map(x => x.name) as [string, ...string[]];

const schema = z.object({
  name: z.string().min(1, "Name is required"),
  recommendedDosage: z.number().positive(),
  dosageMeasurment: z.enum(validMeasurments),
});


export function CreateMedication() {
  const { mutateAsync: createMedicationAsync } = useMutation(createMedication)

  const form = useAppForm({
    defaultValues: {
      name: "",
      recommendedDosage: 0,
      dosageMeasurment: "MEDICATIONMEAUSERMENT_MG"
    },
    validators: {
      onSubmit: schema
    },
    onSubmit: async ({ value }) => {
      const { name, dosageMeasurment, recommendedDosage } = value;

      const dosageMeasurmentEnum = MEDICATIONMEAUSERMENTSchema.values
        .find(x => x.name === dosageMeasurment)?.number
        ?? MEDICATIONMEAUSERMENT.MEDICATIONMEAUSERMENT_UNSPECIFIED;

      console.log(value);
      console.log(dosageMeasurmentEnum);

      await createMedicationAsync({
        name,
        dosageMeasurment: dosageMeasurmentEnum,
        recommendedDosage
      }, {
        onSuccess: (res) => {
          console.log("res", res);
        },

        onError: (error) => {
          console.log("error", error);
        }
      });
    },
  })


  function onMedicationSubmit(e: React.SubmitEvent) {
    e.preventDefault();
    e.stopPropagation();

    form.handleSubmit();
  }

  return <div>
    Create
    <form onSubmit={onMedicationSubmit}>
      <FieldGroup>
        <form.AppField name="name">
          {field => <field.TextField type="text" label="Name" placeholder="" />}
        </form.AppField>

        <div className="grid grid-cols-2 gap-4 items-end">
          <form.AppField name="recommendedDosage" >
            {field => <field.TextField label="Recommended Dosage"
              type="number"
              placeholder="1"
            />}
          </form.AppField>

          <div className="w-fit">
            <form.AppField name="dosageMeasurment" >
              {field => <field.Select label="Meauserment"
                values={MEDICATIONMEAUSERMENTSchema.values
                  .filter(x => x.number !== 0)
                  .map(x => ({ label: measurmentLabels[x.name] ?? x.name, value: x.name }))} />}
            </form.AppField>
          </div>
        </div>
      </FieldGroup>
      <Field className="my-2">
        <Button type="submit">Create</Button>
      </Field>
    </form>
  </div>
}
