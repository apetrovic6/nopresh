package main

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	medicationv1 "nopresh.apetrovic.com/gen/proto/medication/v1"
	"nopresh.apetrovic.com/internal/data"
	"nopresh.apetrovic.com/internal/domain/medication"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type MedicationServer struct {
	models data.Models
	logger *slog.Logger
}

type MedReq[T any] = *connect.Request[T]
type MedRes[T any] = *connect.Response[T]

func (ms *MedicationServer) CreateMedication(
	ctx context.Context,
	req MedReq[medicationv1.CreateMedicationRequest],
) (MedRes[medicationv1.CreateMedicationResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	newMedEntry := medication.New(
		0,
		req.Msg.Name,
		req.Msg.RecommendedDosage,
		medication.ToDomainMeasurement(req.Msg.DosageMeasurement),
		userCtx.JwtClaims.ID,
	)

	savedEntry, err := ms.models.Medication.Insert(ctx, newMedEntry)

	if err != nil {
		ms.logger.Error("couldn't save new medication entry",
			"id", savedEntry.ID,
			"name", savedEntry.Name,
			"userId", savedEntry.UserId,
		)
		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't save new blood pressure entry"))
	}

	res := connect.NewResponse(&medicationv1.CreateMedicationResponse{
		Medication: medFromDomainObject(savedEntry),
	})

	return res, nil
}

func (ms *MedicationServer) GetMedication(
	ctx context.Context,
	req MedReq[medicationv1.GetMedicationRequest],
) (MedRes[medicationv1.GetMedicationResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	medData, err := ms.models.Medication.GetById(ctx, uint(req.Msg.Id), userCtx.JwtClaims.ID)

	if err != nil {
		ms.logger.Error("couldn't get medication",
			"user", userCtx.JwtClaims.ID,
			"medication", req.Msg.Id,
			"error", err.Error(),
		)

		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't fetch medication"))
	}

	resp := connect.NewResponse(&medicationv1.GetMedicationResponse{
		Medication: medFromDomainObject(medData),
	})

	return resp, nil
}

func (ms *MedicationServer) GetMedications(
	ctx context.Context,
	req MedReq[medicationv1.GetMedicationsRequest],
) (MedRes[medicationv1.GetMedicationsResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	medications, err := ms.models.Medication.GetAll(ctx, userCtx.JwtClaims.ID)

	if err != nil {
		return nil, err
	}

	medResponse := lo.Map(medications, func(med *medication.Medication, i int) *medicationv1.Medication {
		return medFromDomainObject(med)
	})

	res := connect.NewResponse(&medicationv1.GetMedicationsResponse{
		Medications: medResponse,
	})

	return res, nil

}

func (ms *MedicationServer) UpdateMedication(
	ctx context.Context,
	req MedReq[medicationv1.UpdateMedicationRequest],
) (MedRes[medicationv1.UpdateMedicationResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	mask := req.Msg.UpdateMask

	var input medication.UpdateInput

	for _, path := range mask.Paths {
		switch path {
		case "name":
			input.Name = &req.Msg.Name
		case "recommended_dosage":
			input.RecommendedDosage = &req.Msg.RecommendedDosage
		case "dosage_measurement":
			input.DosageMeasurement = new(medication.ToDomainMeasurement(req.Msg.DosageMeasurement))
		}
	}

	updatedMedication, err := ms.models.Medication.Update(ctx, uint(req.Msg.Id), input, userCtx.JwtClaims.ID)

	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&medicationv1.UpdateMedicationResponse{
		Medication: medFromDomainObject(updatedMedication),
	})

	return res, nil
}

func (ms *MedicationServer) DeleteMedication(
	ctx context.Context,
	req MedReq[medicationv1.DeleteMedicationRequest],
) (MedRes[emptypb.Empty], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	_, err := ms.models.Medication.DeleteById(ctx, uint(req.Msg.Id), userCtx.JwtClaims.ID)

	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := connect.NewResponse(&emptypb.Empty{})

	return res, nil
}

func toPbMeasurement(m medication.Measurement) medicationv1.MEDICATIONMEAUSEREMENT {
	switch m {
	case medication.Mg:
		return medicationv1.MEDICATIONMEAUSEREMENT_MEDICATIONMEAUSEREMENT_MG
	case medication.G:
		return medicationv1.MEDICATIONMEAUSEREMENT_MEDICATIONMEAUSEREMENT_G
	default:
		return medicationv1.MEDICATIONMEAUSEREMENT_MEDICATIONMEAUSEREMENT_UNSPECIFIED
	}
}

func medFromDomainObject(med *medication.Medication) *medicationv1.Medication {
	return &medicationv1.Medication{
		Id:                uint32(med.ID),
		Name:              med.Name,
		RecommendedDosage: med.RecommendedDosage,
		DosageMeasurement: toPbMeasurement(med.DosageMeasurement),
	}
}
