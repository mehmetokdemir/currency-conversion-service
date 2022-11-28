package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id                  uint   `gorm:"primaryKey;autoIncrement" `
	Username            string `gorm:"uniqueIndex;not null" binding:"required"`
	Email               string `gorm:"uniqueIndex;not null" binding:"required"`
	Password            string `gorm:"not null" binding:"required"`
	DefaultCurrencyCode string
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

//
// Request
//

type RegisterRequest struct {
	Username     string `json:"username" extensions:"x-order=1" example:"john" validate:"required" valid:"required~username|invalid"`          // Username of the creating user
	Email        string `json:"email" extensions:"x-order=2" example:"john@gmail.com" validate:"required" valid:"required~email|invalid"`      // Email of the creating user
	Password     string `json:"password" extensions:"x-order=3" example:"TopSecret!!!" validate:"required" valid:"required~password|invalid"`  // Password of the creating user
	CurrencyCode string `json:"currency_code" extensions:"x-order=4" example:"TRY" validate:"required" valid:"required~currency_code|invalid"` // Currency code for default wallet which is given currency
}

type LoginRequest struct {
	Username string `json:"username" extensions:"x-order=1" example:"john" validate:"required" valid:"required~username|invalid"`         // Username of the user
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
	RegisterResponse `json:",inline" extensions:"x-order=1"`
	TokenHash        string `json:"token_hash" extensions:"x-order=2" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklEIjozLCJVc2VybmFtZSI6Impob24iLCJFbWFpbCI6ImpvaG5AZ21haWwuY29tIiwiUGFzc3dvcmQiOiIkMmEkMTAkRkFUb1ZsS2Y2VmZIRGtYL1dLWmVRT0o2U1kuU3Z0SnNYYmhZV2FlTnBrbjU3S0hlNk4vZTIiLCJEZWZhdWx0Q3VycmVuY3lDb2RlIjoiIiwiY3JlYXRlZF9hdCI6IjIwMjItMTEtMjNUMjI6MjA6MDkuMzk0NzQ3KzAzOjAwIiwidXBkYXRlZF9hdCI6IjIwMjItMTEtMjNUMjI6MjA6MDkuMzk0NzQ3KzAzOjAwIiwiZGVsZXRlZF9hdCI6bnVsbH0sImV4cCI6MTY2OTM4OTM3MH0.b_i6GhYzqOp0VvouVi0rw2VG43UZx7lnJXqNEAKMH8o"`
}
