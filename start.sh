set -e 

echo "run db migration"
/app/db/migration -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"