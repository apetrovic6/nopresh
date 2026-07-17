import { getSettings } from "#/gen/proto/settings/v1/settings-SettingsService_connectquery";
import { useAppForm } from "#/hooks/demo.form"
import { createConnectQueryKey, useMutation, useQuery, useSuspenseQuery } from "@connectrpc/connect-query";
import { FieldGroup } from "../ui/field";
import { getMedications } from "#/gen/proto/medication/v1/medication-MedicationService_connectquery";
import { createBloodPressure, getBloodPressure, updateBloodPressure } from "#/gen/proto/bloodpressure/v1/bloodpressure-BloodPressureService_connectquery";
import z from "zod";
import { queryClient, transport } from "#/integrations/connect";
import { toast } from "sonner";
import { BloodPressureEntrySchema, type BloodPressureEntry } from "#/gen/proto/bloodpressure/v1/bloodpressure_pb";
import {  MedicationUtils, } from "../medication/medication-utils";
import { timestampDate, timestampNow } from "@bufbuild/protobuf/wkt";
import { startOfMinute, } from "date-fns";
import { TZDate } from "@date-fns/tz";
import { dirtyFieldMaskPaths } from "#/lib/field-mask";



interface CreateEditBpProps {
  bpEntry?: BloodPressureEntry | undefined,
  formId?: string,
  onSuccess?: () => void
}

const schema = z.object({
  diastolic: z.number().positive(),
  systolic: z.number().positive(),
  pulse: z.number().positive(),
  dosage: z.number().nonnegative(),
  dateTimeUtc: z.date().nonoptional(),
  medication: z.string().transform(val => Number(val)).pipe(z.number().int().nonnegative()),
  medicationTaken: z.boolean(),
})

type BpFormValues = {
  diastolic: number
  systolic: number
  pulse: number
  dosage: number
  dateTimeUtc: Date
  medication: string
  medicationTaken: boolean
}

