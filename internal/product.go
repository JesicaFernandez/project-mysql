package internal

import (
	"errors"
	"time"
)

// Product is a struct that contains the attributes of a product
type Product struct {
	// Id is the unique identifier of the product
	Id int
	// Name is the name of the product
	Name string
	// Quantity is the quantity of the product
	Quantity int
	// CodeValue is the code value of the product
	CodeValue string
	// IsPublished is the published status of the product
	IsPublished bool
	// Expiration
	Expiration time.Time
	// Price
	Price float64
	// warehouse id
	WarehouseId int
}

type ProductWithWarehouse struct {
	// Id is the unique identifier of the product
	Id int
	// Name is the name of the product
	Name string
	// Quantity is the quantity of the product
	Quantity int
	// CodeValue is the code value of the product
	CodeValue string
	// IsPublished is the published status of the product
	IsPublished bool
	// Expiration
	Expiration time.Time
	// Price
	Price float64
	// name warehouse
	NameWarehouse string
	// address warehouse
	AddressWarehouse string
	// capacity warehouse
	CapacityWarehouse int
	// telephone warehouse
	TelephoneWarehouse string
}
type ProductRepository interface {
	FindById(id int) (product Product, err error)
	Save(product *Product) (err error)
	Update(product *Product) (err error)
	Delete(id int) (err error)
	GetAll() (products []Product, err error)
	GetProductWithWarehouse(id int) (product ProductWithWarehouse, err error)
}

var (
	ErrProductAlreadyExists = errors.New("product: already exists")
	ErrProductNotFound      = errors.New("product: not found")
)
