package location

type LocationAggregate2 struct {
	ID                int64
	Slug              string
	Name              string
	Zone              string
	Address           string
	OperationType     string
	Lat               float64
	Lng               float64
	Motorcycles       int64
	Cars              int64
	Transactions      int64
	Revenue           int64
	Officers          int64
	OfficerName       string
	OfficerStatus     string
	OfficerShiftStart string
	OfficerShiftEnd   string
	TariffMotor       int64
	TariffMobil       int64
	OperationalNote   string
	ShiftTemplates    []ParkingShiftTemplate
}

type ParkingLocation struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Zone              string                 `json:"zone"`
	Address           string                 `json:"address"`
	Lat               float64                `json:"lat"`
	Lng               float64                `json:"lng"`
	OfficerName       string                 `json:"officerName"`
	OfficerShiftStart string                 `json:"officerShiftStart"`
	OfficerShiftEnd   string                 `json:"officerShiftEnd"`
	OfficerStatus     string                 `json:"officerStatus"`
	DismissalReason   string                 `json:"dismissalReason"`
	TariffMotor       int64                  `json:"tariffMotor"`
	TariffMobil       int64                  `json:"tariffMobil"`
	Motorcycles       int64                  `json:"motorcycles"`
	Cars              int64                  `json:"cars"`
	Officers          int64                  `json:"officers"`
	OccupancyLabel    string                 `json:"occupancyLabel"`
	ShiftTemplates    []ParkingShiftTemplate `json:"shiftTemplates"`
}

type ParkingShiftTemplate struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Start string `json:"start"`
	End   string `json:"end"`
}
