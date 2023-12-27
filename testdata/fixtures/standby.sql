CREATE ROLE exporter WITH LOGIN PASSWORD '1234';

CREATE SUBSCRIPTION sub1
  CONNECTION 'host=primary port=5432 user=primary password=primary dbname=primary'
PUBLICATION pub1
WITH (connect = true, create_slot = true, enabled = true);

CREATE SUBSCRIPTION sub2
  CONNECTION 'host=primary port=5432 user=primary password=primary dbname=primary'
PUBLICATION pub2
WITH (connect = true, create_slot = true, enabled = true);
