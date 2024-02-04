package repository

import (
	"database/sql"
	"example1/internal/logger"
	"example1/internal/model"
)

type productRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
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
	query := `INSERT INTO reservation (warehouse_id, product_code, count) VALUES ($1, $2, $3) RETURNING id;`
	row := r.DB.QueryRow(query, reservation.WarehouseID, reservation.ProductCode, reservation.Count)
	err := row.Scan(&reservation.ID)
	if err != nil {
		logger.ErrLog.Println(err)
		return nil, err
	}
	return reservation, nil
}

func (r *productRepository) GetLeftCount(uniqueCode string, warehouseId int) (int, error) {
	query := `SELECT (left_count) FROM warehouse_product WHERE product_code = $1 AND warehouse_id = $2;`
	count := 0
	err := r.DB.QueryRow(query, uniqueCode, warehouseId).Scan(&count)
	if err != nil {
		logger.ErrLog.Println(err)
		return 0, err
	}
	return count, nil
}

func (r *productRepository) ChangeLeftCount(uniqueCode string, warehouseID, count int) error {
	query := `UPDATE warehouse_product SET left_count = $1 WHERE product_code = $2 AND warehouse_id = $3;`
	_, err := r.DB.Exec(query, count, uniqueCode, warehouseID)
	if err != nil {
		logger.ErrLog.Println(err)
		return err
	}
	return nil
}

func (r *productRepository) DeleteReservation(reservationID int) error {
	query := `DELETE FROM reservation WHERE id = $1;`
	_, err := r.DB.Exec(query, reservationID)
	if err != nil {
		logger.ErrLog.Println(err)
		return err
	}
	return nil
}

func (r *productRepository) GetWarehouseByReservationID(resID int) (int, error) {
	query := `SELECT (warehouse_id) FROM reservation WHERE id = $1;`
	id := 0
	err := r.DB.QueryRow(query, resID).Scan(&id)
	if err != nil {
		logger.ErrLog.Println(err)
		return 0, err
	}
	return id, nil
}
