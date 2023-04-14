# Prometheus Postgres Logical Replication Metrics Exporter

Simple prometheus exporter for monitoring PostgreSQL logical replication.

## Building and running

```shell
$ make build
```

```shell
$ ./bin/export --help
Usage of ./bin/exporter:
  -listen-address string
    	The address to listen on for HTTP requests. (default ":9394")
  -log-level string
    	Level of the logs. (default "info")
  -primary-uri string
    	Connection URI of the primary instance host.
  -standby-uri string
    	Connection URI of the standby instance host.
```

### Postgres role setup

Both roles (on primary and standby instances) require superuser priveleges as they access system tables and functions, such as `pg_stat_subscription`, `pg_stat_replication`, `pg_replication_slots` tables and `pg_current_wal_lsn()` function.

## Available metrics

| Metric                                         | Type  | Description                                                                             |
|------------------------------------------------|-------|-----------------------------------------------------------------------------------------|
| pg_logical_replication_subscription_status     | Guage | Status of subscription.                                                                 |
| pg_logical_replication_subscription_lag_bytes  | Guage | The amount of WAL records generated in the primary, but not yet applied in the standby. |
| pg_logical_replication_publication_status      | Guage | Status of publication.                                                                  |
| pg_logical_replication_publication_lag         | Guage | The amount of WAL records generated in the primary, but not yet sent to the standby.    |
| pg_logical_replication_replication_slot_status | Guage | Status of replication slot.                                                             |

## Local development

This repo contains a compose file which facilitates local setup of postgres primary and standby instance. To run instances you will need either `docker`, `podman-compose` or `nerdctl`. To run compose execute following commands:

```shell
$ make start-db
```

To load seed data to databases run following command:

```shell
$ make seed-db
```

To stop database instances run:

```shell
$ make stop-db
```

To run the exporter locally using these databases run:

```shell
$ ./bin/exporter --primary-uri="postgres://primary:primary@localhost:9432/primary?sslmode=disable" --standby-uri="postgres://standby:standby@localhost:9442/standby?sslmode=disable"
```
