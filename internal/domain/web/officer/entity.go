package officer

type ParkingOfficerOption struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Role              string `json:"role"`
	HomeZone          string `json:"homeZone"`
	Availability      string `json:"availability"`
	AvailabilityNote  string `json:"availabilityNote"`
	CurrentAssignment string `json:"currentAssignment"`
	CurrentLocationID string `json:"currentLocationId"`
	CurrentShiftID    string `json:"currentShiftId"`
	Status            string `json:"status"`
	DefaultShiftStart string `json:"defaultShiftStart"`
	DefaultShiftEnd   string `json:"defaultShiftEnd"`
	DefaultStatus     string `json:"defaultStatus"`
}
