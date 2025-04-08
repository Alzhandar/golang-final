package models

import "time"

type User struct {
	ID          int64  `json:"id" db:"id"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Name        string `json:"name" db:"name"`
	LastName    string `json:"last_name" db:"last_name"`
	Language    string `json:"language" db:"language"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

type City struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Restaurant struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CityID    int64  `json:"city_id" db:"city_id"`
	AddressRU string `json:"address_ru" db:"address_ru"`
	AddressKZ string `json:"address_kz" db:"address_kz"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	Map2GIS   string `json:"_2gis_map" db:"_2gis_map"`
}

type Section struct {
	ID           int64  `json:"id" db:"id"`
	RestaurantID int64  `json:"restaurant_id" db:"restaurant_id"`
	Name         string `json:"name" db:"name"`
}

type Table struct {
	ID            int64  `json:"id" db:"id"`
	NumberOfTable int    `json:"number_of_table" db:"number_of_table"`
	SectionID     int64  `json:"section_id" db:"section_id"`
	QR            string `json:"qr" db:"qr"`
}

type MenuType struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Img  string `json:"img" db:"img"`
}

type Menu struct {
	ID           int64  `json:"id" db:"id"`
	RestaurantID int64  `json:"restaurant_id" db:"restaurant_id"`
	NameRU       string `json:"name_ru" db:"name_ru"`
	NameKZ       string `json:"name_kz" db:"name_kz"`
	Img          string `json:"img" db:"img"`
}

type EventType string

const (
	EventTypeWedding   EventType = "wedding"
	EventTypeBirthday  EventType = "birthday"
	EventTypeCorporate EventType = "corporate"
)

type RestaurantEvent struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	EventType   EventType `json:"eventtype" db:"eventtype"`
	Description string    `json:"desc" db:"desc"`
	Price       float64   `json:"price" db:"price"`
	Img         string    `json:"img" db:"img"`
}

type RestaurantEventTable struct {
	EventID     int64     `json:"event_id" db:"event_id"`
	TableID     int64     `json:"table_id" db:"table_id"`
	BookingDate time.Time `json:"booking_date" db:"booking_date"`
}
