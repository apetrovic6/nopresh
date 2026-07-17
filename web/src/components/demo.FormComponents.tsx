import { useFieldContext, useFormContext } from '#/hooks/demo.form-context'

import { Button } from '#/components/ui/button'
import { Input } from '#/components/ui/input'
import { Textarea as ShadcnTextarea } from '#/components/ui/textarea'
import * as ShadcnSelect from '#/components/ui/select'
import { Slider as ShadcnSlider } from '#/components/ui/slider'
import { Switch as ShadcnSwitch } from '#/components/ui/switch'
import { Label } from '#/components/ui/label'
import { Field, FieldGroup, FieldLabel } from './ui/field'
import { useSelector } from '@tanstack/react-store'
import { InputGroup, InputGroupAddon, InputGroupInput } from './ui/input-group'
import { MEDICATIONMEAUSEREMENTSchema, type Medication } from '#/gen/proto/medication/v1/medication_entry_pb'
import { measurementLabels } from './medication/medication-utils'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'
import { Calendar } from './ui/calendar'
import { CalendarIcon, ChevronDownIcon, XIcon } from 'lucide-react'
import { useState, type ChangeEvent } from 'react'
import { format, parse, set } from "date-fns"
import type { DateRange } from 'react-day-picker'
import { ComboboxContent, ComboboxEmpty, ComboboxInput, ComboboxItem, ComboboxList, Combobox as ShadcnCombobox } from './ui/combobox'
import type { Group } from '@base-ui/react/internals/resolveValueLabel'
import type { ComboboxRoot } from '@base-ui/react'
import { TZDate } from '@date-fns/tz'

export function SubscribeButton({ label }: { label: string }) {
  const form = useFormContext()
  return (
    <form.Subscribe selector={(state) => state.isSubmitting}>
      {(isSubmitting) => (
        <Button type="submit" disabled={isSubmitting}>
          {label}
        </Button>
      )}
    </form.Subscribe>
  )
}

function ErrorMessages({
  errors,
}: {
  errors: Array<string | { message: string }>
}) {
  return (
    <>
      {errors.map((error) => (
        <div
          key={typeof error === 'string' ? error : error.message}
          className="mt-1 text-sm font-semibold text-red-600"
        >
          {typeof error === 'string' ? error : error.message}
        </div>
      ))}
    </>
  )
}



interface TextFieldProps extends React.ComponentProps<'input'> {
  label: string
}

export function TextField({
  label,
  type,
  ...props
}: TextFieldProps
) {
  const field = useFieldContext<string | number>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <Field>
      <FieldLabel
        htmlFor={label}
        className="text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </FieldLabel>
      <Input
        value={field.state.value}
        type={type}
        onBlur={field.handleBlur}
        onChange={(e) => field.handleChange(
          type === "number" ? e.target.valueAsNumber : e.target.value
        )}
        {...props}
      />
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </Field>
  )
}

interface UnitInputProps extends React.ComponentProps<'input'> {
  label: string,
  unit?: string,
}

export function UnitInput({
  label,
  unit,
  ...props
}: UnitInputProps) {
  const field = useFieldContext<string | number>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <Field>
      <FieldLabel
        htmlFor={label}
        className="text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </FieldLabel>
      <InputGroup>
        <InputGroupInput
          type="number"
          {...props}
          value={field.state.value}
          onBlur={field.handleBlur}
          onChange={(e) => field.handleChange(e.target.valueAsNumber)}
        />
        <InputGroupAddon align="inline-end">
          {unit}
        </InputGroupAddon>
      </InputGroup>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </Field>
  )
}

export function PasswordField({
  label = "Password",
  placeholder,
}: {
  label: string
  placeholder?: string
}) {
  const field = useFieldContext<string>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <div>
      <Label
        htmlFor={label}
        className="mb-2 text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </Label>
      <Input
        value={field.state.value}
        placeholder={placeholder}
        onBlur={field.handleBlur}
        onChange={(e) => field.handleChange(e.target.value)}
      />
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  )
}

