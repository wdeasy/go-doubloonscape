CREATE TABLE destinations(
  name character varying NOT NULL,
  end_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  amount INTEGER DEFAULT 0 NOT NULL
)