CREATE TABLE IF NOT EXISTS categories (
                                          "id" BIGSERIAL PRIMARY KEY,
                                          "name" TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS category_id_categories_idx ON categories(id);

CREATE TABLE IF NOT EXISTS portfolios (
                                        "id" BIGSERIAL PRIMARY KEY,
                                        "profile_id" BIGINT NOT NULL,
                                        "name" TEXT,
                                        "category_id" BIGINT,
                                        "description" TEXT,
                                        FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE INDEX IF NOT EXISTS profile_id_portfolios_idx ON portfolios(profile_id);
CREATE INDEX IF NOT EXISTS portfolio_id_portfolios_idx ON portfolios(id);

CREATE TABLE IF NOT EXISTS crafts (
                                       "id" BIGSERIAL PRIMARY KEY,
                                       "portfolio_id" BIGINT NOT NULL,
                                       "name" TEXT,
                                       "description" TEXT,
                                       FOREIGN KEY (portfolio_id) REFERENCES portfolios(id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS portfolio_id_crafts_idx ON crafts(portfolio_id);
CREATE INDEX IF NOT EXISTS craft_id_crafts_idx ON crafts(id);

CREATE TABLE IF NOT EXISTS tags (
                                       "id" BIGSERIAL PRIMARY KEY,
                                       "name" TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS tag_id_tags_idx ON tags(id);

CREATE TABLE IF NOT EXISTS crafts_tags (
                                           "craft_id" BIGINT NOT NULL,
                                           "tag_id" BIGINT NOT NULL,
                                           PRIMARY KEY (craft_id, tag_id),
                                           FOREIGN KEY (craft_id) REFERENCES crafts(id) ON DELETE CASCADE,
                                           FOREIGN KEY (tag_id) REFERENCES tags(id)
);

CREATE INDEX IF NOT EXISTS craft_id_crafts_tags_idx ON crafts_tags(craft_id);

CREATE TABLE IF NOT EXISTS contents (
                                        "id" BIGSERIAL PRIMARY KEY,
                                        "craft_id" BIGINT NOT NULL,
                                        "description" TEXT,
                                        "data" BYTEA NOT NULL,
                                        FOREIGN KEY (craft_id) REFERENCES crafts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS craft_id_contents_idx ON contents(craft_id);