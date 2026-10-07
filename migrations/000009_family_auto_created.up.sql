ALTER TABLE families
    ADD COLUMN auto_created BOOLEAN NOT NULL DEFAULT false;

UPDATE families f
SET auto_created = true
WHERE f.name = ''
  AND f.metadata = '{}'::jsonb
  AND EXISTS (
      SELECT 1
      FROM family_relationships fr
      WHERE fr.family_id = f.id
  );
