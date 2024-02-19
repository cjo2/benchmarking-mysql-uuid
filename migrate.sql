CREATE DATABASE IF NOT EXISTS TestDb;

USE TestDb;

DROP TABLE IF EXISTS TestTable;

CREATE TABLE TestTable (
   ID varchar(36) NOT NULL,
   Name varchar(10),
   PRIMARY KEY (ID)
);