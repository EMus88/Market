package model

import uuid "github.com/gofrs/uuid"

type User struct {
	ID       uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" json:"-" `
	Username string    `gorm:"type:varchar(150); not null; unique" json:"username" binding:"required"`
	Phone    string    `gorm:"type:varchar(150); not null; unique" json:"phone" binding:"required" valid:"numeric"`
	Password string    `gorm:"type:varchar(150); not null" json:"password" binding:"required"`
	Role     string    `gorm:"type:varchar(150); default:'user'" json:"role,omitempty"`
}

type Admin struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}
