# Домашнее задание 1

## Структура каталога
1.1. Каталог состоит из пользователей и их окружений — проектов, организаций, корпоративных групп и сообществ.

    Реализовано при помощи моделей:
    
```go
type User struct {
	ID     uuid.UUID
	Name   string
	Envs []Env
}

type Env struct {
	ID    uuid.UUID
	Name  string
	Users []User
}

```
1.2. Одно окружение включает в себя много пользователей.
1.3. Один пользователь входит в несколько окружений.

	Эта возможность предусмотрена в простейшем виде в моделях в виде списка, в хранилище должна быть предусмотрена связь много ко многим.

1.4. Система позволяет:

	1.4.1. Добавлять пользователей и окружения.

	реализовано в методах интерфейсов и их реализаций:
```go
type Env interface {
	Create(ctx context.Context, name string) (env models.Env, err error)
	...

type EnvStorage interface {
	Create(ctx context.Context, name string) (env models.Env, err error)
	...

type UserStorage interface {
	Create(ctx context.Context) (user models.User, err error)
	...

type User interface {
	Create(ctx context.Context, name string) (user models.User, err error)
	...

```

	1.4.2. Назначать и убирать пользователей из окружений.

```go
type Connector interface {
	AddToEnv(ctx context.Context, user models.User, env models.Env) error
	GetByEnv(ctx context.Context, env models.Env) ([]models.User, error)
	GetByUser(ctx context.Context, user models.User) ([]models.Env, error)
	DeleteUserFromEnv(ctx context.Context, user models.User, env models.Env) error
}

```
	1.4.3. Искать юзеров по имени или по названию окружения, куда они входят.

	Описано в примере выше. 

	1.4.4. Осуществлять поиск окружений по их названию или по именам входящих в них пользователей.

	Также описано в интерфейсе выше
1.5. Реализовывать полную функциональность не надо. Достаточно показать методы у подходящих структур (без тела).

2. Покажите, что умеете использовать шаблоны `Data Mapper` и `Unit of Work` внутри пакетов.

Data Mapper использован во всех пакетах storage для маппинга б.д. 
Написан только каркас приложения, поэтому нет возможности показать работу с транзакциями 