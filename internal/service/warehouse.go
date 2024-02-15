package service

import (
	"database/sql"
	"errors"
	"example1/internal/DTO"
	"example1/internal/repository"
	"example1/pkg/logger"
)

var ErrNoProducts = errors.New("no products")

type warehouseService struct {
	Repository repository.WarehouseRepository
	Logger     logger.Logger
}

func NewWarehouseService(r repository.WarehouseRepository) WarehouseService {
	return &warehouseService{r, logger.Get()}
}

type WarehouseService interface {
	GetProducts(req *DTO.ReqGetProducts) (*DTO.ResGetProducts, error)
}

func (s *warehouseService) GetProducts(req *DTO.ReqGetProducts) (*DTO.ResGetProducts, error) {
	s.Logger.Info("start service GetProducts")
	products, err := s.Repository.AllProducts(req.WarehouseID)
	if errors.Is(err, sql.ErrNoRows) || len(products) == 0 {
		return nil, ErrNoProducts
	}
	if err != nil {
		return nil, err
	}
	return &DTO.ResGetProducts{ProductsCodes: products}, nil
}
