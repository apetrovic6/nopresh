import { createMedication, getMedications, updateMedication, } from "#/gen/proto/medication/v1/medication-MedicationService_connectquery";
import { useAppForm } from "#/hooks/demo.form"
import { createQueryOptions, useMutation } from "@connectrpc/connect-query";
import { FieldGroup } from "../../ui/field";
import z from "zod";
import { queryClient, transport } from "#/integrations/connect";
import { dirtyFieldMaskPaths } from "#/lib/field-mask";
import { MEDICATIONMEAUSEREMENTSchema, MedicationSchema, type Medication } from "#/gen/proto/medication/v1/medication_entry_pb";
import { useToast } from "#/hooks/use-toast";
import { MedicationUtils } from "../medication-utils";

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
  const toast = useToast();
  const { mutateAsync: createMedicationAsync } = useMutation(createMedication)
  const { mutateAsync: updateMedicationAsync } = useMutation(updateMedication)

  const defaultValues = medication ? {
    ...medication,
    dosageMeasurement: MedicationUtils.getNameByEnum(medication.dosageMeasurement) ?? "MEDICATIONMEAUSEREMENT_MG",
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
      const paths = dirtyFieldMaskPaths(MedicationSchema, form.state.fieldMeta)

      const { name, dosageMeasurement, recommendedDosage } = value;

      const dosageMeasurementEnum = MedicationUtils.getDosageMeasurementIdByName(dosageMeasurement);

      if (medication !== undefined) {
        await updateMedicationAsync({
          id: medication.id,
          name,
          dosageMeasurement: dosageMeasurementEnum,
          recommendedDosage: recommendedDosage,
          updateMask: {
            paths
          }
        }, {
          onSuccess: (_) => {
            queryClient.invalidateQueries({
              queryKey: createQueryOptions(getMedications, {}, { transport }).queryKey
            });

            toast.success(name, "Medication has been succesfully updated")

            form.reset();
            onSuccess?.();
          },

          onError: (error) => {
            toast.error(`Error`, `Couldn't save ${name} medication\n${error.message}`)
          }
        });


        return;
      }

      await createMedicationAsync({
        name,
        dosageMeasurement: dosageMeasurementEnum,
        recommendedDosage
      }, {
        onSuccess: (_) => {
          queryClient.invalidateQueries({
            queryKey: createQueryOptions(getMedications, {}, { transport }).queryKey
          });

          toast.success(name, "Medication has been succesfully created");

          form.reset();
          onSuccess?.();
        },

        onError: (error) => {
          console.log("error", error);
          toast.error(`Error`, `Couldn't save ${name} medication\n${error.message}`);
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
