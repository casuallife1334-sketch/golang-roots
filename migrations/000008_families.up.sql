CREATE TYPE family_member_role AS ENUM ('partner', 'parent', 'child');

CREATE TABLE families (
    id CHAR(26) PRIMARY KEY,
    tree_id CHAR(26) NOT NULL REFERENCES trees(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT families_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE TABLE family_members (
    family_id CHAR(26) NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    person_id CHAR(26) NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    role family_member_role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (family_id, person_id, role)
);

CREATE TABLE family_relationships (
    family_id CHAR(26) NOT NULL REFERENCES families(id) ON DELETE CASCADE,
    relationship_id CHAR(26) NOT NULL REFERENCES relationships(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (family_id, relationship_id)
);

CREATE INDEX families_tree_id_idx ON families(tree_id);
CREATE INDEX family_members_person_id_idx ON family_members(person_id);
CREATE INDEX family_relationships_relationship_id_idx ON family_relationships(relationship_id);
