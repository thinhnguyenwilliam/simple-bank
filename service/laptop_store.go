// simple-bank/service/laptop_store.go
package service

import (
	"errors"

	"github.com/thinhcompany/simple-bank/model"
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"gorm.io/gorm"
)

var ErrAlreadyExists = errors.New("laptop already exists")

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
}

type DBLaptopStore struct {
	db *gorm.DB
}

func NewDBLaptopStore(db *gorm.DB) *DBLaptopStore {
	return &DBLaptopStore{db: db}
}

func (store *DBLaptopStore) Save(laptop *pb.Laptop) error {
	laptopModel := &model.LaptopModel{
		ID:    laptop.Id,
		Brand: laptop.Brand,
		Name:  laptop.Name,
	}

	err := store.db.Create(laptopModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}
