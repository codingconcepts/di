CREATE TABLE "example" (
  "uuid" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "string" STRING NOT NULL,
  "date" DATE NOT NULL,
  "time" TIME NOT NULL,
  "timestamptz" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "int2" INT2 NOT NULL,
  "int4" INT4 NOT NULL,
  "int8" INT8 NOT NULL,
  "bool" BOOL NOT NULL,
  "float" FLOAT NULL,
  "decimal" DECIMAL NULL
);