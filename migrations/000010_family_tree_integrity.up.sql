-- Fail rather than silently reassign members or links if historical data spans trees.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM family_members fm
        JOIN families f ON f.id = fm.family_id
        JOIN persons p ON p.id = fm.person_id
        WHERE p.tree_id IS DISTINCT FROM f.tree_id
    ) OR EXISTS (
        SELECT 1 FROM family_relationships fr
        JOIN families f ON f.id = fr.family_id
        JOIN relationships r ON r.id = fr.relationship_id
        WHERE r.tree_id IS DISTINCT FROM f.tree_id
    ) OR EXISTS (
        SELECT 1 FROM relationships r
        JOIN persons p1 ON p1.id = r.person1_id
        JOIN persons p2 ON p2.id = r.person2_id
        WHERE r.tree_id IS NOT NULL
          AND (p1.tree_id IS DISTINCT FROM r.tree_id OR p2.tree_id IS DISTINCT FROM r.tree_id)
    ) THEN
        RAISE EXCEPTION 'cross-tree family or relationship data must be resolved before migration 000010';
    END IF;
END $$;

-- Migration 000009 guessed provenance from empty name/metadata. That cannot
-- distinguish manually created families. Preserve all existing families;
-- only families created by the application after this migration are managed.
UPDATE families SET auto_created = false WHERE auto_created;

ALTER TABLE family_members ADD COLUMN tree_id CHAR(26);
ALTER TABLE family_relationships ADD COLUMN tree_id CHAR(26);
UPDATE family_members fm SET tree_id = f.tree_id FROM families f WHERE fm.family_id = f.id;
UPDATE family_relationships fr SET tree_id = f.tree_id FROM families f WHERE fr.family_id = f.id;
ALTER TABLE family_members ALTER COLUMN tree_id SET NOT NULL;
ALTER TABLE family_relationships ALTER COLUMN tree_id SET NOT NULL;

ALTER TABLE families ADD CONSTRAINT families_tree_id_id_key UNIQUE (tree_id, id);
ALTER TABLE persons ADD CONSTRAINT persons_tree_id_id_key UNIQUE (tree_id, id);
ALTER TABLE relationships ADD CONSTRAINT relationships_tree_id_id_key UNIQUE (tree_id, id);

-- Keep the original single-column person FKs on relationships for legacy
-- relationships with NULL tree_id; composite FKs enforce scoped records.
ALTER TABLE relationships ADD CONSTRAINT relationships_person1_tree_fk
    FOREIGN KEY (tree_id, person1_id) REFERENCES persons(tree_id, id) ON DELETE CASCADE;
ALTER TABLE relationships ADD CONSTRAINT relationships_person2_tree_fk
    FOREIGN KEY (tree_id, person2_id) REFERENCES persons(tree_id, id) ON DELETE CASCADE;

ALTER TABLE family_members DROP CONSTRAINT family_members_family_id_fkey;
ALTER TABLE family_members DROP CONSTRAINT family_members_person_id_fkey;
ALTER TABLE family_relationships DROP CONSTRAINT family_relationships_family_id_fkey;
ALTER TABLE family_relationships DROP CONSTRAINT family_relationships_relationship_id_fkey;

ALTER TABLE family_members ADD CONSTRAINT family_members_family_tree_fk
    FOREIGN KEY (tree_id, family_id) REFERENCES families(tree_id, id) ON DELETE CASCADE;
ALTER TABLE family_members ADD CONSTRAINT family_members_person_tree_fk
    FOREIGN KEY (tree_id, person_id) REFERENCES persons(tree_id, id) ON DELETE CASCADE;
ALTER TABLE family_relationships ADD CONSTRAINT family_relationships_family_tree_fk
    FOREIGN KEY (tree_id, family_id) REFERENCES families(tree_id, id) ON DELETE CASCADE;
ALTER TABLE family_relationships ADD CONSTRAINT family_relationships_relationship_tree_fk
    FOREIGN KEY (tree_id, relationship_id) REFERENCES relationships(tree_id, id) ON DELETE CASCADE;
