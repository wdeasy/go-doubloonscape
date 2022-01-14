CREATE TABLE treasure(
  amount INTEGER,
  up BOOLEAN
);

CREATE TABLE events(
  name character varying NOT NULL,
  last TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  up BOOLEAN  
);

CREATE TABLE logs(
  text character varying NOT NULL,  
  time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);