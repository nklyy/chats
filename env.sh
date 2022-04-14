cat > app.env << EOF
PORT=5000
APP_ENV=development
MONGO_DB_NAME=Example
MONGO_DB_URL=http://127.0.0.1
REDIS_HOST=localhost
REDIS_PORT=6380
SALT=salt
EOF