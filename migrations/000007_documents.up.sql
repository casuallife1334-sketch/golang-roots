CREATE TABLE documents (
    id CHAR(26) PRIMARY KEY,
    tree_id CHAR(26) NOT NULL REFERENCES trees(id) ON DELETE CASCADE,
    person_id CHAR(26) REFERENCES persons(id) ON DELETE CASCADE,
    relationship_id CHAR(26) REFERENCES relationships(id) ON DELETE CASCADE,
    file_name TEXT NOT NULL CHECK (length(trim(file_name)) > 0),
    storage_key TEXT NOT NULL UNIQUE,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL CHECK (size_bytes > 0 AND size_bytes <= 20971520),
    created_by CHAR(26) NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (
        (person_id IS NOT NULL AND relationship_id IS NULL)
        OR
        (person_id IS NULL AND relationship_id IS NOT NULL)
    )
);

CREATE INDEX documents_tree_id_idx ON documents(tree_id);
CREATE INDEX documents_person_id_idx ON documents(person_id);
CREATE INDEX documents_relationship_id_idx ON documents(relationship_id);
