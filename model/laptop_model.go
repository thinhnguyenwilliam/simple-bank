// simple-bank/model/laptop_model.go
package model

import (
	"time"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
)

type LaptopModel struct {
	ID string `gorm:"primaryKey;type:uuid"`

	Brand string `gorm:"not null;index"`
	Name  string `gorm:"not null"`

	PriceUSD    float64 `gorm:"not null;index"`
	ReleaseYear uint32  `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToLaptopModel(p *pb.Laptop) *LaptopModel {
	return &LaptopModel{
		ID:          p.Id,
		Brand:       p.Brand,
		Name:        p.Name,
		PriceUSD:    p.PriceUsd,
		ReleaseYear: p.ReleaseYear,
	}
}
