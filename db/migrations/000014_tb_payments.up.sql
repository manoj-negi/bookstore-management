CREATE TYPE payment_status_enum AS ENUM ('Pending', 'In-Process','Completed');

CREATE TABLE "payments" (
  "id" SERIAL PRIMARY KEY,
  "order_id" integer NOT NULL,
  "amount" int NOT NULL,
  "payment_status"  payment_status_enum NOT NULL CHECK (payment_status IN ('Pending', 'In-Process','Completed')),
"is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);