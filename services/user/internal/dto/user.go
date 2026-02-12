package dto

type LoginRequestDTO struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
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
