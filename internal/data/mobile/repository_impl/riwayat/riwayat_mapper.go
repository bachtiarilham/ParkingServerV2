package riwayat

import (
	"database/sql"
	"strings"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

type riwayatRow struct {
	SectionDate       sql.NullString
	TransactionID     sql.NullInt64
	TransactionCode   sql.NullString
	SessionID         sql.NullInt64
	PlateNumber       sql.NullString
	VehicleTypeID     sql.NullInt64
	VehicleTypeCode   sql.NullString
	VehicleTypeName   sql.NullString
	PaymentMethodID   sql.NullInt64
	PaymentMethodCode sql.NullString
	PaymentMethodName sql.NullString
	LocationID        sql.NullInt64
	LocationName      sql.NullString
	LocationAddress   sql.NullString
	AreaID            sql.NullInt64
	AreaName          sql.NullString
	ZoneID            sql.NullInt64
	ZoneName          sql.NullString
	BaseAmount        sql.NullInt64
	DiscountAmount    sql.NullInt64
	FinalAmount       sql.NullInt64
	CompanyShare      sql.NullInt64
	JukirShare        sql.NullInt64
	TaxAmount         sql.NullInt64
	FeeAmount         sql.NullInt64
	TransactionStatus sql.NullString
	OperationType     sql.NullString
	OccurredAt        sql.NullTime
	PaidAt            sql.NullTime
	CreatedAt         sql.NullTime
}

func mapRiwayatRowToItem(row riwayatRow) model.RiwayatItemModel {
	return model.RiwayatItemModel{
		Code:        utils.NullStringValue(row.TransactionCode),
		PlateNumber: utils.NullStringValue(row.PlateNumber),
		VehicleType: utils.NullStringValue(row.VehicleTypeName),
		Time:        utils.FormatRiwayatTime(row.PaidAt, row.OccurredAt, row.CreatedAt),
		Amount:      utils.NullInt64Value(row.FinalAmount),
		IsEntry:     isEntryOperation(row.OperationType),
	}
}

func mapRiwayatSectionDate(row riwayatRow) string {
	return utils.NullStringValue(row.SectionDate)
}

func normalizeFilterCode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "SEMUA"
	}

	return strings.ToUpper(value)
}

func isEntryOperation(operation sql.NullString) bool {
	if !operation.Valid {
		return false
	}
	return strings.EqualFold(operation.String, "ENTRY")
}
