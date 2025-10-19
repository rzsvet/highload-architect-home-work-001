#!/bin/bash
set -e

echo "Создание пользователя репликации и настройка мастера..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER $POSTGRES_REPLICATION_USER WITH REPLICATION ENCRYPTED PASSWORD '$POSTGRES_REPLICATION_PASSWORD';
    CREATE PUBLICATION mypublication FOR ALL TABLES;
EOSQL

echo "Настройка синхронной репликации..."

cat >> /var/lib/postgresql/data/pgdata/postgresql.conf <<-EOF
# Настройки репликации
wal_level = replica
max_wal_senders = 5
max_replication_slots = 5
synchronous_commit = on
synchronous_standby_names = 'FIRST 1 (dbreplica01_slot,dbreplica02_slot)'
# listen_addresses = '*'
EOF

cat >> /var/lib/postgresql/data/pgdata/pg_hba.conf <<-EOF
host    replication     all             0.0.0.0/0               md5
EOF

