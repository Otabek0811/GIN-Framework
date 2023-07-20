package postgres

import (
	"app/models"
	"app/pkg/helper"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) CreateProduct(req *models.CreateProduct) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO product (id, name, price, category_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(query,
		id,
		req.Name,
		req.Price,
		helper.NewNullString(req.Category_id),
	)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (r *productRepo) GetProductByID(req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		resp  models.Product
		query string
	)

	query = `
		SELECT
			id,
			name,
			price,
			COALESCE(category_id::VARCHAR, ''),
			created_at,
			updated_at
		FROM product
		WHERE id = $1
	`

	err := r.db.QueryRow(query, req.Id).Scan(
		&resp.Id,
		&resp.Name,
		&resp.Price,
		&resp.Category_id,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *productRepo) UpdateProduct(req *models.UpdateProduct) (string, error) {
	var (
		id    = req.Id
		query string
	)

	query = `
		Update 
			product 
		set 
			name = $1,
			price = $2,
			category_id= $3,
			updated_at= NOW()
		where id = $4
	`
	_, err := r.db.Exec(query,
		req.Name,
		req.Price,
		helper.NewNullString(req.Category_id),
		id,
	)

	if err != nil {
		return "", err
	}
	return id, nil

}

func (r *productRepo) DeleteProduct(req *models.ProductPrimaryKey) error {
	var (
		id    = req.Id
		query string
	)

	query = `
		DElETE FROM product where id = $1
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepo) GetListProduct(req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)
	query = `
		SELECT
			COUNT(*) OVER(),
		    id,
		    name,
			price,
		    category_id,
		    created_at,
		    updated_at
		FROM product
	`
	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND title ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			product     models.Product
			category_Id sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&product.Id,
			&product.Name,
			&product.Price,
			&category_Id,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		product.Category_id = category_Id.String
		resp.Products = append(resp.Products, &product)
	}
	return resp, nil

}
