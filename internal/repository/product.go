package repository

import (
	"database/sql"
	"example1/internal/model"
	"example1/pkg/logger"
)

type productRepository struct {
	DB     *sql.DB
	Logger logger.Logger
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db, logger.Get()}
}

type ProductRepository interface {
	ReserveProduct(reservation *model.Reservation) (*model.Reservation, error)
	GetLeftCount(uniqueCode string, warehouseId int) (int, error)
	ChangeLeftCount(uniqueCode string, warehouseId, count int) error
	DeleteReservation(reservationID int) error
	GetWarehouseByReservationID(resID int) (int, error)
}

func (r *productRepository) ReduceCountOfProduct() {

}

func (r *productRepository) ReserveProduct(reservation *model.Reservation) (*model.Reservation, error) {
	r.Logger.Info("start repository ReserveProduct")

	query := `INSERT INTO reservation (warehouse_id, product_code, count) VALUES ($1, $2, $3) RETURNING id;`
	row := r.DB.QueryRow(query, reservation.WarehouseID, reservation.ProductCode, reservation.Count)
	err := row.Scan(&reservation.ID)
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}
	return reservation, nil
}

func (r *productRepository) GetLeftCount(uniqueCode string, warehouseId int) (int, error) {
	r.Logger.Info("start repository GetLeftCount")

	query := `SELECT (left_count) FROM warehouse_product WHERE product_code = $1 AND warehouse_id = $2;`
	count := 0
	err := r.DB.QueryRow(query, uniqueCode, warehouseId).Scan(&count)
	if err != nil {
		r.Logger.Error(err)
		return 0, err
	}
	return count, nil
}

func (r *productRepository) ChangeLeftCount(uniqueCode string, warehouseID, count int) error {
	r.Logger.Info("start repository ChangeLeftCount")

	query := `UPDATE warehouse_product SET left_count = $1 WHERE product_code = $2 AND warehouse_id = $3;`
	_, err := r.DB.Exec(query, count, uniqueCode, warehouseID)
	if err != nil {
		r.Logger.Error(err)
		return err
	}
	return nil
}

func (r *productRepository) DeleteReservation(reservationID int) error {
	r.Logger.Info("start repository DeleteReservation")

	query := `DELETE FROM reservation WHERE id = $1;`
	_, err := r.DB.Exec(query, reservationID)
	if err != nil {
		r.Logger.Error(err)
		return err
	}
	return nil
}

func (r *productRepository) GetWarehouseByReservationID(resID int) (int, error) {
	r.Logger.Info("start repository GetWarehouseByReservationID")

	query := `SELECT (warehouse_id) FROM reservation WHERE id = $1;`
	id := 0
	err := r.DB.QueryRow(query, resID).Scan(&id)
	if err != nil {
		r.Logger.Error(err)
		return 0, err
	}
	return id, nil
}
