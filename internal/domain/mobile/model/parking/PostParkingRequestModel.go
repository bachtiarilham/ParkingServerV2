package parking

type PostParkingRequestModel struct {
	OfficerUserId   int64
	PlateNumber     string
	VehicleTypeCode string
	SelectedAreaId  int64
	BiayaParkir     int64
}
