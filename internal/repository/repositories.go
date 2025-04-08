package repository

import (
	"context"
	"time"

	"restaurant-management/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type CityRepository interface {
	Create(ctx context.Context, city *models.City) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.City, error)
	GetByName(ctx context.Context, name string) (*models.City, error)
	Update(ctx context.Context, city *models.City) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.City, error)
}

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *models.Restaurant) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Restaurant, error)
	GetByCity(ctx context.Context, cityID int64) ([]*models.Restaurant, error)
	Update(ctx context.Context, restaurant *models.Restaurant) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, active bool) ([]*models.Restaurant, error)
}

type SectionRepository interface {
	Create(ctx context.Context, section *models.Section) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Section, error)
	GetByRestaurant(ctx context.Context, restaurantID int64) ([]*models.Section, error)
	Update(ctx context.Context, section *models.Section) error
	Delete(ctx context.Context, id int64) error
}

type TableRepository interface {
	Create(ctx context.Context, table *models.Table) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Table, error)
	GetBySection(ctx context.Context, sectionID int64) ([]*models.Table, error)
	Update(ctx context.Context, table *models.Table) error
	Delete(ctx context.Context, id int64) error
	GenerateQR(ctx context.Context, tableID int64) (string, error)
}

type MenuTypeRepository interface {
	Create(ctx context.Context, menuType *models.MenuType) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.MenuType, error)
	Update(ctx context.Context, menuType *models.MenuType) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.MenuType, error)
}

type MenuRepository interface {
	Create(ctx context.Context, menu *models.Menu) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Menu, error)
	GetByRestaurant(ctx context.Context, restaurantID int64) ([]*models.Menu, error)
	Update(ctx context.Context, menu *models.Menu) error
	Delete(ctx context.Context, id int64) error
}

type RestaurantEventRepository interface {
	Create(ctx context.Context, event *models.RestaurantEvent) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.RestaurantEvent, error)
	GetByType(ctx context.Context, eventType models.EventType) ([]*models.RestaurantEvent, error)
	Update(ctx context.Context, event *models.RestaurantEvent) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.RestaurantEvent, error)
}

type RestaurantEventTableRepository interface {
	Create(ctx context.Context, eventTable *models.RestaurantEventTable) error
	GetByEvent(ctx context.Context, eventID int64) ([]*models.RestaurantEventTable, error)
	GetByTable(ctx context.Context, tableID int64) ([]*models.RestaurantEventTable, error)
	Delete(ctx context.Context, eventID, tableID int64, date time.Time) error
	CheckAvailability(ctx context.Context, tableID int64, date time.Time) (bool, error)
}

type Repository struct {
	User                 UserRepository
	City                 CityRepository
	Restaurant           RestaurantRepository
	Section              SectionRepository
	Table                TableRepository
	MenuType             MenuTypeRepository
	Menu                 MenuRepository
	RestaurantEvent      RestaurantEventRepository
	RestaurantEventTable RestaurantEventTableRepository
}
