-- ----------------------------
-- Table structure for mpk_users
-- ----------------------------
CREATE TABLE IF NOT EXISTS "mpk_users" (
	"id" INTEGER NOT NULL PRIMARY KEY ,
	"u_local_id" INTEGER NOT NULL UNIQUE,
	"u_email" TEXT NOT NULL UNIQUE,
	"u_pass_hash" TEXT NOT NULL,
	"u_created_at" INTEGER NOT NULL,
	"u_updated_at" INTEGER NOT NULL
	)