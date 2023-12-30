CREATE TYPE status_enum AS ENUM ('Pending', 'In-Process','Completed');

CREATE TABLE "orders" (
  "id" SERIAL PRIMARY KEY,
  "book_id" integer NOT NULL,
  "user_id" INT NOT NULL,
  "order_no" varchar,
  "quantity" int NOT NULL,
  "total_price" int NOT NULL,
  "status" status_enum NOT NULL CHECK (status IN ('Pending', 'In-Process','Completed')),
  "is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
   CONSTRAINT fk_orders_book FOREIGN KEY (book_id) REFERENCES books(id) ON UPDATE CASCADE ON DELETE RESTRICT,
   CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
