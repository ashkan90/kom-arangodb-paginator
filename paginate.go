package kom_arangodb_paginator

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"log"
	"time"
)

type Paginator struct {
	col driver.Collection

	opts PaginatorOptions
}

type PaginatorOptions struct {
	Limit  int
	Offset int
}

type PaginationResult struct{}

type PaginatorStrategySelector interface {
	Classic() PaginatorClassic
	TypeSafe() PaginatorWithTypeSafeness
}

type PaginatorClassic interface {
	Paginate() PaginationResult
}

type PaginatorWithTypeSafeness interface {
	PaginateSafely(doc Clonable, af func(interface{}))
}

func New(col driver.Collection, opts PaginatorOptions) PaginatorStrategySelector {
	return &Paginator{col: col, opts: opts}
}

func (p *Paginator) Classic() PaginatorClassic {
	return p
}

func (p *Paginator) TypeSafe() PaginatorWithTypeSafeness {
	return p
}

func (p *Paginator) Paginate() PaginationResult {
	return PaginationResult{}
}

// PaginateSafely returns total count of documents
func (p *Paginator) PaginateSafely(doc Clonable, af func(interface{})) {
	var ctx = context.Background()
	var ql = prepareQuery(p.col.Name(), p.opts)

	log.Println(ql)

	if err := p.col.Database().ValidateQuery(ctx, ql); err != nil {
		log.Println(err)
		return // return err
	}

	var cursor, err = p.col.Database().Query(ctx, ql, nil)
	if err != nil {
		log.Println(err)
		return // return err
	}

	defer cursor.Close()
	for {
		var newDoc = doc.Clone()
		_, err = cursor.ReadDocument(ctx, newDoc)

		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
		}

		af(newDoc)
	}
}

func prepareQuery(collectionName string, opts PaginatorOptions) string {
	if opts.Offset == 0 && opts.Limit == 0 {
		opts.Limit = 10
	}
	
	return fmt.Sprintf("FOR d IN %s LIMIT %d,%d RETURN d", collectionName, opts.Offset, opts.Limit)
}

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

type Clonable interface {
	Clone() interface{}
}

func (p *Product) Clone() interface{} {
	return &Product{}
}
