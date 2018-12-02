create table resources (
  uri text NOT NULL,
  type text NOT NULL,
  hash text NOT NULL,
  data json NOT NULL,
  data_b jsonb NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(uri, type)
)


