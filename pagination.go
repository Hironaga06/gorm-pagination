package pagination

import (
	"math"

	"gorm.io/gorm"
)

const (
	DefaultOffset = 1
	DefaultLimit  = 10
)

type (
	pagination struct {
		db     *gorm.DB
		offset int
		limit  int
		order  []string
		models interface{}
	}

	Result struct {
		TotalRecord int64       `json:"totalRecord"`
		TotalPage   int         `json:"totalPage"`
		Records     interface{} `json:"records"`
		Offset      int         `json:"offset"`
		Limit       int         `json:"limit"`
		CurrentPage int         `json:"currentPage"`
		PrevPage    int         `json:"prevPage"`
		NextPage    int         `json:"nextPage"`
	}
)

func New(db *gorm.DB, offset, limit int, order []string, models interface{}, debug bool) *pagination {
	newDB := db
	if debug {
		newDB = db.Debug()
	}
	return &pagination{
		db:     newDB,
		offset: offset,
		limit:  limit,
		order:  order,
		models: models,
	}
}

func (p *pagination) Paging() (*Result, error) {
	db := p.db
	if p.offset < DefaultOffset {
		p.offset = DefaultOffset
	}
	if p.limit == 0 {
		p.limit = DefaultLimit
	}

	var offset int
	if p.offset == 1 {
		offset = 0
	} else {
		offset = (p.offset - 1) * p.limit
	}

	if len(p.order) > 0 {
		for _, o := range p.order {
			db = db.Order(o)
		}
	}

	count, err := p.CountRecords()
	if err != nil {
		return nil, err
	}

	if err := db.Limit(p.limit).Offset(offset).Find(p.models).Error; err != nil {
		return nil, err
	}

	var (
		prevPage, nextPage = p.offset, p.offset
		totalPage          = int(math.Ceil(float64(count) / float64(p.limit)))
	)
	if p.offset > 1 {
		prevPage = p.offset - 1
	}
	if p.offset != totalPage {
		nextPage = p.offset + 1
	}

	return &Result{
		TotalRecord: count,
		TotalPage:   totalPage,
		Records:     p.models,
		Offset:      p.offset,
		Limit:       p.limit,
		CurrentPage: p.offset,
		PrevPage:    prevPage,
		NextPage:    nextPage,
	}, nil
}

func (p *pagination) CountRecords() (int64, error) {
	var count int64
	if err := p.db.Model(p.models).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
