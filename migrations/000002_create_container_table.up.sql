BEGIN;

CREATE TABLE container
(
    id      VARCHAR(36) PRIMARY KEY,
    node_id VARCHAR(36)  NOT NULL,
    image   VARCHAR(256) NOT NULL,
    status  VARCHAR(12)  NOT NULL,
    FOREIGN KEY (node_id) REFERENCES node (id)
);

CREATE UNIQUE INDEX container__id__node_id__unique ON container (id, node_id);

COMMIT;