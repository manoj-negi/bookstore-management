CREATE TABLE "offers" (
  "id" SERIAL PRIMARY KEY,
  "book_id" integer NOT NULL,
  "discount_percentage" varchar,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
 "is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
   CONSTRAINT fk_offer_book FOREIGN KEY (book_id) REFERENCES books(id) ON UPDATE CASCADE ON DELETE RESTRICT
);