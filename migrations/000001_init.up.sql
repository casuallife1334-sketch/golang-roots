CREATE TYPE gender AS ENUM ('male', 'female', 'other');
CREATE TYPE relationship_type AS ENUM ('parent_child', 'spouse');
CREATE TYPE relationship_direction AS ENUM ('parent', 'child');

CREATE TABLE persons (
    id CHAR(26) PRIMARY KEY,
    first_name TEXT NOT NULL CHECK (length(trim(first_name)) > 0),
    last_name TEXT NOT NULL CHECK (length(trim(last_name)) > 0),
    birth_date DATE,
    death_date DATE,
    gender gender,
    photo_url TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    CHECK (death_date IS NULL OR birth_date IS NULL OR death_date >= birth_date)
);

CREATE TABLE relationships (
    id CHAR(26) PRIMARY KEY,
    person1_id CHAR(26) NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    person2_id CHAR(26) NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    type relationship_type NOT NULL,
    direction relationship_direction,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (person1_id <> person2_id),
    CHECK ((type = 'parent_child' AND direction IS NOT NULL) OR (type = 'spouse' AND direction IS NULL))
);

CREATE UNIQUE INDEX relationships_unique ON relationships (person1_id, person2_id, type, direction) NULLS NOT DISTINCT;
CREATE INDEX relationships_person1_idx ON relationships(person1_id);
CREATE INDEX relationships_person2_idx ON relationships(person2_id);
