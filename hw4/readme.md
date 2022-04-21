# Создайте набор RED-метрик для методов Exec, Query и QueryRow структуры sql.Db пакета "database/sql", используя библиотеку системы Prometheus.

## Метрики
Набор метрик реализован в пакете [`metrics`](./metrics/metrics.go)

Изначально писалась для каждой функции обертка, которая возвращала функцию с той же сигнатурой:
```go
func (m *Metr) MesurableExec(e func(string, ...interface{}) (sql.Result, error)) func(query string, args ...interface{}) (sql.Result, error) {
	return func(query string, args ...interface{}) (sql.Result, error) {
		t := time.Now()
		if m.on {
			m.Requests.
				WithLabelValues(query, "Exec").
				Inc()
		}
		res, err := e(query, args...)
		if m.on {
			var e string
			if err != nil {
				m.Errors.WithLabelValues(query, "Exec", err.Error()).Inc()
				e = err.Error()
			}
			m.Duration.
				WithLabelValues(query, "Exec", e).
				Observe(time.Since(t).Seconds())
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
```
Но вызов такого варианта выглядит следующим образом:
```go
res, err := st.MesurableExec(st.Exec)("INSERT INTO student(name, lastname, faculty) values (?, ?, ?)",
		name, lastname, faculty)
```
Вызов выполняется неочевидынм образом. 

Поэтому было решено изменить на сигнатуру функции на:
```go
func (m *Metr) MesurableExec(e func(string, ...interface{}) (sql.Result, error), query string, args ...interface{}) (sql.Result, error) {
    ....
```
Функцию можно использовать в качестве замены стандартной `db.Exec`, добавляется один параметр, собственно функция `db.Exec`.

Функция является методом структуры, что повзовляет добавить поле для контроля режима `bypass`, то есть отключения метрик.

Использование метрик осуществляется через embedding cтруктуры с метриками в функции для [работы с базой](./storage/service.go)
```go
type store struct {
	*sql.DB
	*metrics.Metr
}
```
Таким образом можем получить доступ к функциям метрик напрямую.

Пример вызова в коде: 
```go
res, err := st.MesurableExec(st.Exec, "INSERT INTO student(name, lastname, faculty) values (?, ?, ?)",
		name, lastname, faculty)
```

## База данных
Для проверки реализована простая база данных с функциями, имеющими вызов `db.Exec`, `db.Query` и `db.QueryRow`

## Апи. 
Добавлен простенький апи для проверки метрик. 

Примеры запросов:
+ Добавление пользователя
```bash
curl -X POST localhost:8080/add -H "Content-Type: application/json" -d '{"name":"user","lastname":"userlast", "faculty":"tf"}'
```
+ Получения пользователя по фамилии
```bash
curl localhost:8080/get\?lastname\=userlast
```
+ Получение всех пользователей по факультету
```bash
curl localhost:8080/byfaculty\?faculty\=tf
```
__Для получения метрик используется эндпоинт `/metrics`__

## Не удалось

Изначально были идеи сделать одну функцию для обертки над всеми тремя `db.Exec`, `db.Query` и `db.QueryRow`, хотелось сделать эт с учетом того, что две последние функции как минимум по сути - обертки над функцией 
```go
QueryContext(ctx context.Context, query string, args ...any) (*Rows, error)
```
Но это не удалось сделать, потому что разные значения для возврата, и получить `sql.Row` из `sql.Rows` не было возможности, так как поля `sql.Row` не экспортируемы.