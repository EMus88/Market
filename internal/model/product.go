package model

import uuid "github.com/gofrs/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()"`
	Name        string    `gorm:"index:prod_name;type:varchar(150); not null; unique"`
	Weight      float32   `gorm:"not null"`
	Valume      float32   `gorm:"not null"`
	Description string    `gorm:"type:varchar(255)"`
	Photo       []string  `gorm:"type:text[]"`
	Price       uint64    `gorm:"not null"`
	Visible     bool      `gorm:"default:true"`
	CategoryID  uuid.UUID `gorm:"type:uuid; not null"`
}

type ProductDTO struct {
	Name        string   `json:"name" binding:"required" valid:"alpha"`
	Weight      float32  `json:"weight" binding:"required"`
	Valume      float32  `json:"valume" binding:"required"`
	Description string   `json:"description,omitempty" `
	Photo       []string `json:"photo,omitempty"`
	Price       float32  `json:"price" binding:"required"`
	Visible     bool     `json:"visible,omitempty"`
	Category    string   `json:"category"`
}

type Visible struct {
	Name    string `json:"name" binding:"required" valid:"alpha"`
	Visible bool   `json:"visible"`
}
