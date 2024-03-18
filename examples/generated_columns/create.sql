CREATE DATABASE example
  PRIMARY REGION 'aws-us-east-1'
  REGIONS 'aws-eu-central-1', 'aws-ap-southeast-1';

USE example;

SET enable_super_regions = 'on';
ALTER DATABASE example ADD SUPER REGION "us" VALUES 'aws-us-east-1';
ALTER DATABASE example ADD SUPER REGION "eu" VALUES 'aws-eu-central-1';
ALTER DATABASE example ADD SUPER REGION "ap" VALUES 'aws-ap-southeast-1';

CREATE TABLE "example" (
  "uuid" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "market" STRING NOT NULL,
  "crdb_region" crdb_internal_region AS (
    CASE
      WHEN "market" IN ('de', 'es', 'fr', 'ie', 'uk') THEN 'aws-eu-central-1'
      WHEN "market" IN ('br', 'bs', 'co', 'mx', 'us') THEN 'aws-us-east-1'
      WHEN "market" IN ('hk', 'in', 'ja', 'my', 'sg') THEN 'aws-ap-southeast-1'
      ELSE 'aws-us-east-1'
    END
  ) STORED
) LOCALITY REGIONAL BY ROW;