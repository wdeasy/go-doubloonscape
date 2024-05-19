ALTER TABLE events
DROP COLUMN amount;
DROP COLUMN chance;
DROP COLUMN max;
DROP COLUMN cooldown;

ALTER TABLE destinations
DROP COLUMN chance;
DROP COLUMN low;
DROP COLUMN max;
DROP COLUMN mod_max;
DROP COLUMN duration;

CREATE TABLE treasure(
  amount INTEGER,
  up BOOLEAN,
  prestige DECIMAL DEFAULT 1 NOT NULL
);