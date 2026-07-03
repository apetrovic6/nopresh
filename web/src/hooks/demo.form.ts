import { createFormHook } from '@tanstack/react-form'

import {
  DosageMeasurementPicker,
  MedicationPicker,
  PasswordField,
  Select,
  SubscribeButton,
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
  },
  formComponents: {
    SubscribeButton,
  },
  fieldContext,
  formContext,
})
