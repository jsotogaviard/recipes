-- +migrate Up
-- Create user table
CREATE TABLE IF NOT EXISTS t_user (
  unique_id bigserial primary key,
  login text NOT NULL,
  hashed_password text NOT NULL
);

-- Create recipe table
CREATE TABLE IF NOT EXISTS t_recipe (
  unique_id bigserial primary key,
  name text NOT NULL,
  preparation_time integer NOT NULL,
  difficulty integer NOT NULL CONSTRAINT valid_difficulty CHECK (difficulty > 0 and difficulty <= 3),
  vegetarian boolean NOT NULL,
  created_by INTEGER NOT NULL REFERENCES t_user (unique_id),
  updated_by INTEGER REFERENCES t_user (unique_id)
);

-- Create recipe rating table
CREATE TABLE IF NOT EXISTS t_recipe_rating (
  unique_id bigserial primary key,
  recipe_unique_id INTEGER NOT NULL REFERENCES t_recipe (unique_id) ON DELETE CASCADE,
  rating integer NOT NULL CONSTRAINT valid_rating CHECK (rating > 0 and rating <= 5)
);

-- Add indexes to ease the search
CREATE INDEX vegetarian_idx ON t_recipe (vegetarian);
CREATE INDEX difficulty_idx ON t_recipe (difficulty);

-- +migrate Down
DROP TABLE IF EXISTS t_recipe_rating;
DROP TABLE IF EXISTS t_recipe;
DROP TABLE IF EXISTS t_user ;
