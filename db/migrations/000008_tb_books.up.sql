CREATE TABLE "books" (
  "id" SERIAL PRIMARY KEY,
  "title" varchar NOT NULL,
  "author_id" int NOT NULL,
  "publication_date" date NOT NULL,
  "price" integer NOT NULL,
  "stock_quantity" integer NOT NULL,
  "is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT fk_books_author FOREIGN KEY (author_id) REFERENCES authors(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

