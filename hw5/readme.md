# Урок 5. Оптимизация работы с данными. Репликация, секционирование, сегментирование.

## 1. Создайте модель на языке Golang для доступа к таблице активностей с использованием менеджера шардов по аналогии с таблицей профилей пользователей.

При работе с исходным кодом был создан репозиторий с файловой структурой:
```
.
├── activities
│   └── activities.go
├── cmd
│   └── pooler
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── init
│   └── init.sql
├── manager
│   └── manager.go
├── pool
│   └── pool.go
├── readme.md
└── user
    └── user.go
```
В пакете `activities` находится модель для одноименной таблицы в б.д.:
```go
type Activity struct {
	UserId int
	Date   time.Time
	Name   string
}
``` 
Определены методы:
```go
func (a *Activity) connection(m *manager.Manager, p *pool.Pool) (*sql.DB, error)
func (a *Activity) Create(m *manager.Manager, p *pool.Pool) error
func (a *Activity) Read(m *manager.Manager, p *pool.Pool) error
func (a *Activity) Update(m *manager.Manager, p *pool.Pool) error
func (a *Activity) Delete(m *manager.Manager, p *pool.Pool) error
```
Посмотреть в действии можно при работе со следующими заданиями:

## 2. Модифицируйте обе модели таким образом, чтобы лучше использовать реализованный механизм репликации данных: запись изменений надо осуществлять на master, а чтение данных — с хранилища slave.

Для реализации разделения чтения и записи был добавлен параметр `master` в функцию:
```go
func (m *Manager) ShardById(entityId int, master bool) (*Shard, error) {
```
параметр принимает значение `true`, если работа ведется с `master`, в нашем случае во всех операциях, кроме чтения.
Для простейшей реализации мы договариваемся, что номера всех `master` - шардов будут отднозначными, номера всех реплик - на 10 больше.
В зависимости от инфраструктуры эту политику можно сменить. Но для нашей простейшей организации такая реализация будет самой быстрой.


Также для более удобного логгирования было добавлено поле `Role` в структуру `Shards`. 

При запуске `main.go` мы получаем следующие логи при добавлении, чтении и удалении записей таблицы:

```
-----Creating users-----
2022/04/29 11:15:28 operation on shard #1 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
2022/04/29 11:15:28 operation on shard #0 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
-----Creating activities-----
2022/04/29 11:15:28 operation on shard #1 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
2022/04/29 11:15:28 operation on shard #0 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
-----Reading users-----
2022/04/29 11:15:28 operation on shard #12 role: slave
2022/04/29 11:15:28 {11 Jill Biden 69 1}
2022/04/29 11:15:28 operation on shard #11 role: slave
2022/04/29 11:15:28 {1 Joe Biden 78 10}
2022/04/29 11:15:28 operation on shard #10 role: slave
2022/04/29 11:15:28 {15 Donald Trump 74 25}
2022/04/29 11:15:28 operation on shard #12 role: slave
2022/04/29 11:15:28 {26 Melania Trump 52 13}
-----Reading activities-----
2022/04/29 11:15:28 operation on shard #12 role: slave
2022/04/29 11:15:28 {11 2020-11-03 08:00:00 +0000 +0000 wait at home}
2022/04/29 11:15:28 operation on shard #11 role: slave
2022/04/29 11:15:28 {1 2020-11-03 08:00:00 +0000 +0000 election}
2022/04/29 11:15:28 operation on shard #10 role: slave
2022/04/29 11:15:28 {15 2021-01-06 10:00:00 +0000 +0000 2021 United States Capitol attack}
2022/04/29 11:15:28 operation on shard #12 role: slave
2022/04/29 11:15:28 {26 2021-01-06 10:00:00 +0000 +0000 have no idea}
-----Deleting users-----
2022/04/29 11:15:28 operation on shard #1 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
2022/04/29 11:15:28 operation on shard #0 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
-----Deleting activities-----
2022/04/29 11:15:28 operation on shard #1 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
2022/04/29 11:15:28 operation on shard #0 role: master
2022/04/29 11:15:28 operation on shard #2 role: master
```

Видно, что чтение происходит только из `slave` реплик.

## 3. Модифицируйте решение, полученное в задании №2, таким образом, чтобы чтение данных происходило не только из хранилища slave, но и из хранилища master равновероятно.

Для реализации равновероятного чтения, помимо проверки на `master/slave`, был добавлен простейший балансер - целая переменная с поддержкой конкурентной записи.
```go
// Struct to balance between master and slave
type balancer struct {
	sync.Mutex
	balance int64
}
```
Переменная  должна быть всегда равна 0, чтение из мастера увеличивает ее на 1, из слэйва уменьшает.
```go
	if !master {
		m.b.Lock()
		if m.b.balance > 0 {
			m.b.balance--
			n += 10
		} else {
			m.b.balance++
		}
		m.b.Unlock()
	}
}
```
Работа с балансером происходит в функции `ShardById`.

Проверка переменной `master` по сути является проверкой операции чтения.

В результате с той же функцией `main` получаем следующий вывод:

```
-----Creating users-----
2022/04/29 11:42:00 operation on shard #1 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
2022/04/29 11:42:00 operation on shard #0 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
-----Creating activities-----
2022/04/29 11:42:00 operation on shard #1 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
2022/04/29 11:42:00 operation on shard #0 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
-----Reading users-----
2022/04/29 11:42:00 operation on shard #2 role: master
2022/04/29 11:42:00 {11 Jill Biden 69 1}
2022/04/29 11:42:00 operation on shard #11 role: slave
2022/04/29 11:42:00 {1 Joe Biden 78 10}
2022/04/29 11:42:00 operation on shard #0 role: master
2022/04/29 11:42:00 {15 Donald Trump 74 25}
2022/04/29 11:42:00 operation on shard #12 role: slave
2022/04/29 11:42:00 {26 Melania Trump 52 13}
-----Reading activities-----
2022/04/29 11:42:00 operation on shard #2 role: master
2022/04/29 11:42:00 {11 2020-11-03 08:00:00 +0000 +0000 wait at home}
2022/04/29 11:42:00 operation on shard #11 role: slave
2022/04/29 11:42:00 {1 2020-11-03 08:00:00 +0000 +0000 election}
2022/04/29 11:42:00 operation on shard #0 role: master
2022/04/29 11:42:00 {15 2021-01-06 10:00:00 +0000 +0000 2021 United States Capitol attack}
2022/04/29 11:42:00 operation on shard #12 role: slave
2022/04/29 11:42:00 {26 2021-01-06 10:00:00 +0000 +0000 have no idea}
-----Deleting users-----
2022/04/29 11:42:00 operation on shard #1 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
2022/04/29 11:42:00 operation on shard #0 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
-----Deleting activities-----
2022/04/29 11:42:00 operation on shard #1 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
2022/04/29 11:42:00 operation on shard #0 role: master
2022/04/29 11:42:00 operation on shard #2 role: master
```