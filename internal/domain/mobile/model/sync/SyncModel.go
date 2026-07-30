package sync

type SyncModel struct {
	LocationID        int64
	AreaID            int64
	ZoneID            int64
	OfficerUserID     int64
	VehicleTypeCode   string
	Amount            int64
	PartnerShare      int64
	CompanyShare      int64
	GovShare          int64
	PaymentMethodCode string
}
