ALTER TABLE captains ALTER COLUMN gold TYPE BIGINT;

ALTER TABLE captains
ADD COLUMN prestige DECIMAL DEFAULT 1 NOT NULL;