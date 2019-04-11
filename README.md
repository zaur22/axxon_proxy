# Прокси с историей
Дз для собеседования

Было у меня некоторое подозрение, что нужен был jsonRPC бэкенд. Надеюсь всё таки правильно понял ТЗ.

## Доступные операции
Доступны 3 url для операций
### Операция создания задачи /fetch
json должен иметь поля
{method, path, headers(опц), body(опц))
в ответе придёт
{ID, http-status, headers, content-length}

пример
```
curl --header "Content-Type: application/json" --request POST  --data '{"method":"GET","path":"https://google.com"}'   http://localhost:3000/fetch
```

### Операция получения предыдущих задач /get
json должен иметь поля
{offset(опц), count(опц))
в ответе придёт
[{ID, http-status, headers, content-length}]

пример
```
curl --header "Content-Type: application/json" --request POST  --data '{"count":1}'   http://localhost:3000/get
```

### Операция для удаления задач /delete
json должен иметь поля
{id}
если такого task не существует, придёт bad request
иначе 200 ok

пример
```
curl --header "Content-Type: application/json" --request POST  --data '{"id":"some-id"}'   http://localhost:3000/delete
```