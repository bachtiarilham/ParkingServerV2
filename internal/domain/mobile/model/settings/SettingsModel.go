package settings

type SettingsModel struct {
	UserId          int64
	FotoProfil      *string `json:"foto_profil,omitempty"`
	Email           *string `json:"email,omitempty"`
	NoTelp          *string `json:"noTelp,omitempty"`
	FotoProfilBytes []byte
	FotoProfilName  string
}
