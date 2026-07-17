package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	bloodpressurev1 "nopresh.apetrovic.com/gen/proto/bloodpressure/v1"
	"nopresh.apetrovic.com/internal/data"
	"nopresh.apetrovic.com/internal/domain/bp"
	"nopresh.apetrovic.com/internal/utils/auth"
	"nopresh.apetrovic.com/internal/utils/database"
)

type BloodPressureServer struct {
	models data.Models
	logger *slog.Logger
}

type BpReq[T any] = *connect.Request[T]
type BpRes[T any] = *connect.Response[T]

func (bps *BloodPressureServer) GetBloodPressure(
	ctx context.Context,
	req BpReq[bloodpressurev1.GetBloodPressureRequest],
) (BpRes[bloodpressurev1.GetBloodPressureResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	if req.Msg.PageInfo.PageSize == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("page size must be > 0"))
	}

	// isFirstPage := req.Msg.PageInfo.Cursor == ""
	// pointsNext := false

	var startDate time.Time
	var endDate time.Time

	if req.Msg.DateFilter != nil {
		start := req.Msg.DateFilter.StartDate
		end := req.Msg.DateFilter.EndDate

		if start != nil {
			startDate = start.AsTime()
		}

		if end != nil {
			endDate = end.AsTime()
		}

	}

	entries, cursor, err := bps.models.BloodPressure.Get(ctx, userCtx.JwtClaims.ID, uint(req.Msg.PageInfo.PageSize), req.Msg.PageInfo.Cursor, toDomainSortOrder(req.Msg.PageInfo.SortOrder), startDate, endDate)

	encodedCursor := database.EncodeCursor(cursor)

	if err != nil {
		bps.logger.Error("couldn't retrieve blood pressure info.",
			"user id", userCtx.JwtClaims.ID,
		)
		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't retrieve blood pressure info"))
	}

	entriesMapped := lo.Map(entries, func(item *bp.BloodPressure, index int) *bloodpressurev1.BloodPressureEntry {
		return fromDomainObject(item)
	})

	res := connect.NewResponse(&bloodpressurev1.GetBloodPressureResponse{
		BloodPressure: entriesMapped,
		PageInfo: &bloodpressurev1.PageInfo{
			Cursor:    encodedCursor,
			PageSize:  req.Msg.PageInfo.PageSize,
			SortOrder: req.Msg.PageInfo.SortOrder,
		},
	})

	return res, nil
}

func (bps *BloodPressureServer) CreateBloodPressure(
	ctx context.Context,
	req BpReq[bloodpressurev1.CreateBloodPressureRequest],
) (BpRes[bloodpressurev1.CreateBloodPressureResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	newBpEntry := bp.New(
		0,
		userCtx.JwtClaims.ID,
		req.Msg.DateTimeUtc.AsTime(),
		uint16(req.Msg.Systolic),
		uint16(req.Msg.Diastolic),
		uint16(req.Msg.Pulse),
		req.Msg.Dosage,
		uint(req.Msg.MedicationId),
		req.Msg.MedicationTaken,
	)

	savedBpEntry, err := bps.models.BloodPressure.Insert(&ctx, newBpEntry)

	if err != nil {
		bps.logger.Error("couldn't save new bp entry",
			"id", savedBpEntry.ID,
			"userId", savedBpEntry.UserId,
		)
		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't save new blood pressure entry"))
	}

	res := connect.NewResponse(&bloodpressurev1.CreateBloodPressureResponse{
		Entry: fromDomainObject(savedBpEntry),
	})

	return res, nil
}

func (bps *BloodPressureServer) UpdateBloodPressure(
	ctx context.Context,
	req BpReq[bloodpressurev1.UpdateBloodPressureRequest],
) (BpRes[emptypb.Empty], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	mask := req.Msg.UpdateMask

	if mask == nil || len(mask.Paths) == 0 {
		res := connect.NewResponse(&emptypb.Empty{})
		return res, nil
	}

	var input bp.UpdateDto

	for _, path := range mask.Paths {
		switch path {
		case "systolic":
			input.Systolic = new(uint16(req.Msg.Systolic))
		case "diastolic":
			input.Diastolic = new(uint16(req.Msg.Diastolic))
		case "pulse":
			input.Pulse = new(uint16(req.Msg.Pulse))
		case "date_time_utc":
			input.DateTimeUtc = new(req.Msg.DateTimeUtc.AsTime())
		case "medication_id":
			input.MedicationId = new(uint(req.Msg.MedicationId))
		case "medication_taken":
			input.MedicationTaken = new(req.Msg.MedicationTaken)
		}

	}

	err := bps.models.BloodPressure.Update(&ctx, uint(req.Msg.Id), &input, userCtx.JwtClaims.ID)

	if err != nil {
		bps.logger.Error("couldn't update blood pressure entry",
			"id", req.Msg.Id,
			"userId", userCtx.JwtClaims.ID,
			"error", err.Error(),
		)

		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't update the entry"))
	}

	res := connect.NewResponse(&emptypb.Empty{})

	return res, nil
}

func (bps *BloodPressureServer) DeleteBloodPressure(
	ctx context.Context,
	req BpReq[bloodpressurev1.DeleteBloodPressureRequest],
) (BpRes[emptypb.Empty], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	err := bps.models.BloodPressure.Delete(ctx, uint(req.Msg.Id), userCtx.JwtClaims.ID)

	if err != nil {
		bps.logger.Error("couldn't delete blood pressure entry",
			"bp entry id", req.Msg.Id,
			"userId", userCtx.JwtClaims.ID,
			"Error", err.Error(),
		)

		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't delete the entry"))
	}

	res := connect.NewResponse(&emptypb.Empty{})

	return res, nil
}

func toDomainSortOrder(m bloodpressurev1.SORTORDER) database.SortOrder {
	switch m {
	case bloodpressurev1.SORTORDER_SORTORDER_ASC:
		return database.ASC
	case bloodpressurev1.SORTORDER_SORTORDER_DESC:
		return database.DESC
	default:
		return database.DESC
	}
}

func sortOrderFromDomain(sortOrder database.SortOrder) bloodpressurev1.SORTORDER {
	switch sortOrder {
	case database.ASC:
		return bloodpressurev1.SORTORDER_SORTORDER_ASC
	case database.DESC:
		return bloodpressurev1.SORTORDER_SORTORDER_DESC
	default:
		return bloodpressurev1.SORTORDER_SORTORDER_DESC
	}
}

func fromDomainObject(bp *bp.BloodPressure) *bloodpressurev1.BloodPressureEntry {
	return &bloodpressurev1.BloodPressureEntry{
		Id:              uint32(bp.ID),
		Systolic:        uint32(bp.Systolic),
		Diastolic:       uint32(bp.Diastolic),
		Pulse:           uint32(bp.Pulse),
		UserId:          uint32(bp.UserId),
		DateTimeUtc:     timestamppb.New(bp.DateTimeUtc),
		CreatedAt:       timestamppb.New(bp.CreatedAt),
		UpdatedAt:       timestamppb.New(bp.UpdatedAt),
		MedicationId:    uint32(bp.MedicationId),
		Dosage:          bp.Dosage,
		MedicationTaken: bp.MedicationTaken,
		Medication:      medFromDomainObject(&bp.Medication),
	}
}

type BloodPressureUpdateDto struct {
	Systolic        *uint16
	Diastolic       *uint16
	Pulse           *uint16
	UserId          *uint16
	DateTimeUtc     *time.Time
	MedicationId    *uint
	dosage          *float32
	MedicationTaken *bool
}
