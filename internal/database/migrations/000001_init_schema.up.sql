CREATE TABLE "cats" (
"id" BIGSERIAL  PRIMARY KEY,
"name" VARCHAR NOT NULL,
"years_of_experience" SMALLINT NOT NULL,
"breed" VARCHAR NOT NULL,
"salary" DECIMAL NOT NULL,
"created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE "missions" (
"id" BIGSERIAL PRIMARY KEY,
"name" VARCHAR NOT NULL,
"cat_id" BIGINT DEFAULT NULL,
"is_completed" BOOLEAN NOT NULL DEFAULT FALSE,
"created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);


CREATE TABLE "targets" (
"id" BIGSERIAL PRIMARY KEY,
"mission_id" BIGINT NOT NULL,
"name" VARCHAR NOT NULL,
"country" VARCHAR NOT NULL,
"notes" VARCHAR NOT NULL,
"is_completed" BOOLEAN NOT NULL DEFAULT FALSE,
"created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);


ALTER TABLE "missions" ADD FOREIGN KEY ("cat_id") REFERENCES "cats" ("id");

ALTER TABLE "targets" ADD FOREIGN KEY ("mission_id") REFERENCES "missions" ("id") ON DELETE CASCADE;



