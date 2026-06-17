package status

func mapStatusTone(status string) string {
	switch status {
	case "Aktif":
		return "green"
	case "Istirahat":
		return "blue"
	case "Diberhentikan":
		return "gray"
	default:
		return "blue"
	}
}

func mapOccupancyTone(label string) string {
	switch label {
	case "Zona Padat":
		return "gold"
	case "Ramai":
		return "blue"
	case "Normal":
		return "green"
	default:
		return "gray"
	}
}
