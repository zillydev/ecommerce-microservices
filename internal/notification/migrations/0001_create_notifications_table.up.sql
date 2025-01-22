DROP TYPE IF EXISTS notification_type;
CREATE TYPE notification_type AS ENUM (
  'promotions',
  'order_updates',
  'recommendations'
);

CREATE TABLE IF NOT EXISTS notifications (
  id SERIAL PRIMARY KEY,
  userId INTEGER NOT NULL,
  type notification_type NOT NULL,
  content TEXT NOT NULL,
  sentAt TIMESTAMP NOT NULL,
  read BOOLEAN NOT NULL
);
