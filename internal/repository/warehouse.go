package repository

import (
	"database/sql"
	"example1/internal/logger"
)

type warehouseRepository struct {
	DB *sql.DB
}

func NewWarehouseRepository(db *sql.DB) WarehouseRepository {
	return &warehouseRepository{db}
}

type WarehouseRepository interface {
	CheckAvailable(warehouseID int) (bool, error)
	AllProducts(warehouseID int) ([]string, error)
}

func (r *warehouseRepository) CheckAvailable(warehouseID int) (bool, error) {
	query := `SELECT (available) FROM warehouse WHERE id = $1;`
	available := false
	err := r.DB.QueryRow(query, warehouseID).Scan(&available)
	if err != nil {
		logger.ErrLog.Println(err)
		return false, err
	}
	return available, nil
}

func (r *warehouseRepository) AllProducts(warehouseID int) ([]string, error) {
	query := `SELECT (product_code) FROM warehouse_product WHERE warehouse_id = $1;`
	var codes []string
	rows, err := r.DB.Query(query, warehouseID)
	defer rows.Close()

	if err != nil {
		logger.ErrLog.Println(err)
		return nil, err
	}

	for rows.Next() {
		var code string
		err = rows.Scan(&code)
		if err != nil {
			logger.ErrLog.Println(err)
			return nil, err
		}
		codes = append(codes, code)
	}

	return codes, nil
}
