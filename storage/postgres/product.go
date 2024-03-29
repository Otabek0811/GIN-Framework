package postgres

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(ctx context.Context, req *models.CreateProduct) (string, error) {

	trx, err := r.db.Begin(ctx)
	if err != nil {
		return "", nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO product(id, name, price, category_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err = trx.Exec(ctx, query,
		id,
		req.Name,
		req.Price,
		helper.NewNullString(req.CategoryId),
	)

	if err != nil {
		return "", err
	}

	// Market to Product Relation -> Many to Many
	if len(req.MarketIds) > 0 {
		marketProductInsertQuery := `
			INSERT INTO 
				market_product_relation(product_id, market_id) 
			VALUES`

		marketProductInsertQuery, args := helper.InsertMultiple(marketProductInsertQuery, id, req.MarketIds)
		_, err = trx.Exec(ctx, marketProductInsertQuery, args...)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (r *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		query string

		id         sql.NullString
		name       sql.NullString
		price      sql.NullFloat64
		categoryID sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString

		marketObj pgtype.JSONB
	)

	query = `
		WITH market_product AS (
			SELECT
				JSON_AGG(
					JSON_BUILD_OBJECT (
						'id', m.id,
						'name', m.name,
						'address', m.address,
						'phone_number', m.phone_number,
						'created_at', m.created_at,
						'updated_at', m.updated_at
					)
				)  AS markets,
				mpr.product_id AS product_id

			FROM market AS m
			JOIN market_product_relation AS mpr ON mpr.market_id = m.id
			WHERE mpr.product_id = $1
			GROUP BY mpr.product_id
		)
		SELECT
			p.id,
			p.name,
			p.price,
			p.category_id,
			p.created_at,
			p.updated_at,

			mp.markets
			
		FROM product AS p
		LEFT JOIN market_product AS mp ON mp.product_id = p.id
		WHERE p.id =  $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&price,
		&categoryID,
		&createdAt,
		&updatedAt,
		&marketObj,
	)

	if err != nil {
		return nil, err
	}

	markets := []*models.Market{}
	marketObj.AssignTo(&markets)

	return &models.Product{
		Id:         id.String,
		Name:       name.String,
		Price:      price.Float64,
		CategoryId: categoryID.String,
		CreatedAt:  createdAt.String,
		UpdatedAt:  updatedAt.String,
		Markets:    markets,
	}, nil
}

func (r *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {

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
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			price      sql.NullFloat64
			categoryID sql.NullString
			createdAt  sql.NullString
			updatedAt  sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&price,
			&categoryID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, &models.Product{
			Id:         id.String,
			Name:       name.String,
			Price:      price.Float64,
			CategoryId: categoryID.String,
			CreatedAt:  createdAt.String,
			UpdatedAt:  updatedAt.String,
		})
	}

	return resp, nil
}

func (r *productRepo) Update(ctx context.Context, req *models.UpdateProduct) (int64, error) {

	trx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			product
		SET
			name = :name,
			price = :price,
			category_id = :category_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":          req.Id,
		"name":        req.Name,
		"price":       req.Price,
		"category_id": req.CategoryId,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := trx.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	// Market to Product Relation -> Many to Many
	if len(req.MarketIds) > 0 {
		var count int
		marketProductRelationCountQuery := `
			SELECT COUNT(*) FROM market_product_relation WHERE market_id = $1 
		`

		err := trx.QueryRow(ctx, marketProductRelationCountQuery, req.Id).Scan(&count)
		if err != nil {
			return 0, err
		}

		if count > 0 {
			marketProductRelationDeleteQuery := `
				DELETE FROM market_product_relation WHERE market_id = $1 
			`

			_, err := trx.Exec(ctx, marketProductRelationDeleteQuery, req.Id)
			if err != nil {
				return 0, err
			}
		}

		marketProductInsertQuery := `
				INSERT INTO 
					market_product_relation(product_id, market_id) 
				VALUES`

		marketProductInsertQuery, args := helper.InsertMultiple(marketProductInsertQuery, req.Id, req.MarketIds)
		_, err = trx.Exec(ctx, marketProductInsertQuery, args...)
		if err != nil {
			return 0, err
		}
	}

	return result.RowsAffected(), nil
}

func (r *productRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

	var (
		query string
		set   string
	)

	if len(req.Fields) <= 0 {
		return 0, errors.New("no fields")
	}

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s, ", key, key)
	}

	query = `
		UPDATE
			product
		SET ` + set + ` updated_at = now()
		WHERE id = :id
	`

	req.Fields["id"] = req.ID

	fmt.Println(query)

	query, args := helper.ReplaceQueryParams(query, req.Fields)
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *productRepo) Delete(ctx context.Context, req *models.ProductPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM market_product_relation WHERE product_id = $1", req.Id)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, "DELETE FROM product WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
