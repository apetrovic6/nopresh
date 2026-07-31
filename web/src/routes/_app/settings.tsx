import { Button } from '#/components/ui/button';
import { FieldGroup } from '#/components/ui/field';
import { Spinner } from '#/components/ui/spinner';
import { getMedications } from '#/gen/proto/medication/v1/medication-MedicationService_connectquery'
import { createSettings, getSettings, updateSettings } from '#/gen/proto/settings/v1/settings-SettingsService_connectquery';
import type { Medication } from '#/gen/proto/medication/v1/medication_entry_pb';
import { SettingsSchema, type Settings } from '#/gen/proto/settings/v1/settings_pb';
import { useAppForm } from '#/hooks/demo.form';
import { queryClient, transport } from '#/integrations/connect';
import { createQueryOptions, useMutation, useQuery } from '@connectrpc/connect-query'
import { createFileRoute } from '@tanstack/react-router'
import z from 'zod';
import { dirtyFieldMaskPaths } from '#/lib/field-mask';
import { useToast } from '#/hooks/use-toast';

const timeZones = Intl.supportedValuesOf('timeZone').map(timezone => ({ label: timezone, value: timezone }));

export interface SettingsFormProps {
  medications: Medication[],
  settings: Settings | undefined,
  formId: string,
}

export const Route = createFileRoute('/_app/settings')({
  component: RouteComponent,
})

const schema = z.object({
  defaultMedicationId: z.string().transform(val => Number(val)).pipe(z.number().int().nonnegative()),
  timezone: z.string(),
});

function RouteComponent() {
  const { data: medications, isLoading: isLoadingMedications } = useQuery(getMedications, {},);
  const { data: settings, isLoading: isLoadingSettings } = useQuery(getSettings, {}, {});

  if (isLoadingMedications || isLoadingSettings) {
    return <Spinner />
  }

  return (
    <SettingsForm
      medications={medications?.medications ?? []}
      settings={settings?.settings}
      formId='create-edit-settings-form'
    />
  )
}


function SettingsForm({ medications, settings, formId }: SettingsFormProps) {
  const toast = useToast();
  const { mutateAsync: createSettingsAsync, isPending: createSettingsPending } = useMutation(createSettings)
  const { mutateAsync: updateSettingsAsync, isPending: updateSettingsPending } = useMutation(updateSettings)

  async function createSettingsEntry(value: typeof form.state.values) {
    await createSettingsAsync({
      defaultMedicationId: Number(value.defaultMedicationId),
      timezone: value.timezone
    }, {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: createQueryOptions(getSettings, {}, { transport }).queryKey
        });

        toast.success("Settings saved", "Your settings have been succesfully updated");
      },

      onError: (error) => {
        toast.error("Error", `Couldn't save settings \n${error.message}`);
      }
    });
  }

  async function updateSettingsEntry(value: typeof form.state.values) {
    const paths = dirtyFieldMaskPaths(SettingsSchema, form.state.fieldMeta);
    
    await updateSettingsAsync({
      defaultMedicationId: Number(value.defaultMedicationId),
      timezone: value.timezone,
      updateMask: {
        paths
      }
    }, {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: createQueryOptions(getSettings, {}, { transport }).queryKey
        });

        toast.success("Settings saved", "Your settings have been succesfully updated");
      },

      onError: (error) => {
        toast.error("Error", `Couldn't save settings \n${error.message}`);
      }
    })
  }

  console.log("settings ugala bugala", settings);
  const form = useAppForm({
    defaultValues: {
      defaultMedicationId: settings?.defaultMedicationId?.toString() ?? "",
      timezone: settings?.timezone ?? "",
    },
    validators: {
      onSubmit: schema,
    },
    onSubmit: async ({ value }) => {
      if (settings === undefined) {
        await createSettingsEntry(value);
      } else {
        await updateSettingsEntry(value);
      }
    },
  });

  function onSettingsSubmit(e: React.SubmitEvent): void {
    e.preventDefault();
    e.stopPropagation();
    form.handleSubmit();
  }

  return (
    <section className='m-2 flex-1 min-h-0 overflow-auto'>
      <h1 className='scroll-m-20 text-left text-4xl font-extrabold tracking-tight text-balance'>Settings</h1>
      <div>
        <form onSubmit={onSettingsSubmit} id={formId} >
          <FieldGroup className='lg:max-w-xs my-2'>
            <form.AppField name="defaultMedicationId">
              {field => <field.MedicationPicker label="Default Medication"
                values={medications}
              />}
            </form.AppField>

            <form.AppField name="timezone">
              {field => <field.Combobox<string> label='Timezone'
                items={timeZones.map(x => x.value)}
              />}
            </form.AppField>
          </FieldGroup>
          <Button disabled={createSettingsPending || updateSettingsPending} type="submit">
            {((createSettingsPending || updateSettingsPending) &&
              <>
                <Spinner />  Saving Settings
              </>)
              ||
              "Save Settings"
            }
          </Button>
        </form>
      </div>
    </section>
  )
}
