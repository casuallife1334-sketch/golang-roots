CREATE TYPE tree_role AS ENUM ('owner', 'editor', 'viewer');

CREATE TABLE trees (
    id CHAR(26) PRIMARY KEY,
    owner_id CHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tree_members (
    tree_id CHAR(26) NOT NULL REFERENCES trees(id) ON DELETE CASCADE,
    user_id CHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role tree_role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tree_id, user_id)
);

CREATE INDEX trees_owner_id_idx ON trees(owner_id);
CREATE INDEX tree_members_user_id_idx ON tree_members(user_id);

ALTER TABLE persons ADD COLUMN tree_id CHAR(26) REFERENCES trees(id) ON DELETE CASCADE;
ALTER TABLE relationships ADD COLUMN tree_id CHAR(26) REFERENCES trees(id) ON DELETE CASCADE;

CREATE INDEX persons_tree_id_idx ON persons(tree_id);
CREATE INDEX relationships_tree_id_idx ON relationships(tree_id);
