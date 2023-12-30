CREATE TABLE "countries" (
  "id" SERIAL PRIMARY KEY,
  "iso2" varchar NOT NULL,
  "short_name" varchar NOT NULL,
  "long_name" varchar NOT NULL,
  "numcode" varchar,
  "calling_code" varchar NOT NULL,
  "cctld" varchar NOT NULL,
  "is_deleted" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);




