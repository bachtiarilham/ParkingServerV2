package settings

type SettingsDto struct {
	FotoProfil *string `json:"foto_profil,omitempty"`
	Email      *string `json:"email,omitempty"`
	NoTelp     *string `json:"noTelp,omitempty"`
}
