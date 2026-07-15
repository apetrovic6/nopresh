import { createFormHook } from '@tanstack/react-form'

import {
    DateRangePicker,
    DateTimePicker,
  DosageMeasurementPicker,
  MedicationPicker,
  PasswordField,
  Select,
  SubscribeButton,
  Switch,
  TextArea,
  TextField,
  UnitInput,
} from '../components/demo.FormComponents'
import { fieldContext, formContext } from './demo.form-context'

export const { useAppForm } = createFormHook({
  fieldComponents: {
    TextField,
    UnitInput,
    PasswordField,
    Select,
    TextArea,
    MedicationPicker,
    DosageMeasurementPicker,
    Switch,
    DateTimePicker,
    DateRangePicker,
  },
  formComponents: {
    SubscribeButton,
  },
  fieldContext,
  formContext,
})
