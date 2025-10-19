#!/bin/bash
touch /var/lib/postgresql/data/pgdata/standby.signal

cat >> /var/lib/postgresql/data/pgdata/postgresql.conf <<-EOF
# Настройки репликации
primary_conninfo = 'host=db-master port=5432 user=$POSTGRES_REPLICATION_USER password=$POSTGRES_REPLICATION_PASSWORD application_name=dbreplica02_slot'
EOF
