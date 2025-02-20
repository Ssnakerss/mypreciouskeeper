	CREATE TABLE "mpk_assets" (
	"id" INTEGER NOT NULL PRIMARY KEY,
	"a_user_id" INTEGER NOT NULL,
	"a_type" TEXT NOT NULL,
	"a_sticker" TEXT NOT NULL,
	"a_body" blob NOT NULL,
	"a_created_at" INTEGER,
	"a_updated_at" INTEGER,
	"a_deleted_yn" TEXT,
	"a_deleted_at" INTEGER

	)