# Создайте сервис, который по REST API и GRPC получает запросы на добавление, изменение, удаление и чтение прайс-листов по их уникальному идентификатору. 

* __Каждый прайс-лист состоит из набора товаров и цен к ним.__
Для хранения используется файловая система. Реализован простой CRUD интерфейс хранилища.
```go
type Storage interface {
	Create(ctx context.Context, list models.List) error
	Read(ctx context.Context, id uuid.UUID) (list *models.List, err error)
	Update(ctx context.Context, id uuid.UUID, items []*models.Item) error
	Delete(ctx context.Context, id uuid.UUID) error
}
```
Хранилище используется обоими сервисами (REST и GRPC).

* Для выполнения задания используйте кодогенераторы имплементации API из этого урока, а также напишите сами спецификации API (OpenAPI, GRPC).

1) Для реализации REST-интерфейса была использована библиотека `oapi-codegen` и вариант реализации сервера с помощью `echo`. Генерировалось на основе [openapi спецификации](./api/api.yml)

После генерации файлов была изменена структура моделей и вместо сгенерированных, вставлены модели, описанные в пакете [`models`](./pkg/models/models.go).

```go
// CreateListParams defines parameters for CreateList.
type CreateListParams struct {
	List *models.List `form:"list,omitempty" json:"list,omitempty"`
}

// UpdateListObjectJSONBody defines parameters for UpdateListObject.
type UpdateListObjectJSONBody = models.List

// UpdateListObjectJSONRequestBody defines body for UpdateListObject for application/json ContentType.
type UpdateListObjectJSONRequestBody = UpdateListObjectJSONBody

```

2) GRPC сервер сгенеирован при помощи библиотеки [protoc-gen-go](google.golang.org/protobuf/cmd/protoc-gen-go@v1.26) на основе [proto-файла](./service.proto)

## Запуск и тестирование
Запуск можно осуществлять при помощи `Makefile`:

* Для запуска rest-сервера
```bash
make run-rest
```
* Для запуска grpc
```bash
make run-rpc
```
* Для запуска обоих серверов
```bash
make run-both
```
Для тестирования grpc можно использовать следуюзие команды:
```bash
grpcurl -plaintext -d '{"id":"","items":[{"name":"test1","price":10},{"name":"test2","price":20}]}' localhost:9000 ListService/Create

grpcurl -plaintext -d '{"id":"fcab6956-5f4d-4009-942c-7122a7976316"}' localhost:9000 ListService/Read

grpcurl -plaintext -d '{"id":"fcab6956-5f4d-4009-942c-7122a7976316","items":[{"name":"test2","price":3000}]}' localhost:9000 ListService/Update

grpcurl -plaintext -d '{"id":"a00a7e9c-29e9-4e21-a28c-333319a2b02e"}' localhost:8000 ListService/Delete

```

Для тестирования rest можно использовать следующие команды:
```bash
curl -H "Content_Type: application/json" -X POST localhost:8080/list/create \
-d '{"id":"311e0467-a754-4286-bf89-47c3a83eeb68","items":[{"name":"test1","price":10},{"name":"test2","price":20}]}'

curl localhost:8080/list/311e0467-a754-4286-bf89-47c3a83eeb68

curl -X DELETE localhost:8080/list/delete/311e0467-a754-4286-bf89-47c3a83eeb68

curl -H "Content_Type: application/json" -X PATCH localhost:8080/list/update \
-d '{"id":"311e0467-a754-4286-bf89-47c3a83eeb68","items":[{"name":"test1","price":100},{"name":"test2","price":20}]}'
```
Для запуска тестов можно использовать:
```bash
make test-all
```