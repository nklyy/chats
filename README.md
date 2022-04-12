# Get started

### 1. Setup envs in .env file
```
APP_PORT=
APP_ENV=

MONGO_PORT=
MONGO_DB_NAME=
MONGO_DB_URL=(if you will app via docker-compose, you don't need to setup this env)

REDIS_HOST_AUTH=
REDIS_PORT_AUTH=

REDIS_HOST_CHAT=
REDIS_PORT_CHAT=

SALT=

JWT_SECRET_ACCESS=
JWT_EXPIRY_ACCESS=
JWT_SECRET_REFRESH=
JWT_EXPIRY_REFRESH=
AUTO_LOGOUT=
```

### 2. Start tests
``` makefile
make test
```

### 3. Start app manually
``` makefile
make run
```
>:warning: **When you start app not via docker, please set redis and mongo envs from your instances or docker containers**

### 4. Start app via docker
``` makefile
run-docker
```

***If you want to see front-end part, visit [repository](https://github.com/nn-labs/noname-realtime-support-chat)***

<img src="https://github.com/nn-labs/noname-realtime-support-chat/workflows/development/badge.svg?branch=dev"><br>