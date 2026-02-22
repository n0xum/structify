CREATE TABLE "product" (
    "i_d" BIGINT PRIMARY KEY,
    "category_i_d" BIGINT NOT NULL REFERENCES "product_category"("id") ON DELETE RESTRICT,
    "s_k_u" VARCHAR(255) UNIQUE CHECK (length(sku) > 0),
    "name" VARCHAR(255) CHECK (length(name) > 0),
    "description" VARCHAR(255),
    "price" DOUBLE PRECISION NOT NULL CHECK (price >= 0),
    "stock_qty" INTEGER NOT NULL CHECK (stock_qty >= 0) DEFAULT 0,
    "status" VARCHAR(255) CHECK ("status" IN ('active', 'inactive', 'discontinued', 'draft')) DEFAULT 'draft',
    "attributes" JSONB DEFAULT '{}',
    "featured_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL
);

CREATE INDEX "idx_product_name" ON "product" ("name");

CREATE TABLE "product_category" (
    "i_d" BIGINT PRIMARY KEY,
    "name" VARCHAR(255) UNIQUE CHECK (length(name) > 0),
    "slug" VARCHAR(255),
    "description" VARCHAR(255),
    "parent_i_d" BIGINT NOT NULL REFERENCES "product_category"("id") ON DELETE SET NULL,
    "sort_order" INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX "idx_category_slug" ON "product_category" ("slug");

CREATE TABLE "order" (
    "i_d" BIGINT PRIMARY KEY,
    "customer_i_d" BIGINT NOT NULL REFERENCES "customer"("id") ON DELETE CASCADE,
    "status" VARCHAR(255) CHECK ("status" IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'refunded')) DEFAULT 'pending',
    "total_amount" DOUBLE PRECISION NOT NULL CHECK (total_amount >= 0),
    "currency" VARCHAR(255) CHECK (length(currency) = 3) DEFAULT 'USD',
    "created_at" BIGINT NOT NULL DEFAULT extract(epoch from now()),
    "shipped_at" BIGINT NOT NULL,
    "delivered_at" BIGINT NOT NULL
);

CREATE TABLE "customer" (
    "i_d" BIGINT PRIMARY KEY,
    "email" VARCHAR(255) UNIQUE CHECK (email ~* '^[^@]+@[^@]+\\.[^@]+$'),
    "username" VARCHAR(255),
    "full_name" VARCHAR(255) CHECK (length(full_name) > 0),
    "phone" VARCHAR(255),
    "tier" VARCHAR(255) CHECK ("tier" IN ('free', 'silver', 'gold', 'platinum')) DEFAULT 'free',
    "active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX "idx_customer_username" ON "customer" ("username");

CREATE TABLE "order_item" (
    "order_i_d" BIGINT NOT NULL REFERENCES "order"("id") ON DELETE CASCADE,
    "product_i_d" BIGINT NOT NULL REFERENCES "product"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    "quantity" INTEGER NOT NULL CHECK (quantity > 0),
    "unit_price" DOUBLE PRECISION NOT NULL CHECK (unit_price >= 0),
    "discount" DOUBLE PRECISION NOT NULL CHECK (discount >= 0) DEFAULT 0,
    PRIMARY KEY ("order_i_d", "product_i_d")
);

