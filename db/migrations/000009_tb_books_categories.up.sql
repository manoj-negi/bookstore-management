CREATE TABLE "books_categories" (
  "id" SERIAL PRIMARY KEY,
  "book_id" integer NOT NULL,
  "category_id" integer NOT NULL,
 "is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT fk_book_category_book FOREIGN KEY (book_id) REFERENCES books(id) ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_book_category_category FOREIGN KEY (category_id) REFERENCES categories(id) ON UPDATE CASCADE ON DELETE RESTRICT
);