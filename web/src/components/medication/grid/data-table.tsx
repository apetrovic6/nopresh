import { useState } from "react";
import { createColumnHelper, getCoreRowModel, flexRender, useReactTable, getPaginationRowModel, type PaginationState, type ColumnFiltersState, getFilteredRowModel } from "@tanstack/react-table";
import { fuzzyFilter } from "#/lib/table-filters";
import { measurementLabels } from "../medication-utils";
import { createQueryOptions, useMutation, useQuery } from "@connectrpc/connect-query";
import { deleteMedication, getMedications } from "#/gen/proto/medication/v1/medication-MedicationService_connectquery";
import { Spinner } from "#/components/ui/spinner";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "#/components/ui/table";
import { Button } from "#/components/ui/button";
import { ChevronsLeftIcon, ChevronsRightIcon, Edit, Trash2 } from "lucide-react";
import { ActionMenu } from "#/components/action-menu";
import type { Medication } from "#/gen/proto/medication/v1/medication_entry_pb";
import { CreateEditMedication } from "../forms/create-medication";
import { queryClient, transport } from "#/integrations/connect";
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "#/components/ui/pagination";
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "#/components/ui/select";
import { Input } from "#/components/ui/input";
import { PaginationPages } from "#/components/pagination";
import { useToast } from "#/hooks/use-toast";

const helper = createColumnHelper<Medication>()

const columns = [
  helper.accessor('name', {
    header: "Name",
    cell: info => info.getValue(),
    footer: props => props.column.id,
  }),

  helper.accessor('recommendedDosage', {
    cell: info => info.getValue(),
    header: "Recommended Dosage",
    footer: props => props.column.id,
  }),

  helper.accessor('dosageMeasurement', {
    cell: info => measurementLabels[info.getValue()] ?? info.getValue(),
    header: "Dosage Measurement",
    footer: props => props.column.id,
  }),

  helper.display({
    id: "actions",
    cell: props => <MedicationTableActionMenu row={props.row.original} />
  }),
]

const fallbackData: Medication[] = [];


function MedicationTableActionMenu({ row }: { row: Medication }) {
  const toast = useToast();
  const { mutateAsync: deleteMedicationAsync } = useMutation(deleteMedication);

  async function onMedicationDelete() {
    await deleteMedicationAsync({
      id: row.id
    }, {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: createQueryOptions(getMedications, {}, { transport }).queryKey
        });

        toast.success(`Medication: ${row.name}`, "Medication has been deleted");
      },

      onError: (error) => {
        toast.error(`Error`, `Couldn't delete ${row.name} medication\n${error.message}`)
      }
    })
  }

  return (
    <ActionMenu>
      <ActionMenu.Menu>
        <ActionMenu.Item action="edit" icon={<Edit />}>Edit</ActionMenu.Item>
        <ActionMenu.Item action="delete" icon={<Trash2 />}>Delete</ActionMenu.Item>
      </ActionMenu.Menu>

      <ActionMenu.Dialog
        action="edit"
        title="Edit Medication"
        description="Make changes to your medication here. Click save when you're done."
        cancel={<Button variant="outline">Cancel</Button>}
        submit={<Button form="create-edit-medication" type="submit">Save changes</Button>}
      >
        {(close) => (
          <CreateEditMedication medication={row} formId="create-edit-medication" onSuccess={close} />
        )}
      </ActionMenu.Dialog>

      <ActionMenu.Alert
        action="delete"
        title="Are you sure?"
        description="This action cannot be undone. This will permanently delete the data."
        ok="Delete"
        cancel="Cancel"
        onConfirm={onMedicationDelete}
      />
    </ActionMenu>
  )
}


const pageSizes = [5, 10, 15, 20];

export default function MedicationTable() {
  const { data, isFetching } = useQuery(getMedications, {});
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  });
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])

  const table = useReactTable({
    columns,
    data: data?.medications ?? fallbackData,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onPaginationChange: setPagination,
    onColumnFiltersChange: setColumnFilters,
    getFilteredRowModel: getFilteredRowModel(),
    state: { pagination, columnFilters },
    filterFns: { fuzzy: fuzzyFilter },
  })

  if (isFetching) {
    return <div className="w-full flex justify-center">
      <Spinner className="size-8" />
    </div>
  }

  const currentPage = table.getState().pagination.pageIndex + 1

  return (
    <div className="overflow-hidden rounded-md border">
      <div className="flex items-center py-4">
        <Input
          className="max-w-sm mx-2"
          placeholder="Filter name..."
          value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
          onChange={(event) =>
            table.getColumn("name")?.setFilterValue(event.target.value)
          }
        />
      </div>
      <Table className="">
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>


      <div className="flex">
        <div />
        <Pagination>

          <PaginationContent>
            <PaginationItem>
              <PaginationLink
                aria-label="Go to first page"
                onClick={() => table.setPageIndex(0)}
                aria-disabled={!table.getCanPreviousPage()}
                className={!table.getCanPreviousPage() ? "pointer-events-none opacity-50" : ""}
              >
                <ChevronsLeftIcon />
              </PaginationLink>
            </PaginationItem>
            <PaginationItem>
              <PaginationPrevious
                onClick={() => table.previousPage()}
                aria-disabled={!table.getCanPreviousPage()}
                className={!table.getCanPreviousPage() ? "pointer-events-none opacity-50" : ""}
              />
            </PaginationItem>

            <PaginationPages currentPage={currentPage}
              pageCount={table.getPageCount()}
              setPage={(page) => table.setPageIndex(page)} />

            <PaginationItem>
              <PaginationNext
                onClick={() => table.nextPage()}
                aria-disabled={!table.getCanNextPage()}
                className={!table.getCanNextPage() ? "pointer-events-none opacity-50" : ""}
              />
            </PaginationItem>
            <PaginationItem>
              <PaginationLink
                aria-label="Go to last page"
                onClick={() => table.setPageIndex(table.getPageCount() - 1)}
                aria-disabled={!table.getCanNextPage()}
                className={!table.getCanNextPage() ? "pointer-events-none opacity-50" : ""}
              >
                <ChevronsRightIcon />
              </PaginationLink>
            </PaginationItem>
          </PaginationContent>
        </Pagination>

        <Select
          value={table.getState().pagination.pageSize.toString()}
          onValueChange={(value) => table.setPageSize(Number(value))}>
          <SelectTrigger className="w-fit max-w-48 m-2">
            <SelectValue placeholder="Pick page size" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              {pageSizes.map(pageSize =>
                <SelectItem key={pageSize} value={pageSize.toString()}>{pageSize}</SelectItem>)
              }
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
