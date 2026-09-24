BEGIN;

ALTER TABLE persons
    ADD COLUMN patronymic TEXT
    CHECK (patronymic IS NULL OR length(trim(patronymic)) BETWEEN 1 AND 200);

UPDATE persons
SET patronymic = NULLIF(trim(metadata->>'patronymic'), '')
WHERE jsonb_typeof(metadata->'patronymic') = 'string';

UPDATE persons
SET metadata = metadata - 'patronymic'
WHERE jsonb_typeof(metadata->'patronymic') IN ('string', 'null');

COMMIT;
