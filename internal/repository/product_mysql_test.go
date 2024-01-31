package repository_test

import (
	"database/sql"
	"product-testing/internal"
	"product-testing/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// init initializes the database connection configuration and registers the txdb driver.
func init() {
	// database config connection
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "3425216965Yesii",
		Addr:      "localhost:3306",
		Net:       "tcp",
		DBName:    "product_test_db",
		ParseTime: true,
	}

	// register txdb driver
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestProductMySQL_FindByID(t *testing.T) {

	t.Run("success 01 - product was found", func(t *testing.T) {
		// arrange
		// - db
		db, err := sql.Open("txdb", "product_test_db_success")
		require.NoError(t, err)
		defer db.Close() // rollback

		// - set-up
		err = func(db *sql.DB) error {
			_, err := db.Exec(
				"INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`,`id_warehouse` ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				1, "product 01", 10, "code 01", true, time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), 10.0, 0,
			)
			return err
		}(db)
		require.NoError(t, err)
		// - repository
		rp := repository.NewProductMySQL(db)

		// act
		product, err := rp.FindById(1)

		// assert
		expectedProduct := internal.Product{
			Id:          1,
			Name:        "product 01",
			Quantity:    10,
			CodeValue:   "code 01",
			IsPublished: true,
			Expiration:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Price:       10.0,
			WarehouseId: 0,
		}
		require.NoError(t, err)
		require.Equal(t, expectedProduct, product)
	})

	t.Run("failure 01 - product was not found", func(t *testing.T) {
		// arrange
		// - db
		db, err := sql.Open("txdb", "product_test_db_failer")
		require.NoError(t, err)
		defer db.Close() // rollback
		// - repository
		rp := repository.NewProductMySQL(db)

		// act
		product, err := rp.FindById(1)

		// assert
		expectedTask := internal.Product{}
		expectedErr := internal.ErrProductNotFound
		require.Equal(t, expectedTask, product)
		require.ErrorIs(t, err, expectedErr)
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestProductMySQL_Save(t *testing.T) {

	t.Run("success 01 - product was saved", func(t *testing.T) {
		// arrange
		// - db
		db, err := sql.Open("txdb", "product_test_db_success")
		require.NoError(t, err)
		defer db.Close() // rollback

		defer func(db *sql.DB) {
			// rollback
			_, err = db.Exec("DELETE FROM `products`")
			require.NoError(t, err)
			// reset auto increment
			_, err = db.Exec("ALTER TABLE `products` AUTO_INCREMENT = 1")
			require.NoError(t, err)
		}(db)

		// - repository
		rp := repository.NewProductMySQL(db)

		// act
		pr := internal.Product{
			Name:        "product 01",
			Quantity:    10,
			CodeValue:   "code 01",
			IsPublished: true,
			Expiration:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Price:       10.0,
			WarehouseId: 0,
		}

		err = rp.Save(&pr)

		expectedId := 1
		require.NoError(t, err)
		require.Equal(t, expectedId, pr.Id)
	})

	t.Run("failure 01 - product was not saved - duplicate entry", func(t *testing.T) {
		// arrange
		// - db
		db, err := sql.Open("txdb", "product_test_db_failer")
		require.NoError(t, err)
		defer db.Close() // rollback

		defer func(db *sql.DB) {
			// rollback
			_, err = db.Exec("DELETE FROM `products`")
			require.NoError(t, err)
			// reset auto increment
			_, err = db.Exec("ALTER TABLE `products` AUTO_INCREMENT = 1")
			require.NoError(t, err)
		}(db)

		// set-up
		err = func(db *sql.DB) error {
			_, err := db.Exec(
				"INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`,`id_warehouse` ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				1, "product 01", 10, "code 01", true, time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), 10.0, 0,
			)
			return err
		}(db)
		require.NoError(t, err)

		// - repository
		rp := repository.NewProductMySQL(db)

		// act
		pr := internal.Product{
			Name:        "product 01",
			Quantity:    10,
			CodeValue:   "code 01",
			IsPublished: true,
			Expiration:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Price:       10.0,
			WarehouseId: 0,
		}

		err = rp.Save(&pr)

		require.ErrorIs(t, err, internal.ErrProductAlreadyExists)
		require.EqualError(t, err, internal.ErrProductAlreadyExists.Error())
	})
}

func TestProductsMySQL_GetAll(t *testing.T) {

	t.Run("success 01 - empty db", func(t *testing.T) {

		// arrange
		db, err := sql.Open("txdb", "product_test_db_empty")
		assert.NoError(t, err)
		defer db.Close()

		// product repository
		rp := repository.NewProductMySQL(db)

		// act
		p, err := rp.GetAll()

		// assert
		expectedProducts := []internal.Product(nil)
		expectedErr := error(nil)
		assert.Equal(t, expectedProducts, p)
		assert.Equal(t, expectedErr, err)
	})

}
