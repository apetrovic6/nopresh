package main

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
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
	req MedReq[medicationv1.CreateRequest],
) (MedRes[medicationv1.CreateResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing auth info"))
	}

	newMedEntry := medication.New(
		0,
		req.Msg.Name,
		req.Msg.RecommendedDosage,
		medication.ToDomainMeasurment(req.Msg.DosageMeasurment),
		userCtx.JwtClaims.ID,
	)

	savedEntry, err := ms.models.Medication.Insert(newMedEntry)

	if err != nil {
		ms.logger.Error("couldn't save new medication entry",
			"id", savedEntry.ID,
			"name", savedEntry.Name,
			"userId", savedEntry.UserId,
		)
		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't save new blood pressure entry"))
	}

	res := connect.NewResponse(&medicationv1.CreateResponse{
		Entry: medFromDomainObject(savedEntry),
	})

	return res, nil
}

func medFromDomainObject(med *medication.Medication) *medicationv1.MedicationEntry {
	return &medicationv1.MedicationEntry{
		Id:                uint32(med.ID),
		Name:              med.Name,
		RecommendedDosage: med.RecommendedDosage,
		DosageMeasurment:  medication.ToPbMeasurment(med.DosageMeasurment),
	}
}
