CREATE TABLE "banners" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "image" varchar,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "offer_id" int NOT NULL,
 "is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
   CONSTRAINT fk_banner_offer FOREIGN KEY (offer_id) REFERENCES offers(id) ON UPDATE CASCADE ON DELETE RESTRICT
);