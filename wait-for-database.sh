# wait for database to start up
while ! pg_isready -h localhost -p 5432 > /dev/null 2> /dev/null; do
    sleep 1
done