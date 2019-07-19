package dao

import (
	"errors"
	"github.com/llairunhao/abiu"
)

type Product struct {
	ProductID    *abiu.Int64  `json:"productID,omitempty" orm:"product_id"`
	Name         *abiu.String `json:"name,omitempty" orm:"product_name"`
	MaxPrice     *abiu.UInt32 `json:"maxPrice,omitempty" orm:"max_price"`
	MinPrice     *abiu.UInt32 `json:"minPrice,omitempty" orm:"min_price"`
	Shipping     *abiu.String `json:"shipping,omitempty" orm:"shipping"`
	Rate         *abiu.String `json:"rate,omitempty" orm:"rate"`
	OriginURL    *abiu.String `json:"originURL,omitempty" orm:"origin_url"`
	ReviewCount  *abiu.UInt   `json:"reviewCount,omitempty" orm:"review_count"`
	OrderCount   *abiu.UInt   `json:"orderCount,omitempty" orm:"order_count"`
	ModifiedDate *abiu.Date   `json:"modifiedDate,omitempty" orm:"gmt_modified"`
	CreateDate   *abiu.Date   `json:"createDate,omitempty" orm:"gmt_create"`
}

func NewProduct() *Product {
	return &Product{
		ProductID:    abiu.NewInt64(0),
		Name:         abiu.NewString(""),
		MaxPrice:     abiu.NewUInt32(0),
		MinPrice:     abiu.NewUInt32(0),
		Shipping:     abiu.NewString(""),
		Rate:         abiu.NewString(""),
		OriginURL:    abiu.NewString(""),
		ReviewCount:  abiu.NewUInt(0),
		OrderCount:   abiu.NewUInt(0),
		ModifiedDate: &abiu.Date{},
		CreateDate:   &abiu.Date{},
	}
}

func (p *Product) SaveValid() error {
	if p.Name.IsEmpty() {
		return errors.New("name 不能为空")
	}
	if p.MaxPrice.IsZero() {
		return errors.New("maxPrice 不能为空")
	}
	if p.MinPrice.IsZero() {
		return errors.New("minPrice 不能为空")
	}
	if p.Shipping.IsEmpty() {
		return errors.New("shipping 不能为空")
	}
	if p.Rate.IsEmpty() {
		return errors.New("rate 不能为空")
	}
	if p.OriginURL.IsEmpty() {
		return errors.New("originURL 不能为空")
	}
	return nil
}

func (p *Product) UpdateValid() error {
	if p.ProductID.IsZero() {
		return errors.New("productId 不能为0")
	}
	if p.Name == nil &&
		p.MaxPrice == nil &&
		p.MinPrice == nil &&
		p.Shipping == nil &&
		p.Rate == nil &&
		p.OrderCount == nil &&
		p.ReviewCount == nil {
		return errors.New("没有更新字段")
	}
	return nil
}

func (p *Product) RemoveValid() error {
	if p.ProductID.IsZero() {
		return errors.New("productId 不能为0")
	}
	return nil
}
