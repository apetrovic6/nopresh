import { MEDICATIONMEAUSEREMENT } from "#/gen/proto/medication/v1/medication_pb";

export const measurementLabels: Partial<Record<number, string>> = {
  [MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_MG]: "mg",
  [MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_G]: "g",
}
