package user_service

import (
	identity_v1 "github.com/ducthangng/geofleet-proto/gen/go/identity/v1"
	"github.com/ducthangng/geofleet/gateway/app/role"
)

// copy from protobuf
type UserCreation struct {
	Fullname string               `protobuf:"bytes,1,opt,name=fullname,proto3" json:"fullname,omitempty"`
	Email    string               `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Phone    string               `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	Password string               `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	Address  string               `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	Bod      string               `protobuf:"bytes,6,opt,name=bod,proto3" json:"bod,omitempty"` // Dùng Timestamp thay vì string
	Role     identity_v1.UserRole `protobuf:"varint,7,opt,name=role,proto3,enum=geofleet.identity.v1.UserRole" json:"role,omitempty"`
}

// copy from protobuf
type CheckDuplicatePhoneRequest struct {
	Phone string `json:"phone"`
}

// copy from protobuf
type CheckDuplicatePhoneResponse struct {
	Phone        string `json:"phone"`
	IsDuplicated bool   `json:"isDuplicated"`
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
}

// copy from protobuf
type LoginRequest struct {
	Phone    string   `json:"phone"`
	Password Password `json:"password"`
}

type Password struct {
	Value string `json:"value"`
}

// copy from protobuf
type LoginResponse struct {
	IsValid  bool   `json:"isValid"`
	User     User   `json:"user"`
	JWTToken string `json:"jwtToken"`
}

type User struct {
	UserId   string    `json:"userId"`
	Fullname string    `json:"fullname"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Bod      string    `json:"bod"`
	Score    float64   `json:"score"`
	Address  string    `json:"address"`
	Role     role.Role `json:"role"`
}
