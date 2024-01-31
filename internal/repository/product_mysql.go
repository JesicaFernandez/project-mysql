package repository

import (
	// sql
	"database/sql"
	// mysql
	"github.com/go-sql-driver/mysql"
	// internal
	"product-testing/internal"
	// errors
	"errors"
)

func NewProductMySQL(db *sql.DB) internal.ProductRepository {
	return &ProductMySQL{db}
}

type ProductMySQL struct {
	// DB is the database connection
	db *sql.DB
}

// find by id (buscar por id)
func (p *ProductMySQL) FindById(id int) (product internal.Product, err error) {

	// query
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE id = ?"

	// execute query
	row := p.db.QueryRow(query, id)
	// check errors
	if row.Err() != nil {
		err = row.Err()
		return
	}

	// serialize product
	err = row.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)
	if err != nil {
		// check if product not found
		if err == sql.ErrNoRows {
			err = internal.ErrProductNotFound
			return
		}
		return
	}
	return
}

func (p *ProductMySQL) Save(product *internal.Product) (err error) {

	// query
	query := "INSERT INTO products (name, quantity, code_value, is_published, expiration, price, id_warehouse) VALUES (?, ?, ?, ?, ?, ?, ?)"

	// execute query
	result, err := p.db.Exec(query, product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price, product.WarehouseId)
	if err != nil {
		// check if product already exists
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok {
			if mysqlErr.Number == 1062 {
				err = internal.ErrProductAlreadyExists
				return
			}
		}
		return
	}

	// get last inserted id
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return
	}
	(*product).Id = int(lastInsertId)

	return
}

func (p *ProductMySQL) Update(product *internal.Product) (err error) {

	// execute query
	_, err = p.db.Exec(
		"UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ? WHERE id = ?",
		product.Name, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.Price, product.Id,
	)
	if err != nil {
		// create mysql error
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			// check if product already exists
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrProductAlreadyExists
				return
			default:
				return
			}
		}
	}
	return
}

func (p *ProductMySQL) Delete(id int) (err error) {

	// execute query
	result, err := p.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return
	}

	// check if product not found
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	// check rows affected
	if rowsAffected == 0 {
		err = internal.ErrProductNotFound
		return
	} else if rowsAffected > 1 {
		err = errors.New("more than one product was deleted")
		return
	}
	return
}

func (p *ProductMySQL) GetAll() (products []internal.Product, err error) {

	// query
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price, id_warehouse FROM products"

	// prepare statement
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return
	}

	defer stmt.Close()

	// execute query
	rows, err := stmt.Query()
	if err != nil {
		return
	}

	defer rows.Close()

	// serialize products
	for rows.Next() {
		var product internal.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price, &product.WarehouseId)
		if err != nil {
			return
		}
		products = append(products, product)
	}
	return
}

func (p *ProductMySQL) GetProductWithWarehouse(id int) (product internal.ProductWithWarehouse, err error) {

	// query
	query := "SELECT p.id, p.name, p.quantity, p.code_value, p.is_published, p.expiration, p.price, w.name, w.address, w.capacity, w.telephone FROM products p INNER JOIN warehouses w ON p.id_warehouse = w.id WHERE p.id = ?"

	// prepare statement
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	// execute query
	row := stmt.QueryRow(id)
	// check errors
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			err = internal.ErrProductNotFound
		} else {
			err = row.Err()
		}
		return
	}

	// serialize product
	err = row.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price, &product.NameWarehouse, &product.AddressWarehouse, &product.CapacityWarehouse, &product.TelephoneWarehouse)
	if err != nil {
		return
	}
	return
}
