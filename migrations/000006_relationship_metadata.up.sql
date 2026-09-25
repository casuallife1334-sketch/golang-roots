ALTER TABLE relationships
    ADD COLUMN metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN updated_at TIMESTAMPTZ,
    ADD CONSTRAINT relationships_metadata_object_check
        CHECK (jsonb_typeof(metadata) = 'object'),
    ADD CONSTRAINT relationships_metadata_comment_check
        CHECK (
            NOT (metadata ? 'comment')
            OR (
                jsonb_typeof(metadata->'comment') = 'string'
                AND char_length(metadata->>'comment') <= 5000
            )
        );
