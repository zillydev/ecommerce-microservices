DROP TYPE IF EXISTS notification_type;
CREATE TYPE notification_type AS ENUM (
  'promotions',
  'order_updates',
  'recommendations'
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(150) UNIQUE NOT NULL,
  preferred_notifications notification_type[]
);
