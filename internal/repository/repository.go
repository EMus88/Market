package repository

import (
	"JWT_auth/internal/model"
	"context"
	"errors"
	"log"
)

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddCategory(m *model.Category) error {
	var id string
	q := `INSERT INTO categories(category)
 		VALUES($1)
RETURNING id;`
	r.db.QueryRow(context.Background(), q, m.Name).Scan(&id)
	if id == "" {
		return errors.New("error: internal db error")
	}
	return nil
}

func (r *Repository) AddProduct(m *model.ProductDTO) error {
	var id string
	q := `INSERT INTO products(name,weight,valume,description,photo,price,visible,category_id)
 		VALUES($1,$2,$3,$4,$5,$6,$7,
		(SELECT id FROM categories
			WHERE category=$8))
RETURNING id;`
	r.db.QueryRow(context.Background(), q, m.Name, m.Weight, m.Valume, m.Description, m.Photo, m.Price, m.Visible, m.Category).Scan(&id)
	if id == "" {
		return errors.New("error: internal db error")
	}
	return nil
}

func (r *Repository) ChangeVisible(v *model.Visible) error {
	q := `UPDATE products 
	SET visible=$1
		WHERE name=$2;`
	_, err := r.db.Exec(context.Background(), q, v.Visible, v.Name)
	if err != nil {
		return errors.New("error: internal db error")
	}

	return nil
}

func (r *Repository) GetCatalog() ([]model.ProductDTO, error) {
	var catalog []model.ProductDTO
	q := `SELECT name,weight,valume,description,photo,price,category
	FROM products
	JOIN categories ON category_id=categories.id  
	WHERE visible=true
	ORDER BY name`
	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		return nil, errors.New("error: internal db error")
	}
	for rows.Next() {
		var product model.ProductDTO
		var price int
		err := rows.Scan(&product.Name, &product.Weight, &product.Valume, &product.Description, &product.Photo, &price, &product.Category)
		product.Price = float32(price / 100)
		if err != nil {
			log.Println(err)
			return nil, errors.New("error: internal db error")
		}
		catalog = append(catalog, product)
	}

	return catalog, nil
}

func (r *Repository) GetByCategory(productName string, category string) ([]model.ProductDTO, error) {
	var result []model.ProductDTO
	q := `SELECT name,weight,valume,description,photo,price
	FROM products
		WHERE name @@ $1 AND category_id=
		(SELECT id FROM categories WHERE category=$2) AND visible=true;`
	rows, err := r.db.Query(context.Background(), q, productName, category)
	if err != nil {
		return nil, errors.New("error: internal db error")
	}
	for rows.Next() {
		var product model.ProductDTO
		var price int
		err := rows.Scan(&product.Name, &product.Weight, &product.Valume, &product.Description, &product.Photo, &price)
		product.Price = float32(price / 100)
		if err != nil {
			log.Println(err)
			return nil, errors.New("error: internal db error")
		}
		log.Println(product)
		result = append(result, product)

	}
	return result, nil
}

func (r *Repository) GetByAllCategories(productName string) ([]model.ProductDTO, error) {
	var result []model.ProductDTO
	q := `SELECT name,weight,valume,description,photo,price
	FROM products
		WHERE name @@ $1 AND visible=true;`
	rows, err := r.db.Query(context.Background(), q, productName)
	if err != nil {
		return nil, errors.New("error: internal db error")
	}
	for rows.Next() {
		var product model.ProductDTO
		var price int
		err := rows.Scan(&product.Name, &product.Weight, &product.Valume, &product.Description, &product.Photo, &price)
		product.Price = float32(price / 100)
		if err != nil {
			log.Println(err)
			return nil, errors.New("error: internal db error")
		}
		log.Println(product)
		result = append(result, product)

	}
	return result, nil
}
