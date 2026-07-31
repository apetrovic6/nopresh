import type { BloodPressureEntry } from "#/gen/proto/bloodpressure/v1/bloodpressure_pb"
import { MedicationUtils} from "../medication/medication-utils"
import { Card, CardContent, CardHeader, } from "../ui/card"
import { Checkbox } from "../ui/checkbox"
import { Field, FieldContent, FieldGroup, FieldLabel } from "../ui/field"
import { Label } from "../ui/label"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu"
import { Button } from "../ui/button"
import { Clock, Edit, EllipsisVertical, Trash2 } from "lucide-react"
import { timestampDate } from "@bufbuild/protobuf/wkt"
import { format } from "date-fns"
import { TZDate } from "@date-fns/tz"
import { useSuspenseQuery } from "@connectrpc/connect-query"
import { getSettings } from "#/gen/proto/settings/v1/settings-SettingsService_connectquery"

interface BloodPressureItemProps {
  item: BloodPressureEntry
  onEdit: (item: BloodPressureEntry) => void
  onDelete: (item: BloodPressureEntry) => void
}

export function BloodPressureItem({ item, onEdit, onDelete }: BloodPressureItemProps) {
  const { data: settings } = useSuspenseQuery(getSettings, {}, { staleTime: 5 * 60 * 1000 });
  const dateTime = item.dateTimeUtc ? timestampDate(item.dateTimeUtc) : undefined;

  let dateText = "";

  if (dateTime !== undefined) {
    const localDateTime = new TZDate(dateTime, settings.settings?.timezone);
    dateText = format(localDateTime, "HH:mm");
  } else {
    dateText = "Invalid Time";
  }

  return (
    <Card>
      <CardHeader className="flex justify-between">
        <span>
          {dateTime && (
            <time
              dateTime={dateTime.toISOString()}
              className="flex items-center gap-1.5 text-sm font-medium text-muted-foreground tabular-nums"
            >
              <Clock className="size-4.5" />
              <span className="text-lg ">
                {dateText}
              </span>
            </time>
          )}
        </span>
        <span>
          <BloodPressureItemMenu item={item} onEdit={onEdit} onDelete={onDelete} />
        </span>
      </CardHeader>
      <CardContent>
        <div className="flex">
          <FieldGroup>
            <Field>
              <FieldLabel>Systolic</FieldLabel>
              <FieldContent>
                {item.systolic} mm/Hg
              </FieldContent>
            </Field>

            <Field>
              <FieldLabel>Diastolic</FieldLabel>
              <FieldContent>
                {item.diastolic} mm/Hg
              </FieldContent>
            </Field>


          </FieldGroup>


          <FieldGroup>
            <Field>
              <FieldLabel>Pulse</FieldLabel>
              <FieldContent>
                {item.pulse} bpm
              </FieldContent>
            </Field>

            <Field className="w-fit">
              <FieldLabel>Dosage</FieldLabel>
              <FieldContent>
                {item.dosage}  {MedicationUtils.getMeasurementByEnum(item.medication?.dosageMeasurement ?? 0)}
              </FieldContent>
            </Field>

          </FieldGroup>

          <FieldGroup className="flex flex-col justify-between">
            <Field>
              <FieldLabel>Comment</FieldLabel>
              <FieldContent>
                {item.comment}
              </FieldContent>
            </Field>

            <Field orientation="horizontal">
              <Checkbox id="medication-taken-checkbox" checked={item.medicationTaken} />
              <Label htmlFor="medication-taken-checkbox">Medication Taken ({item.medication?.name})</Label>
            </Field>

          </FieldGroup>

        </div>
      </CardContent>
    </Card>
  )
}

function BloodPressureItemMenu({ item, onEdit, onDelete }: BloodPressureItemProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost">
          <EllipsisVertical />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem onSelect={() => onEdit(item)}>
          <Edit />
          Edit
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => onDelete(item)}>
          <Trash2 />
          Delete
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
