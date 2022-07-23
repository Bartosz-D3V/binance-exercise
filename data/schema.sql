PRAGMA foreign_keys= OFF;
BEGIN TRANSACTION;
DROP TABLE IF EXISTS transactions;
CREATE TABLE IF NOT EXISTS transactions
(
    id        INTEGER PRIMARY KEY,
    order_id  INTEGER NOT NULL,
    price     TEXT    NOT NULL,
    quantity  TEXT    NOT NULL,
    timestamp TEXT    NOT NULL
);
COMMIT;