export function TextArea({
  label,
  rows = 3,
}: {
  label: string
  rows?: number
}) {
  const field = useFieldContext<string>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <div>
      <Label
        htmlFor={label}
        className="mb-2 text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </Label>
      <ShadcnTextarea
        id={label}
        value={field.state.value}
        onBlur={field.handleBlur}
        rows={rows}
        onChange={(e) => field.handleChange(e.target.value)}
      />
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  )
}

interface ComboboxProps {
  label: string
  placeholder?: string,
  noItemsMessage?: string,
  items?: any[] | Group<any>[] | undefined,
}

export function Combobox<T>({ label, items, placeholder, noItemsMessage = "No items found." }: ComboboxProps) {
  const field = useFieldContext<T>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  function onValueChange(value: T, _: ComboboxRoot.ChangeEventDetails): void {
    console.log("Combobox: ", value);
    field.handleChange(value);
  }

  return (
    <Field>

      <Label
        htmlFor={label}
        className="mb-2 text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </Label>

      <ShadcnCombobox items={items} value={field.state.value} onValueChange={(value, details) => onValueChange(value as T, details)}>
        <ComboboxInput placeholder={placeholder} />
        <ComboboxContent>
          <ComboboxEmpty>{noItemsMessage}</ComboboxEmpty>
          <ComboboxList>
            {(item) => (
              <ComboboxItem key={item} value={item}>
                {item}
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </ShadcnCombobox>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </Field>
  )
}

export function Select({
  label,
  values,
  placeholder,
}: {
  label: string
  values: Array<{ label: string; value: string }>
  placeholder?: string
}) {
  const field = useFieldContext<string>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <Field>
      <FieldLabel htmlFor={label}
        className="text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </FieldLabel>
      <ShadcnSelect.Select
        name={field.name}
        value={field.state.value}
        onValueChange={(value) => field.handleChange(value)}
      >
        <ShadcnSelect.SelectTrigger className="w-full">
          <ShadcnSelect.SelectValue placeholder={placeholder} />
        </ShadcnSelect.SelectTrigger>
        <ShadcnSelect.SelectContent className="bg-background text-foreground">
          <ShadcnSelect.SelectGroup>
            <ShadcnSelect.SelectLabel>{label}</ShadcnSelect.SelectLabel>
            {values.map((value) => (
              <ShadcnSelect.SelectItem
                key={value.value}
                value={value.value}
                className="text-foreground"
              >
                {value.label}
              </ShadcnSelect.SelectItem>
            ))}
          </ShadcnSelect.SelectGroup>
        </ShadcnSelect.SelectContent>
      </ShadcnSelect.Select>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </Field>
  )
}

export function MedicationPicker({ label, values }: { label: string, values: Medication[] }) {
  const vals = values.map(x => ({ label: x.name, value: x.id.toString() }));

  return (
    <Select label={label} values={vals} />
  )
}

export function DosageMeasurementPicker({ label }: { label: string }) {
  const vals =
    MEDICATIONMEAUSEREMENTSchema.values
      .filter(x => x.number !== 0)
      .map(x => ({ label: measurementLabels[x.number] ?? x.name, value: x.name }))

  return (
    <Select label={label} values={vals} />
  )
}

export function Slider({ label }: { label: string }) {
  const field = useFieldContext<number>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <div>
      <FieldLabel
        htmlFor={label}
        className="mb-2 text-sm font-semibold text-[var(--sea-ink)]"
      >
        {label}
      </FieldLabel>
      <ShadcnSlider
        id={label}
        onBlur={field.handleBlur}
        value={[field.state.value]}
        onValueChange={(value) => field.handleChange(value[0])}
      />
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  )
}

export function Switch({ label }: { label: string }) {
  const field = useFieldContext<boolean>()
  const errors = useSelector(field.store, (state) => state.meta.errors)

  return (
    <div>
      <div className="flex items-center gap-2">
        <ShadcnSwitch
          id={label}
          onBlur={field.handleBlur}
          checked={field.state.value}
          onCheckedChange={(checked) => field.handleChange(checked)}
        />
        <Label htmlFor={label}>{label}</Label>
      </div>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  )
}

