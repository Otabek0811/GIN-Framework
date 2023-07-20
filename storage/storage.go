package storage

import "app/models"

type StorageI interface {
	Close()
	User() UserRepoI
	Category() CategoryRepoI
	Product() ProductRepoI
}

type UserRepoI interface {
	CreateUser(*models.CreateUser) (string, error)
	GetUserByID(*models.UserPrimaryKey) (*models.User, error)
	GetListUser(*models.UserGetListRequest) (*models.UserGetListResponse, error)
	UpdateUser(*models.UpdateUser) (string, error)
	DeleteUser(*models.UserPrimaryKey) error
}

type CategoryRepoI interface {
	CreateCategory(*models.CreateCategory) (string, error)
	GetCategoryByID(*models.CategoryPrimaryKey) (*models.Category, error)
	GetCategoryList(*models.CategoryGetListRequest) (*models.CategoryGetListResponse, error)
	UpdateCategory(*models.UpdateCategory) (string, error)
	DeleteCategory(*models.CategoryPrimaryKey) error
}

type ProductRepoI interface {
	CreateProduct(*models.CreateProduct) (string, error)
	GetListProduct(*models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	GetProductByID(req *models.ProductPrimaryKey) (*models.Product, error)
	UpdateProduct(*models.UpdateProduct) (string, error)
	DeleteProduct(*models.ProductPrimaryKey) error
}
