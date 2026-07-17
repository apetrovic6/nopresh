import { BloodPressureTable } from '#/components/blood-pressure/blood-pressure-table'
import { CreateEditBp } from '#/components/blood-pressure/create-edit-bp-form'
import Dialog from '#/components/dialog'
import AlertDialog from '#/components/alert-dialog'
import { Button } from '#/components/ui/button'
import { Spinner } from '#/components/ui/spinner'
import { createFileRoute } from '@tanstack/react-router'
import { Suspense, useState } from 'react'
import { deleteBloodPressure, getBloodPressure } from '#/gen/proto/bloodpressure/v1/bloodpressure-BloodPressureService_connectquery'
import { createConnectQueryKey, useMutation } from '@connectrpc/connect-query'
import { queryClient, transport } from '#/integrations/connect'
import { toast } from 'sonner'
import { timestampDate } from '@bufbuild/protobuf/wkt'
import { format } from 'date-fns'
import type { BloodPressureEntry } from '#/gen/proto/bloodpressure/v1/bloodpressure_pb'

export const Route = createFileRoute('/_app/blood-pressure')({
  component: RouteComponent,
})

function RouteComponent() {
  const [creating, setCreating] = useState(false)
  const [editing, setEditing] = useState<BloodPressureEntry | null>(null)
  const [deleting, setDeleting] = useState<BloodPressureEntry | null>(null)

  const { mutateAsync: deleteBpEntryAsync } = useMutation(deleteBloodPressure)

  async function onConfirmDelete() {
    if (!deleting) return
    const entry = deleting
    setDeleting(null)

    await deleteBpEntryAsync({ id: entry.id }, {
      onSuccess: () => {
        queryClient.invalidateQueries({
          // cardinality: undefined matches both finite and infinite (list) queries.
          queryKey: createConnectQueryKey({ schema: getBloodPressure, transport, cardinality: undefined }),
        })

        const label = entry.dateTimeUtc
          ? format(timestampDate(entry.dateTimeUtc), 'yyyy.MM.dd hh:mm')
          : entry.id.toString()

        toast.success(`Blood Pressure Entry: ${label}`, {
          description: 'Entry has been deleted',
          richColors: true,
          dismissible: true,
          closeButton: true,
        })
      },
      onError: (error) => {
        toast.error('Error', {
          description: `Couldn't delete blood pressure entry \n${error.message}`,
          richColors: true,
          dismissible: true,
          closeButton: true,
        })
      },
    })
  }

  return (
    <>
      <div className='flex justify-between m-2'>
        <h1 className="scroll-m-20 text-left text-4xl font-extrabold tracking-tight text-balance">
          Blood Pressure
        </h1>
        <Button  onClick={() => setCreating(true)}>New Entry</Button>
      </div>

      <BloodPressureTable onEdit={setEditing} onDelete={setDeleting} />

      { /* The dialogs are here because it glitches out if it's in the table component */ }
      
      <Dialog
        open={creating}
        onOpenChange={setCreating}
        title="Create Blood Pressure Entry"
        submit={<Button form="create-bp" type="submit">Create</Button>}
        cancel={<Button variant="outline">Cancel</Button>}
      >
        <Suspense fallback={<Spinner />}>
          <CreateEditBp formId="create-bp" onSuccess={() => setCreating(false)} />
        </Suspense>
      </Dialog>

      <Dialog
        open={editing !== null}
        onOpenChange={(open) => { if (!open) setEditing(null) }}
        title="Edit Blood Pressure Entry"
        submit={<Button form="edit-bp" type="submit">Save Changes</Button>}
        cancel={<Button variant="outline">Cancel</Button>}
      >
        <Suspense fallback={<Spinner />}>
          {editing && <CreateEditBp formId="edit-bp" bpEntry={editing} onSuccess={() => setEditing(null)} />}
        </Suspense>
      </Dialog>

      <AlertDialog
        opened={deleting !== null}
        onOpenChange={(open) => { if (!open) setDeleting(null) }}
        title="Are you sure?"
        description="This action cannot be undone. This will permanently delete the data."
        ok="Delete"
        cancel="Cancel"
        onCancel={() => setDeleting(null)}
        onConfirm={onConfirmDelete}
      />
    </>
  )
}
