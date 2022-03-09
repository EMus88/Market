package models

import uuid "github.com/gofrs/uuid"

type Category struct {
	ID   uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; default:uuid_generate_v4()" json:"id,omitempty" `
	Name string    `gorm:"index:indx_category;type:varchar(150);column:category; not null; unique" json:"name" binding:"required"`
}
