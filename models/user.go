package models

import "github.com/personal-project/pitch-league/utils"

type UserRole int

const (
	UserRoleNormal UserRole = 1
	UserRoleAdmin  UserRole = 10
)

type User struct {
	BaseModel
	Email    string   `json:"email" bun:"email"`
	Phone    string   `json:"phone" bun:"phone"`
	Name     string   `json:"name" bun:"name"`
	Surname  string   `json:"surname" bun:"surname"`
	UserName string   `json:"username" bun:"username"`
	Password string   `json:"-" bun:"password"`
	Role     UserRole `json:"role" bun:"role"`
}

// Create için kullanılacak model
type UserCreate struct {
	Email    string `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Name     string `json:"name" validate:"required,max=100"`
	Surname  string `json:"surname" validate:"required,max=100"`
	UserName string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

// ToModel creates a User from UserCreate
func (u UserCreate) ToModel() User {
	hashedPassword, _ := utils.HashPassword(u.Password)
	return User{
		Email:    utils.CleanEmail(u.Email),
		Phone:    utils.CleanPhone(u.Phone),
		Name:     utils.ToTitle(u.Name),
		Surname:  utils.ToTitle(u.Surname),
		UserName: u.UserName,
		Password: hashedPassword,
	}
}

// Update için kullanılacak model
type UserUpdate struct {
	Email    string   `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string   `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Name     string   `json:"name" validate:"required,max=100"`
	Surname  string   `json:"surname" validate:"required,max=100"`
	UserName string   `json:"username" validate:"required,max=20"`
	Role     UserRole `json:"role" validate:"required"`
	Password string   `json:"password" validate:"max=100"`
}

// ToModel updates an existing User from UserUpdate
func (u UserUpdate) ToModel(existing User) User {
	existing.Email = utils.CleanEmail(u.Email)
	existing.Phone = utils.CleanPhone(u.Phone)
	existing.Name = utils.ToTitle(u.Name)
	existing.Surname = utils.ToTitle(u.Surname)
	existing.UserName = u.UserName
	existing.Role = u.Role

	if u.Password != "" {
		hashedPassword, _ := utils.HashPassword(u.Password)
		existing.Password = hashedPassword
	}

	return existing
}

// Response için kullanılacak model
type UserResponse struct {
	ID       int64    `json:"id"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	UserName string   `json:"username"`
	Role     UserRole `json:"role"`
}

func ToUserResponse(u User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Phone:    u.Phone,
		Name:     u.Name,
		Surname:  u.Surname,
		UserName: u.UserName,
		Role:     u.Role,
	}
}

func (User) ModelName() string {
	return "user"
}

func (u User) String() string {
	return u.Name + " " + u.Surname
}

func (r UserRole) String() string {
	switch r {
	case UserRoleNormal:
		return "normal"
	case UserRoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

func (User) TableName() string {
	return "users"
}
