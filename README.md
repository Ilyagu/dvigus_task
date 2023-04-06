# dvigus_task

Тестовое задание для Dvigus. Основной алгоритм для рейт лимитов находится в internal->pkg->middlewares->rate_limit.go.
В продакшене я бы хранил статистику в Redis или Tarantool, но решил воспользоваться просто гошной мапой.

## Makefile

### `make run_api`

Run server in console for debug

### `make docker_build`

Stop and delete old container then build new image.

### `make docker_run`

Start new container.

## Environment values

### `COOLDOWN`

Cooldown - int value in seconds. Timeout after error 429.

### `MAX_REQ`

Maximum requests - int value. Maximum requests in set timeout.

### `PERIOD`

Period - int value in seconds. Time limit for maximum requests

### `MASK_SIZE`

Mask size - int value. Can be from 0 to 32.
