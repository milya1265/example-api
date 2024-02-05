package service

import (
	"database/sql"
	"errors"
	"example1/internal/DTO"
	"example1/internal/logger"
	"example1/internal/model"
	"example1/internal/repository"
)

var ErrInvalidUniqueCode = errors.New("invalid unique code")
var ErrInternal = errors.New("internal error")
var ErrInvalidWarehouse = errors.New("invalid warehouse id")
var ErrWarehouseUnavailable = errors.New("warehouse is unavailable")
var ErrNotEnoughProduct = errors.New("not enough product")
var ErrNonExistReservationId = errors.New("non-existent reservation id")

type productService struct {
	ProductRepository   repository.ProductRepository
	WarehouseRepository repository.WarehouseRepository
}

func NewProductService(r1 repository.ProductRepository, r2 repository.WarehouseRepository) ProductService {
	return &productService{r1, r2}
}

//go:generate mockgen -source=product.go -destination=mocks/mock.go

type ProductService interface {
	Reserve(reservation *DTO.ReqReserveProduct) (*DTO.ResReserveProduct, error)
	FreeReservation(reservations *DTO.ReqFreeReservation) (*DTO.ResFreeReservation, error)
}

func (s *productService) Reserve(reservation *DTO.ReqReserveProduct) (*DTO.ResReserveProduct, error) {

	result := &DTO.ResReserveProduct{
		Successful:   make([]DTO.Successful, 0),
		Unsuccessful: make([]string, 0),
		Errors:       make([]string, 0),
	}

	available, err := s.WarehouseRepository.CheckAvailable(reservation.WarehouseID)
	if err != nil {
		return nil, ErrInvalidWarehouse
	}

	if available == false {
		logger.ErrLog.Println(err)
		return nil, ErrWarehouseUnavailable
	}

	for i := 0; i < len(reservation.UniqueCodes); i++ {
		re := &model.Reservation{
			WarehouseID: reservation.WarehouseID,
			ProductCode: reservation.UniqueCodes[i],
			Count:       reservation.Counts[i],
		}

		leftCount, err := s.ProductRepository.GetLeftCount(re.ProductCode, re.WarehouseID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				result.Unsuccessful = append(result.Unsuccessful, re.ProductCode)
				result.Errors = append(result.Errors, ErrInvalidUniqueCode.Error())
				continue
			}
			result.Unsuccessful = append(result.Unsuccessful, re.ProductCode)
			result.Errors = append(result.Errors, ErrInternal.Error())
			continue
		}

		if leftCount < re.Count {
			result.Unsuccessful = append(result.Unsuccessful, re.ProductCode)
			result.Errors = append(result.Errors, ErrNotEnoughProduct.Error())
			continue
		}
		err = s.ProductRepository.ChangeLeftCount(re.ProductCode, re.WarehouseID, leftCount-re.Count)
		if err != nil {
			result.Unsuccessful = append(result.Unsuccessful, re.ProductCode)
			result.Errors = append(result.Errors, ErrInternal.Error())
			continue
		}
		re, err = s.ProductRepository.ReserveProduct(re)
		if err != nil {
			result.Unsuccessful = append(result.Unsuccessful, re.ProductCode)
			result.Errors = append(result.Errors, ErrInternal.Error())
		} else {
			s := DTO.Successful{
				ID:         re.ID,
				UniqueCode: re.ProductCode,
			}
			result.Successful = append(result.Successful, s)
		}

	}

	return result, nil
}

func (s *productService) FreeReservation(reservations *DTO.ReqFreeReservation) (*DTO.ResFreeReservation, error) {
	result := DTO.ResFreeReservation{
		Successful:   make([]int, 0),
		Unsuccessful: make([]int, 0),
		Errors:       make([]string, 0),
	}

	for i := 0; i < len(reservations.ID); i++ {
		warehouseID, err := s.ProductRepository.GetWarehouseByReservationID(reservations.ID[i])
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				result.Unsuccessful = append(result.Unsuccessful, reservations.ID[i])
				result.Errors = append(result.Errors, ErrNonExistReservationId.Error())
				continue
			}
			result.Unsuccessful = append(result.Unsuccessful, reservations.ID[i])
			result.Errors = append(result.Errors, ErrInternal.Error())
			continue
		}

		available, err := s.WarehouseRepository.CheckAvailable(warehouseID)
		if err != nil {
			result.Unsuccessful = append(result.Unsuccessful, reservations.ID[i])
			result.Errors = append(result.Errors, ErrInternal.Error())
			continue
		}

		if available == false {
			result.Unsuccessful = append(result.Unsuccessful, reservations.ID[i])
			result.Errors = append(result.Errors, ErrWarehouseUnavailable.Error())
			continue
		}

		err = s.ProductRepository.DeleteReservation(reservations.ID[i])
		if err != nil {
			result.Unsuccessful = append(result.Unsuccessful, reservations.ID[i])
			result.Errors = append(result.Errors, ErrInternal.Error())
			continue
		}

		result.Successful = append(result.Successful, reservations.ID[i])

	}

	return &result, nil
}
