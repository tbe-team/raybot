package paging

type Params struct {
	PageSize uint `validate:"required,min=1"`
	Page     uint `validate:"required,min=1"`
}

type Page uint
type PageSize uint

func NewParams(page Page, pageSize PageSize, opts ...ParamsOption) Params {
	p := Params{
		Page:     uint(page),
		PageSize: uint(pageSize),
	}
	if page == 0 {
		p.Page = 1
	}
	if pageSize == 0 {
		p.PageSize = 10
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func (p Params) Offset() uint {
	return (p.Page - 1) * p.PageSize
}

func (p Params) Limit() uint {
	return p.PageSize
}

type ParamsOption func(*Params)

func WithMaxPageSize(size uint) ParamsOption {
	return func(p *Params) {
		if p.PageSize > size {
			p.PageSize = size
		}
	}
}
