package service

import (
	"database/sql"
	"errors"
	"example1/internal/DTO"
	"example1/internal/repository"
)

const (
	ErrNoProducts = "no products"
)

type warehouseService struct {
	Repository repository.WarehouseRepository
}

func NewWarehouseService(r repository.WarehouseRepository) WarehouseService {
	return &warehouseService{r}
}

type WarehouseService interface {
	GetProducts(req *DTO.ReqGetProducts) (*DTO.ResGetProducts, error)
}

func (s *warehouseService) GetProducts(req *DTO.ReqGetProducts) (*DTO.ResGetProducts, error) {
	products, err := s.Repository.AllProducts(req.WarehouseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(ErrNoProducts)
		}
		return nil, err
	}
	return &DTO.ResGetProducts{ProductsCodes: products}, nil
}
