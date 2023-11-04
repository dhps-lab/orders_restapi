
Run
```
go run main.go
```

Environments Variables:
```
DB_HOST="localhost"
DB_USER="enerBit"
DB_PASSWORD="admin123"
DB_NAME="enerBit_orders"
DB_PORT="5432"
REDIS_URL="localhost:6379"
REDIS_PASS="enerBit"
REDIS_QUEUE="complete_orders"
```

Set up Databases
```
docker-compose -up -d
```