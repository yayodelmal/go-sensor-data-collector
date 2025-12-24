-- 1. sht31_readings table + hypertable
CREATE TABLE IF NOT EXISTS sht31_readings (
  id          SERIAL       PRIMARY KEY,
  temperature DOUBLE PRECISION NOT NULL,
  humidity    DOUBLE PRECISION NOT NULL,
  ts          TIMESTAMPTZ  NOT NULL DEFAULT now()
);
SELECT create_hypertable('sht31_readings', 'ts', if_not_exists => TRUE);

-- 2. ds18b20_readings table + hypertable
CREATE TABLE IF NOT EXISTS ds18b20_readings (
  id          SERIAL       PRIMARY KEY,
  temperature DOUBLE PRECISION NOT NULL,
  ts          TIMESTAMPTZ  NOT NULL DEFAULT now()
);
SELECT create_hypertable('ds18b20_readings', 'ts', if_not_exists => TRUE);