export function CreateEditBp({ bpEntry, formId, onSuccess }: CreateEditBpProps) {
  const { data: settings } = useSuspenseQuery(getSettings, {}, { staleTime: 5 * 60 * 1000 });
  const { data: medications } = useQuery(getMedications, {});
  const { mutateAsync: createBloodPressureAsync } = useMutation(createBloodPressure);
  const { mutateAsync: updateBloodPressureAsync } = useMutation(updateBloodPressure);


  const defaultFormValues: BpFormValues = bpEntry ? {
    diastolic: bpEntry.diastolic,
    systolic: bpEntry.systolic,
    pulse: bpEntry.pulse,
    dosage: bpEntry.dosage,
    dateTimeUtc: new TZDate(timestampDate(bpEntry.dateTimeUtc ?? timestampNow()), settings.settings?.timezone),
    medication: bpEntry.medicationId.toString(),
    medicationTaken: bpEntry.medicationTaken,
  } : {
    diastolic: 0,
    systolic: 0,
    pulse: 0,
    dosage: MedicationUtils.getDosage(settings.settings?.defaultMedicationId ?? 0, medications?.medications ?? []) ?? 0,
    dateTimeUtc: startOfMinute(new TZDate(new Date(), settings.settings?.timezone)),
    medication: (settings.settings?.defaultMedicationId ?? 0).toString(),
    medicationTaken: false,
  }

  async function createBloodPressureEntry(values: typeof form.state.values) {
    const { dosage, diastolic, medication, pulse, systolic, medicationTaken, dateTimeUtc } = values;
    const instant = dateTimeUtc ?? new TZDate();

    await createBloodPressureAsync({
      systolic,
      diastolic,
      pulse,
      medicationId: Number(medication),
      dosage: Number(dosage),
      medicationTaken,
      dateTimeUtc: instant.toTimestamp(),
    }, {
      onSuccess() {
        queryClient.invalidateQueries({
          // cardinality: undefined -> omitted from the key, so this matches BOTH
          // the finite and the infinite (list) getBloodPressure queries. Omitting
          // `input` matches every filter/page variant.
          queryKey: createConnectQueryKey({ schema: getBloodPressure, transport, cardinality: undefined })
        });

        toast.success("Success!", {
          description: "Blood pressure entry has been successfully created",
          richColors: true,
          dismissible: true,
          closeButton: true,
        });

        onSuccess?.();
      },

      onError: (error) => {
        toast.error(`Error`, {
          description: `Couldn't save the blood pressure entry\n${error.message}`,
          richColors: true,
          dismissible: true,
          closeButton: true,
        })
      }
    })
  }

  async function updateBloodPressureEntry(values: typeof form.state.values) {
    const { dosage, diastolic, medication, pulse, systolic, medicationTaken, dateTimeUtc } = values;
    const instant = dateTimeUtc ?? new TZDate();

    const paths = dirtyFieldMaskPaths(BloodPressureEntrySchema, form.state.fieldMeta);

    await updateBloodPressureAsync({
      id: bpEntry?.id,
      systolic,
      diastolic,
      pulse,
      medicationId: Number(medication),
      dosage: Number(dosage),
      medicationTaken,
      dateTimeUtc: instant.toTimestamp(),
      updateMask: {
        paths
      }
    }, {
      onSuccess() {
        queryClient.invalidateQueries({
          // cardinality: undefined -> omitted from the key, so this matches BOTH
          // the finite and the infinite (list) getBloodPressure queries. Omitting
          // `input` matches every filter/page variant.
          queryKey: createConnectQueryKey({ schema: getBloodPressure, transport, cardinality: undefined })
        });

        toast.success("Success!", {
          description: "Blood pressure entry has been successfully created",
          richColors: true,
          dismissible: true,
          closeButton: true,
        });

        onSuccess?.();
      },

      onError: (error) => {
        toast.error(`Error`, {
          description: `Couldn't save the blood pressure entry\n${error.message}`,
          richColors: true,
          dismissible: true,
          closeButton: true,
        })
      }
    });
  }


  const form = useAppForm({
    defaultValues: defaultFormValues,
    validators: {
      onSubmit: schema,
    },

    onSubmit: async ({ value }) => {

      if (bpEntry === undefined) {
        await createBloodPressureEntry(value);
      } else {
        await updateBloodPressureEntry(value);
      }
    },
  });

  function handleFormSubmit(e: React.SubmitEvent) {
    e.preventDefault();
    e.stopPropagation();

    form.handleSubmit();
  }


  return (
    <>
      <form id={formId} onSubmit={handleFormSubmit}>
        <FieldGroup>
          <form.Subscribe selector={(state) => state.values.dateTimeUtc}>
            {(dateTime) => <div>{dateTime?.toString()}</div>}
          </form.Subscribe>

          <form.AppField name="dateTimeUtc">
            {field => <field.DateTimePicker label="Date" />}
          </form.AppField>

          <form.AppField name="systolic">
            {field => <field.UnitInput label="Systolic" type="number" min="0" unit="mm/Hg" />}
          </form.AppField>

          <form.AppField name="diastolic">
            {field => <field.UnitInput label="Diastolic" type="number" min={0} unit="mm/Hg" />}
          </form.AppField>

          <form.AppField name="pulse">
            {field => <field.UnitInput label="Pulse" type="number" min="0" unit="bpm" />}
          </form.AppField>

          <form.AppField name="dosage">
            {(field) => (
              <form.Subscribe selector={state => state.values.medication}>
                {(medicationId) => {
                  const medUnit = MedicationUtils.getMeasurementNameByMedId(Number(medicationId), medications?.medications ?? []);
                  return <field.UnitInput label="Medication dosage" min="0" step="0.1" unit={medUnit} />;
                }}
              </form.Subscribe>
            )}
          </form.AppField>

          <form.AppField
            name="medication"
            listeners={{
              onChange: ({ value }) => {
                const med = MedicationUtils.findById(Number(value), medications?.medications ?? []);
                form.setFieldValue("dosage", med?.recommendedDosage ?? 0);
              },
            }}
          >
            {field => <field.MedicationPicker label="Medication" values={medications?.medications ?? []} />}
          </form.AppField>

          <form.AppField name="medicationTaken">
            {field => <field.Switch label="Medication Taken" />}
          </form.AppField>
        </FieldGroup>
      </form>
    </>
  )
}
