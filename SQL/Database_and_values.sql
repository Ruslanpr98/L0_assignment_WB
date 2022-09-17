DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
  order_uid char(19) NOT NULL PRIMARY KEY,
  track_number varchar(50) NOT NULL,
  entry_name varchar(50) NOT NULL,
  locale varchar(20) NOT NULL,
  internal_signature varchar(25) NOT NULL,
  customer_id varchar(100) NOT NULL,
  delivery_service varchar(50) NOT NULL,
  shardkey varchar(50) NOT NULL,
  sm_id INT NOT NULL,
  date_created varchar(50) NOT NULL,
  oof_shard varchar(50) NOT NULL
);

INSERT INTO orders (order_uid, track_number, entry_name, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES
('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '', 'test', 'meest', '9', 99, '2021-11-26 06:22:19', '1'),
('a547mar1b7b83b9test', 'WBILMTESTTRACK2', 'WBTOL', 'ru', '', 'test345', 'meesto7', '2', 65, '2021-12-26 12:23:25', '2'),
('n921dec0b4b42b8test', 'WBILMTESTTRACK267', 'WBTOLIN', 'ge', '', 'test102', 'meestok', '3', 45, '2021-10-15 21:45:35', '3');