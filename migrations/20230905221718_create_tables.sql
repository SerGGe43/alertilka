-- +goose Up

-- +goose StatementBegin
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tracking_type1') THEN
            CREATE TYPE tracking_type1 AS ENUM ('more than price', 'lower than price', 'more than MA');
        END IF;
    END
$$;
CREATE TABLE IF NOT EXISTS "user" (
                                      id SERIAL PRIMARY KEY,
                                      name varchar(32),
                                      chatID bigint
);
CREATE TABLE IF NOT EXISTS alert (
                                     id SERIAL PRIMARY KEY,
                                     ticker varchar(4),
                                     name varchar(32),
                                     userID integer REFERENCES "user"(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS indicator (
                                         id SERIAL PRIMARY KEY,
                                         alertID int REFERENCES alert(id) ON DELETE CASCADE,
                                         indicatorID tracking_type1,
                                         value int
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE indicator;
DROP TABLE alert;
DROP TABLE "user";
DROP TYPE tracking_type1
-- +goose StatementEnd
