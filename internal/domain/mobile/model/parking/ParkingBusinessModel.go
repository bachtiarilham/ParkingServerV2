package parking

type ParkingBusinessModel struct {
	//step 1 jukir input kendaraan
	PlateNumber     string
	OfficerUserId   int64
	LocationId      int64
	LocationName    string
	Address         string
	ZoneId          int64
	ZoneName        string
	AreaId          int64
	AreaName        string
	VehicleTypeId   int64
	VehicleTypeCode string
	VehicleTypeName string
	Amount          int64
	ParkingStatusId int64
	PaymentStatusId int64
	//step 2 insert parking session
	SessionCode     string
	TransactionCode string
	//step 3 simpan qris ke parking session
	SessionId         int64
	PaymentCode       string
	QrisString        string
	ExternalReference string
	ProviderName      string
}
