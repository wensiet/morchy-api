BEGIN;

CREATE TABLE node
(
    id     VARCHAR(36) PRIMARY KEY,
    status VARCHAR(12) NOT NULL
);

COMMIT;