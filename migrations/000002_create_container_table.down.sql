BEGIN;

DROP INDEX container__id__node_id__unique;
DROP TABLE container;

COMMIT;