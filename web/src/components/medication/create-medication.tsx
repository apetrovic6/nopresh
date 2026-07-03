import { createMedication, getMedications, updateMedication, } from "#/gen/proto/medication/v1/medication-MedicationService_connectquery";
import { useAppForm } from "#/hooks/demo.form"
import { createQueryOptions, useMutation } from "@connectrpc/connect-query";
import { FieldGroup } from "../ui/field";
import z from "zod";
import { MEDICATIONMEAUSEREMENT, MEDICATIONMEAUSEREMENTSchema, type Medication } from "#/gen/proto/medication/v1/medication_pb";
import { queryClient, transport } from "#/integrations/connect";
import { measurementLabels } from "./medication-utils";
import { toast } from "sonner";

const validMeasurments = MEDICATIONMEAUSEREMENTSchema.values
  .filter(x => x.number !== 0)
  .map(x => x.name) as [string, ...string[]];

const schema = z.object({
  name: z.string().min(1, "Name is required"),
  recommendedDosage: z.float32().positive().min(0),
  dosageMeasurement: z.enum(validMeasurments),
});


interface CreateEditMedicationProps {
  medication?: Medication | undefined,
  formId?: string,
  onSuccess?: () => void
}

export function CreateEditMedication({ medication, formId = "create-medication", onSuccess }: CreateEditMedicationProps) {
  const { mutateAsync: createMedicationAsync } = useMutation(createMedication)
  const { mutateAsync: updateMedicationAsync } = useMutation(updateMedication)

  const defaultValues = medication ? {
    ...medication,
    dosageMeasurement: MEDICATIONMEAUSEREMENTSchema.values
      .find(x => x.number === medication.dosageMeasurement)?.name
      ?? "MEDICATIONMEAUSEREMENT_MG"
  } : {
    name: "",
    recommendedDosage: 0,
    dosageMeasurement: "MEDICATIONMEAUSEREMENT_MG"
  }

  const form = useAppForm({
    defaultValues,
    validators: {
      onSubmit: schema
    },
    onSubmit: async ({ value }) => {
      const fieldNameMap: Record<string, string> = {
        name: "name",
        recommendedDosage: "recommended_dosage",
        dosageMeasurement: "dosage_measurement",
      }


      const paths = (Object.keys(form.state.fieldMeta) as (keyof typeof value)[])
        .filter(key => form.state.fieldMeta[key]?.isDirty)
        .map(key => fieldNameMap[key])
        .filter(Boolean)

      const { name, dosageMeasurement, recommendedDosage } = value;

      const dosageMeasurmentEnum = MEDICATIONMEAUSEREMENTSchema.values
        .find(x => x.name === dosageMeasurement)?.number
        ?? MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_UNSPECIFIED;

      if (medication !== undefined) {
        await updateMedicationAsync({
          id: medication.id,
          name,
          dosageMeasurement: dosageMeasurmentEnum,
          recommendedDosage: recommendedDosage,
          updateMask: {
            paths
          }
        }, {
          onSuccess: (_) => {
            queryClient.invalidateQueries({
              queryKey: createQueryOptions(getMedications, {}, { transport }).queryKey
            });

            toast.success(name, {
              description: "Medication has been succesfully updated",
              richColors: true,
              dismissible: true,
              closeButton: true,
            })

            form.reset();
            onSuccess?.();
          },

          onError: (error) => {
            toast.error(`Error`, {
              description: `Couldn't save ${name} medication\n${error.message}`,
              richColors: true,
              dismissible: true,
              closeButton: true,
            })
          }
        });


        return;
      }

      await createMedicationAsync({
        name,
        dosageMeasurement: dosageMeasurmentEnum,
        recommendedDosage
      }, {
        onSuccess: (_) => {
          queryClient.invalidateQueries({
            queryKey: createQueryOptions(getMedications, {}, { transport }).queryKey
          });

          toast.success(name, {
            description: "Medication has been succesfully created",
            richColors: true,
            dismissible: true,
            closeButton: true,
          })

          form.reset();
          onSuccess?.();
        },

        onError: (error) => {
          console.log("error", error);

          toast.error(`Error`, {
            description: `Couldn't save ${name} medication\n${error.message}`,
            richColors: true,
            dismissible: true,
            closeButton: true,
          })
        }
      });
    },
  })


  function onMedicationSubmit(e: React.SubmitEvent) {
    e.preventDefault();
    e.stopPropagation();

    form.handleSubmit();
  }

  return (
    <form id={formId} onSubmit={onMedicationSubmit}>
      <FieldGroup>
        <form.AppField name="name">
          {field => <field.TextField type="text" label="Name" placeholder="" />}
        </form.AppField>

        <div className="grid grid-cols-2 gap-4 items-start">
          <form.AppField name="recommendedDosage" >
            {field => <field.TextField label="Recommended Dosage"
              type="number"
              step={"0.01"}
              placeholder="1"
              min={0}
            />}
          </form.AppField>

          <form.AppField name="dosageMeasurement" >
            {field => <field.DosageMeasurementPicker label="Measurement"
            />}
          </form.AppField>
        </div>
      </FieldGroup>
    </form>
  )
}
