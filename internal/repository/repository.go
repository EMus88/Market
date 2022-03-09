package repository

import (
	"context"
	"errors"

	"github.com/EMus88/Market/internal/models"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db     DB
	logger *logrus.Logger
}

func NewRepository(db DB, logger *logrus.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) AddCategory(m *models.Category) error {
	var id string
	q := `INSERT INTO categories(category)
 		VALUES($1)
RETURNING id;`
	row := r.db.QueryRow(context.Background(), q, m.Name).Scan(&id)
	if id == "" {
		r.logger.Error(row.Error())
		return errors.New("error: internal db error")
	}
	return nil
}

func (r *Repository) AddProduct(m *models.ProductDTO) error {
	var price uint64
	var id string
	q := `INSERT INTO products(name,weight,valume,description,photo,price,visible,category_id)
 		VALUES($1,$2,$3,$4,$5,$6,$7,
		(SELECT id FROM categories
			WHERE category=$8))
RETURNING id;`
	//convert price to uint64
	price = uint64(m.Price * 100)
	row := r.db.QueryRow(context.Background(), q, m.Name, m.Weight, m.Valume, m.Description, m.Photo, price, m.Visible, m.Category).Scan(&id)
	if id == "" {
		r.logger.Error(row.Error())
		return errors.New("error: internal db error")
	}
	return nil
}

func (r *Repository) ChangeVisible(v *models.Visible) error {
	q := `UPDATE products 
	SET visible=$1
		WHERE name=$2;`
	_, err := r.db.Exec(context.Background(), q, v.Visible, v.Name)
	if err != nil {
		r.logger.Error(err)
		return errors.New("error: internal db error")
	}

	return nil
}

func (r *Repository) GetCatalog() ([]models.ProductDTO, error) {
	var catalog []models.ProductDTO
	q := `SELECT name,weight,valume,description,photo,price,category
	FROM products
	JOIN categories ON category_id=categories.id  
	WHERE visible=true
	ORDER BY name`
	rows, err := r.db.Query(context.Background(), q)
	if err != nil {
		r.logger.Error(err)
		return nil, errors.New("error: internal db error")
	}
	for rows.Next() {
		var product models.ProductDTO
		var price int
		err := rows.Scan(&product.Name, &product.Weight, &product.Valume, &product.Description, &product.Photo, &price, &product.Category)
		product.Price = float64(price) / 100
		if err != nil {
			r.logger.Error(err)
			return nil, errors.New("error: internal db error")
		}
		catalog = append(catalog, product)
	}

	return catalog, nil
}

func (r *Repository) GetByCategory(productName string, category string) ([]models.ProductDTO, error) {
	var result []models.ProductDTO
	q := `SELECT name,weight,valume,description,photo,price
	FROM products
		WHERE name @@ $1 AND category_id=
		(SELECT id FROM categories WHERE category=$2) AND visible=true;`
	rows, err := r.db.Query(context.Background(), q, productName, category)
	if err != nil {
		r.logger.Error(err)
		return nil, errors.New("error: internal db error")
	}
	for rows.Next() {
		var product models.ProductDTO
		var price int
		err := rows.Scan(&product.Name, &product.Weight, &product.Valume, &product.Description, &product.Photo, &price)
		product.Price = float64(price / 100)
		if err != nil {
			r.logger.Error(err)
			return nil, errors.New("error: internal db error")
		}
		result = append(result, product)

	}
	return result, nil
}

func (r *Repository) GetByAllCategories(productName string) ([]models.ProductDTO, error) {
	var result []models.ProductDTO
	q := `SELECT name,weight,valume,description,photo,price
	FROM products
		WHERE name @@ $1 AND visible=true;`
	rows, err := r.db.Query(context.Background(), q, productName)
	if err != nil {
		r.logger.Error(err)
		return nil, errors.New("error: internal db error")
	}
	for rows.Next() {
		var product models.ProductDTO
		var price int
		err := rows.Scan(&product.Name, &product.Weight, &product.Valume, &product.Description, &product.Photo, &price)
		product.Price = float64(price / 100)
		if err != nil {
			r.logger.Error(err)
			return nil, errors.New("error: internal db error")
		}
		result = append(result, product)

	}
	return result, nil
}
