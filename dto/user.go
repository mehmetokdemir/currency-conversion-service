package dto

//
// Request
//

type RegisterRequest struct {
	Username string `json:"username" extensions:"x-order=1" example:"john" validate:"required" valid:"required~username|invalid"`         // Username of the creating user
	Email    string `json:"email" extensions:"x-order=2" example:"john@gmail.com" validate:"required" valid:"required~email|invalid"`     // Email of the creating user
	Password string `json:"password" extensions:"x-order=3" example:"TopSecret!!!" validate:"required" valid:"required~password|invalid"` // Password of the creating user
}

type LoginRequest struct {
	Username string `json:"username" extensions:"x-order=1" example:"okdemir" validate:"required" valid:"required~username|invalid"`      // Username of the user
	Password string `json:"password" extensions:"x-order=2" example:"TopSecret!!!" validate:"required" valid:"required~password|invalid"` // Password of the user
}

//
// Response
//

type RegisterResponse struct {
	Username string `json:"username" extensions:"x-order=1" example:"john"`
	Email    string `json:"email" extensions:"x-order=2" example:"john@gmail.com"`
}

type LoginResponse struct {
	RegisterResponse `json:",inline"`
	TokenHash        string `json:"token_hash" example:""`
}