export function DateTimePicker({ label, tz }: { label: string, tz?: string }) {
  const field = useFieldContext<TZDate | undefined>()
  const errors = useSelector(field.store, (state) => state.meta.errors)
  const [open, setOpen] = useState(false)

  function setTime(e: ChangeEvent<HTMLInputElement>) {
    const time = e.target.value;
    if (!time) return;

    field.handleChange((date) =>
      set(parse(time, "HH:mm", date ?? new TZDate(new Date(), tz)), { seconds: 0, milliseconds: 0 }),
    );
  }

  function onSelect(date: Date | undefined) {
    if (!date) {
      setOpen(false)
      return
    }

    field.handleChange((prev) =>
      prev
        ? set(prev, { year: date.getFullYear(), month: date.getMonth(), date: date.getDate() })
        : new TZDate(date, tz),
    )

    setOpen(false)
  }

  return (
    <FieldGroup className="mx-auto max-w-xs flex-row">
      <Field>
        <FieldLabel htmlFor="date-picker-optional">{label}</FieldLabel>
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button variant="outline" id="date-picker-optional" className="w-32 justify-between font-normal">
              {field.state.value ? format(field.state.value, "PPP") : "Select date"}
              <ChevronDownIcon />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto overflow-hidden p-0" align="start">
            <Calendar
              mode="single"
              selected={field.state.value}
              captionLayout="dropdown"
              defaultMonth={field.state.value ?? new TZDate()}
              onSelect={onSelect}
            />
          </PopoverContent>
        </Popover>
      </Field>
      <Field className="w-32">
        <FieldLabel htmlFor="time-picker-optional">Time</FieldLabel>
        <Input
          type="time"
          id="time-picker-optional"
          className="appearance-none bg-background [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
          value={field.state.value ? format(field.state.value, "HH:mm") : ""}
          onChange={setTime}
        />
      </Field>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </FieldGroup>
  )
}

interface DateRangePickerProps {
  label: string,
  fieldGroupClass?: string,
  numberOfMonths?: number,
  onClear?: () => void,
}

export function DateRangePicker({ label, fieldGroupClass, onClear: onClearCallback, numberOfMonths = 2 }: DateRangePickerProps) {
  const field = useFieldContext<DateRange | undefined>();

  const errors = useSelector(field.store, (state) => state.meta.errors)
  const [open, setOpen] = useState(false)

  function onDateRangeChange(dateRange: DateRange | undefined) {
    field.handleChange(dateRange);
  }


  function onClear(event: React.MouseEvent<HTMLButtonElement, MouseEvent>): void {
    event.preventDefault();
    field.clearValues();
    onClearCallback?.();
  }

  const shouldShowClearButton = field.state.value !== undefined
    && field.state.value.from !== undefined
    && field.state.value.to !== undefined;

  const clearButton = <Button
    type="button"
    onClick={onClear}
    variant="ghost"
    size="icon"
    className="absolute right-1 top-1/2 size-7 -translate-y-1/2"
  >
    <XIcon />
  </Button>

  return (
    <FieldGroup className={fieldGroupClass}>
      <Field >
        <FieldLabel htmlFor="date-picker-optional">{label}</FieldLabel>
        <div className="relative">
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button variant="outline" id="date-picker-optional" className="w-full px-2.5 pr-9 flex justify-between font-normal">
                <span className='flex space-x-2'>
                  <CalendarIcon className='self-center' />
                  {field.state.value?.from ? (
                    field.state.value.to ? (
                      <>
                        <>
                          {format(field.state.value.from, "LLL dd, y")} -{" "}
                          {format(field.state.value.to, "LLL dd, y")}
                        </>
                      </>
                    ) : (
                      <>
                        {format(field.state.value.from, "LLL dd, y")}
                      </>
                    )
                  ) : (
                    <span>Pick a date</span>
                  )}
                </span>
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto overflow-hidden p-0" align="start">
              <Calendar
                mode="range"
                selected={field.state.value}
                captionLayout="dropdown"
                numberOfMonths={numberOfMonths}
                onSelect={onDateRangeChange}
              />
            </PopoverContent>
          </Popover>
          {shouldShowClearButton && clearButton}
        </div>
      </Field>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </FieldGroup>
  )
}

