package dto

type LoginRequestDTO struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
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

type RefreshResponseDTO struct {
	AccessToken string `json:"accessToken"`
}
