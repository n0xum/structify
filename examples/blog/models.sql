CREATE TABLE "author" (
    "i_d" BIGINT PRIMARY KEY,
    "email" VARCHAR(255) UNIQUE CHECK (email ~* '^[^@]+@[^@]+\\.[^@]+$'),
    "username" VARCHAR(255) CHECK (length(username) >= 3),
    "full_name" VARCHAR(255) CHECK (length(full_name) > 0),
    "bio" VARCHAR(255),
    "avatar_u_r_l" VARCHAR(255),
    "role" VARCHAR(255) CHECK ("role" IN ('reader', 'author', 'editor', 'admin')) DEFAULT 'author',
    "metadata" JSONB DEFAULT '{}',
    "active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX "idx_author_username" ON "author" ("username");

CREATE TABLE "post" (
    "i_d" BIGINT PRIMARY KEY,
    "author_i_d" BIGINT NOT NULL REFERENCES "author"("id") ON DELETE CASCADE,
    "category_i_d" BIGINT NOT NULL REFERENCES "category"("id") ON DELETE SET NULL,
    "title" VARCHAR(255) CHECK (length(title) >= 5),
    "slug" VARCHAR(255) UNIQUE CHECK (length(slug) > 0),
    "content" VARCHAR(255) CHECK (length(content) > 0),
    "summary" VARCHAR(255),
    "status" VARCHAR(255) CHECK ("status" IN ('draft', 'review', 'published', 'archived')) DEFAULT 'draft',
    "view_count" INTEGER NOT NULL CHECK (view_count >= 0) DEFAULT 0,
    "featured" BOOLEAN NOT NULL DEFAULT false,
    "published_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL
);

CREATE TABLE "category" (
    "i_d" BIGINT PRIMARY KEY,
    "name" VARCHAR(255) UNIQUE CHECK (length(name) > 0),
    "slug" VARCHAR(255) CHECK (length(slug) > 0),
    "description" VARCHAR(255),
    "color" VARCHAR(255) DEFAULT '#6b7280',
    "sort_order" INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX "idx_category_slug" ON "category" ("slug");

CREATE TABLE "tag" (
    "i_d" BIGINT PRIMARY KEY,
    "name" VARCHAR(255) UNIQUE CHECK (length(name) > 0),
    "slug" VARCHAR(255) UNIQUE
);

CREATE TABLE "post_tag" (
    "post_i_d" BIGINT NOT NULL REFERENCES "post"("id") ON DELETE CASCADE,
    "tag_i_d" BIGINT NOT NULL REFERENCES "tag"("id") ON DELETE CASCADE,
    "sort_order" INTEGER NOT NULL CHECK (sort_order >= 0) DEFAULT 0,
    PRIMARY KEY ("post_i_d", "tag_i_d")
);

CREATE TABLE "comment" (
    "i_d" BIGINT PRIMARY KEY,
    "post_i_d" BIGINT NOT NULL REFERENCES "post"("id") ON DELETE CASCADE,
    "author_i_d" BIGINT NOT NULL REFERENCES "author"("id") ON DELETE SET NULL,
    "parent_i_d" BIGINT NOT NULL REFERENCES "comment"("id") ON DELETE CASCADE,
    "content" VARCHAR(255) CHECK (length(content) > 0),
    "approved" BOOLEAN NOT NULL DEFAULT false,
    "likes" INTEGER NOT NULL CHECK (likes >= 0) DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT now()
);

