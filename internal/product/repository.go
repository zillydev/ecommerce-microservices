package product

import (
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAllProducts() ([]Product, error) {
	rows, err := r.DB.Query("SELECT id, name, category, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *Repository) GetProductByID(id int) (*Product, error) {
	var product Product
	err := r.DB.QueryRow("SELECT id, name, category, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Category, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *Repository) AddProduct(product *NewProductRequest) (*Product, error) {
	var id int
	err := r.DB.QueryRow("INSERT INTO products (name, category, price) VALUES ($1, $2, $3) RETURNING id", product.Name, product.Category, product.Price).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &Product{ID: id, Name: product.Name, Category: product.Category, Price: product.Price}, nil
}
