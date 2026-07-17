import Dialog from '#/components/dialog'
import { CreateEditMedication } from '#/components/medication/forms/create-medication'
import MedicationTable from '#/components/medication/grid/data-table'
import { Button } from '#/components/ui/button'
import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_app/medication')({
  component: RouteComponent,
})


function RouteComponent() {
  const createEditFormId = "create-edit-medication"  ;
  const [openEdit, setOpenEdit] = useState(false);

  function onSuccess() {
    setOpenEdit(false);
  }

  return <div className='space-y-8'>
    <div className='flex justify-between m-2'>
      <h1 className="scroll-m-20 text-left text-4xl font-extrabold tracking-tight text-balance">
        Medications
      </h1>

    <Button variant="default" onClick={() => setOpenEdit(true)}>New Medication</Button>
    </div>

    <MedicationTable />

    <Dialog open={openEdit}
      onOpenChange={setOpenEdit}
      title="Create Medication"
      description="Make changes to your medication here. Click save when you&apos;re done."
      cancel={
        <Button variant="outline">Cancel</Button>
      }

      submit={
        <Button form={createEditFormId} type="submit">Create</Button>
      }>
      <CreateEditMedication formId={createEditFormId} onSuccess={onSuccess} />
    </Dialog>
  </div>
}
