ALTER TABLE relationships
    DROP CONSTRAINT relationships_metadata_comment_check,
    DROP CONSTRAINT relationships_metadata_object_check,
    DROP COLUMN updated_at,
    DROP COLUMN metadata;
