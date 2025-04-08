CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    language VARCHAR(10) DEFAULT 'ru',
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    city_id INTEGER NOT NULL REFERENCES cities(id),
    address_ru TEXT NOT NULL,
    address_kz TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    _2gis_map TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sections (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    UNIQUE (restaurant_id, name)
);

CREATE TABLE IF NOT EXISTS tables (
    id SERIAL PRIMARY KEY,
    number_of_table INTEGER NOT NULL,
    section_id INTEGER NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    qr TEXT,
    UNIQUE (section_id, number_of_table)
);

CREATE TABLE IF NOT EXISTS menu_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    img TEXT
);

CREATE TABLE IF NOT EXISTS menus (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    name_ru VARCHAR(255) NOT NULL,
    name_kz VARCHAR(255),
    img TEXT
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
        CREATE TYPE event_type AS ENUM ('wedding', 'birthday', 'corporate');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS restaurant_events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    eventtype event_type NOT NULL,
    "desc" TEXT,
    price NUMERIC(10,2) NOT NULL,
    img TEXT
);

CREATE TABLE IF NOT EXISTS restaurant_event_tables (
    event_id INTEGER NOT NULL REFERENCES restaurant_events(id) ON DELETE CASCADE,
    table_id INTEGER NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    booking_date TIMESTAMP NOT NULL,
    PRIMARY KEY (event_id, table_id, booking_date)
);

CREATE INDEX IF NOT EXISTS idx_restaurants_city_id ON restaurants(city_id);
CREATE INDEX IF NOT EXISTS idx_sections_restaurant_id ON sections(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_tables_section_id ON tables(section_id);
CREATE INDEX IF NOT EXISTS idx_menus_restaurant_id ON menus(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_restaurant_event_tables_event_id ON restaurant_event_tables(event_id);
CREATE INDEX IF NOT EXISTS idx_restaurant_event_tables_table_id ON restaurant_event_tables(table_id);
CREATE INDEX IF NOT EXISTS idx_restaurant_event_tables_booking_date ON restaurant_event_tables(booking_date);