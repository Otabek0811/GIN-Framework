package postgres

import (
	"app/models"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(req *models.CreateUser) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO users (id, first_name, last_name, balans, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(query,
		id,
		req.FirstName,
		req.LastName,
		req.Balans,
	)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (r *userRepo) GetUserByID(req *models.UserPrimaryKey) (*models.User, error) {
	var (
		resp  models.User
		query string
	)

	query = `
		SELECT
			id,
			first_name,
			last_name,
			balans,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, req.Id).Scan(
		&resp.Id,
		&resp.FirstName,
		&resp.LastName,
		&resp.Balans,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *userRepo) UpdateUser(req *models.UpdateUser) (string, error) {
	var (
		id    = req.Id
		query string
	)

	query = `
		Update 
			users 
		set 
			first_name = $1,
			last_name = $2,
			balans= $3,
			updated_at= NOW()
		where id = $4
	`
	_, err := r.db.Exec(query,
		req.FirstName,
		req.LastName,
		req.Balans,
		id,
	)

	if err != nil {
		return "", err
	}
	return id, nil

}

func (r *userRepo) DeleteUser(req *models.UserPrimaryKey) error {
	var (
		id    = req.Id
		query string
	)

	query = `
		DElETE FROM users where id = $1
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetListUser(req *models.UserGetListRequest) (*models.UserGetListResponse, error) {
	var (
		resp   = &models.UserGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)
	query = `
		SELECT
			COUNT(*) OVER(),
		    id,
		    first_name,
			last_name,
		    balans,
		    created_at,
		    updated_at
		FROM users
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
			user models.User
		)

		err := rows.Scan(
			&resp.Count,
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Balans,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, &user)
	}
	return resp, nil

}
