BEGIN;

UPDATE persons
SET metadata = jsonb_set(metadata, '{patronymic}', to_jsonb(patronymic), true)
WHERE patronymic IS NOT NULL;

ALTER TABLE persons DROP COLUMN patronymic;

COMMIT;
