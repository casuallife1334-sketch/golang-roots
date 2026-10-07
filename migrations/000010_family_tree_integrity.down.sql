ALTER TABLE family_members DROP CONSTRAINT family_members_family_tree_fk;
ALTER TABLE family_members DROP CONSTRAINT family_members_person_tree_fk;
ALTER TABLE family_relationships DROP CONSTRAINT family_relationships_family_tree_fk;
ALTER TABLE family_relationships DROP CONSTRAINT family_relationships_relationship_tree_fk;
ALTER TABLE relationships DROP CONSTRAINT relationships_person1_tree_fk;
ALTER TABLE relationships DROP CONSTRAINT relationships_person2_tree_fk;

ALTER TABLE family_members ADD CONSTRAINT family_members_family_id_fkey
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;
ALTER TABLE family_members ADD CONSTRAINT family_members_person_id_fkey
    FOREIGN KEY (person_id) REFERENCES persons(id) ON DELETE CASCADE;
ALTER TABLE family_relationships ADD CONSTRAINT family_relationships_family_id_fkey
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;
ALTER TABLE family_relationships ADD CONSTRAINT family_relationships_relationship_id_fkey
    FOREIGN KEY (relationship_id) REFERENCES relationships(id) ON DELETE CASCADE;

ALTER TABLE family_members DROP COLUMN tree_id;
ALTER TABLE family_relationships DROP COLUMN tree_id;
ALTER TABLE families DROP CONSTRAINT families_tree_id_id_key;
ALTER TABLE relationships DROP CONSTRAINT relationships_tree_id_id_key;
ALTER TABLE persons DROP CONSTRAINT persons_tree_id_id_key;

-- Previous auto_created guesses cannot be safely restored on rollback.
