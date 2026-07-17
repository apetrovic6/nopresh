import { MEDICATIONMEAUSEREMENT, MEDICATIONMEAUSEREMENTSchema, type Medication } from "#/gen/proto/medication/v1/medication_entry_pb";

export const measurementLabels: Partial<Record<number, string>> = {
  [MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_MG]: "mg",
  [MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_G]: "g",
}


export namespace MedicationUtils {
  export function findById(medicationId: number, medications: Medication[]) {
    return medications.find(x => x.id === medicationId);
  }

  export function getMeasurementNameByMedId(medicationId: number, medications: Medication[]) {
    const measurement = findById(medicationId, medications)?.dosageMeasurement ?? MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_UNSPECIFIED;
    return measurementLabels[measurement];
  }

  export function getMeasurementByEnum(medicationMeasurement: MEDICATIONMEAUSEREMENT) {
    return measurementLabels[medicationMeasurement];
  }


  export function getMedicationDosage(medicationId: number, medications: Medication[]) {
    return findById(medicationId, medications)?.recommendedDosage
  }

  export function getDosageMeasurementIdByName(name: string) {
    return MEDICATIONMEAUSEREMENTSchema.values
      .find(x => x.name === name)?.number
      ?? MEDICATIONMEAUSEREMENT.MEDICATIONMEAUSEREMENT_UNSPECIFIED;
  }

  export function getNameByEnum(value: MEDICATIONMEAUSEREMENT) {
    return MEDICATIONMEAUSEREMENTSchema.values
      .find(x => x.number === value)?.name
  }
}
