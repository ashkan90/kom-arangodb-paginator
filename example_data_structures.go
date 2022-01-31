package kom_arangodb_paginator

import "time"

type Product struct {
	ID        int            `json:"Id,omitempty" bson:"_id,omitempty" gorm:"primary_key"`
	Title     string         `json:"Title" bson:"Title" validate:"required"`
	Detail    ProductDetail  `json:"Detail" bson:"Detail" gorm:"-"`
	Prices    []ProductPrice `json:"Prices" bson:"Prices"`
	Inventory int            `json:"Inventory" bson:"Inventory"`
	CreatedAt time.Time      `json:"CreatedAt" bson:"UpdatedAt"`
	UpdatedAt *time.Time     `json:"UpdatedAt" bson:"UpdatedAt" gorm:"index"`
}

type ProductDetail struct {
	ID        int        `json:"-,omitempty" bson:"_id,omitempty" gorm:"primary_key"`
	ProductID int        `json:"-" bson:"ProductId" gorm:"foreign_key"`
	Active    bool       `json:"Active" bson:"active"`
	Barcode   string     `json:"Barcode" bson:"barcode"`
	Brand     string     `json:"Brand" bson:"brand"`
	Image     string     `json:"Image" bson:"image"`
	Name      string     `json:"Name" bson:"name"`
	UpdatedAt *time.Time `json:"UpdatedAt" bson:"UpdatedAt" gorm:"index"`
}

type ProductPrice struct {
	ProductID int `json:"ProductId" bson:"ProductId" gorm:"foreign_key"`
	Currency  string
	Value     float64
}

func (p *Product) Clone() interface{} {
	return &Product{}
}
