CREATE DATABASE store
  PRIMARY REGION 'aws-us-east-1'
  REGIONS 'aws-eu-central-1', 'aws-ap-southeast-1';

USE store;

CREATE TABLE customer (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" STRING NOT NULL,

  UNIQUE("email")
) LOCALITY REGIONAL BY ROW;