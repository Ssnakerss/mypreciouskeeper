/*
 Navicat Premium Data Transfer

 Source Server         : testDB
 Source Server Type    : SQLite
 Source Server Version : 3035005
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3035005
 File Encoding         : 65001

 Date: 04/02/2025 09:23:04
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for mpk_assets
-- ----------------------------
DROP TABLE IF EXISTS "mpk_assets";
CREATE TABLE "mpk_assets" (
  "id" INTEGER NOT NULL,
  "a_user_id" INTEGER NOT NULL,
  "a_type" TEXT NOT NULL,
  "a_sticker" TEXT NOT NULL,
  "a_body" blob NOT NULL,
  "a_created_at" integer,
  "a_updated_at" integer,
  "a_deleted_yn" TEXT,
  "a_deleted_at" integer,
  PRIMARY KEY ("id")
);

PRAGMA foreign_keys = true;
