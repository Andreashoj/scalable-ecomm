package dto

type LoginRequestDTO struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterUserDTO struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LogoutRequestDTO struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refreshToken"`
}
