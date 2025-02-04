/*
 Navicat Premium Data Transfer

 Source Server         : testDB
 Source Server Type    : SQLite
 Source Server Version : 3035005
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3035005
 File Encoding         : 65001

 Date: 04/02/2025 09:23:11
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for mpk_users
-- ----------------------------
DROP TABLE IF EXISTS "mpk_users";
CREATE TABLE "mpk_users" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "u_email" text NOT NULL,
  "u_pass_hash" text NOT NULL,
  "u_created_at" integer NOT NULL,
  "u_updated_at" integer NOT NULL
);

-- ----------------------------
-- Auto increment value for mpk_users
-- ----------------------------
UPDATE "sqlite_sequence" SET seq = 4 WHERE name = 'mpk_users';

PRAGMA foreign_keys = true;
