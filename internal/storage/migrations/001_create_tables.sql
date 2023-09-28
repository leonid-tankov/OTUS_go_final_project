-- +goose Up
CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    description VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS slots (
    id SERIAL PRIMARY KEY,
    description VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS social_dem_groups (
    id SERIAL PRIMARY KEY,
    description VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS clicks (
    banner_id INT,
    slot_id INT,
    social_dem_group_id INT,
    count INT NOT NULL,
    FOREIGN KEY (banner_id) REFERENCES banners (id),
    FOREIGN KEY (slot_id) REFERENCES slots (id),
    FOREIGN KEY (social_dem_group_id) REFERENCES social_dem_groups (id)
);
