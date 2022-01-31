package kom_arangodb_paginator

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
)

type Paginator struct {
	col driver.Collection

	opts PaginatorOptions
}

type PaginatorOptions struct {
	Limit  int
	Offset int

	ShowQuery bool
}

type PaginationSafeResult struct {
	CurrentPage  int
	PrevPage     int
	NextPage     int
	TotalPage    int
	TotalRecords int64
}

type PaginationResult struct {
	Result PaginationSafeResult
	Data   []interface{}
}

type PaginatorStrategySelector interface {
	Classic() PaginatorClassic
	TypeSafe() PaginatorWithTypeSafeness
}

type PaginatorClassic interface {
	Paginate() PaginationResult
}

type PaginatorWithTypeSafeness interface {
	PaginateSafely(doc Clonable, af func(interface{})) PaginationSafeResult
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
func (p *Paginator) PaginateSafely(doc Clonable, af func(interface{})) PaginationSafeResult {
	var ctx = context.Background()
	var ql = prepareQuery(p.col.Name(), p.opts)

	if err := p.col.Database().ValidateQuery(ctx, ql); err != nil {
		return PaginationSafeResult{} // return err
	}

	var cursor, err = p.col.Database().Query(ctx, ql, nil)
	if err != nil {
		return PaginationSafeResult{} // return err
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

	return PaginationSafeResult{}
}

func prepareQuery(collectionName string, opts PaginatorOptions) string {
	if opts.Offset == 0 && opts.Limit == 0 {
		opts.Limit = 10
	}

	return fmt.Sprintf("FOR d IN %s LIMIT %d,%d RETURN d", collectionName, opts.Offset, opts.Limit)
}
