# Урок 7. Асинхрон, брокеры сообщений. Kafka
1. Запустите большее число копий HTTP-сервера. Эмулируйте нагрузку так, чтобы она поступала от балансировщика на все копии сервиса равномерно. Оцените полученные результаты:

* как меняется rps сервиса;
* как изменяется величина задержки.

api были запущены локально с kafka и redis в контейнерах. К сожалению, после большого количства попыток, сервисы не смогли получить доступ к данным брокера из контейнеров, но получали доступ при запуске локально.

__Данные при одном http-ceрвере.__
```bash
Running 1m test @ http://127.0.0.1:8081
  5 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.01s     9.90ms   1.06s    86.78%
    Req/Sec     0.17      0.37     1.00     83.39%
  295 requests in 1.00m, 21.61KB read
Requests/sec:      4.91
Transfer/sec:     368.17B
```

Для балансировки нагрузки был использован сервер nginx. И простейший балансировщик, с конфигурацией:

```nginx
events { worker_connections 1024; }
http {
  upstream myproject {
    server 127.0.0.1:8081;
    server 127.0.0.1:8082;
    server 127.0.0.1:8083;
  }

  server {
    listen 8080;
    server_name www.domain.com;
    location / {
      proxy_pass http://myproject;
    }
  }
}
```

Так как это первый опыт работы с nginx, то я проверил логи и увидел, что данные запросов действительно проходили через него.
```bash
...
127.0.0.1 - - [21/May/2022:23:55:55 +0300] "POST /rate?rate=10 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:55 +0300] "POST /rate?rate=10 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:56 +0300] "POST /rate?rate=5 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:56 +0300] "POST /rate?rate=5 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:56 +0300] "POST /rate?rate=1 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:56 +0300] "POST /rate?rate=5 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:56 +0300] "POST /rate?rate=5 HTTP/1.1" 200 0 "-" "-"
127.0.0.1 - - [21/May/2022:23:55:57 +0300] "POST /rate?rate=1 HTTP/1.1" 200 0 "-" "-"
...
```

__Данные при трех http-серверах.__
```bash
wrk -c5 -t5 -d1m -s ./wrk.lua 'http://127.0.0.1:8080'
Running 1m test @ http://127.0.0.1:8080
  5 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.01s     4.40ms   1.04s    85.42%
    Req/Sec     0.24      0.43     1.00     76.27%
  295 requests in 1.00m, 37.45KB read
Requests/sec:      4.91
Transfer/sec:     638.20B
```

Можно сделать вывод, что задержка и rps не изменились значительно при увеличении количства http-серверов.

## 2. Самостоятельно замените брокер сообщений, используя образ RabbitMQ или образ NATS. Оцените полученные результаты:

* как меняется rps сервиса;
* как изменяется величина задержки.

Для замены выбрал `NATS`. Так как хотелось поработать с перспективным новым инструментом + очень прост в работе.

```bash
wrk -c5 -t5 -d1m -s ./wrk.lua 'http://127.0.0.1:8081'
Running 1m test @ http://127.0.0.1:8081
  5 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.96ms    2.57ms  44.01ms   89.44%
    Req/Sec   746.65    197.87     1.54k    68.05%
  223150 requests in 1.00m, 15.96MB read
Requests/sec:   3713.58
Transfer/sec:    271.99KB
```

Очень сложно оценить результаты, потому что сервисы были запущены без контейнеров в первом случае, а во втором в контейнерах.
Поэтому разница в разы. При этом значения latency сопоставимы.

## 3. Модифицируйте HTTP-сервис так, чтобы запустить несколько обработчиков process, каждый из которых будет обрабатывать своё сообщение.

Понимая задание таким образом, что должен быть обработчик под каждое сообщение, я придумал только слежующий способ:
- Сервис `process` перестает быть отдельным сервисом и становится пакетом сервиса `api`.
- Мы создаем обработчик каждый раз, когда получаем сообщение, то есть внутри хэндлера пост-запрсоса:
```go
func PostRateHandler(w http.ResponseWriter, r *http.Request) {

	rate := r.FormValue("rate")
	if _, err := strconv.Atoi(rate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go func() {
		atomic.AddInt64(&opened, 1)
		process.New().Proceed(&opened)
	}()
...
```
`opened` - переменная, которая считает количество открытых и закрытых соединений простым изменением счетчика.
Посмотреть его текущее состояние можно используя get-запрос к эндпоинту `/opened`. 

Стандартный тест дает следующий результат:
```bash
wrk -c5 -t5 -d1m -s ./wrk.lua 'http://127.0.0.1:8081'
Running 1m test @ http://127.0.0.1:8081
  5 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    27.56ms   65.78ms 956.21ms   92.21%
    Req/Sec   147.44    168.30     1.40k    88.72%
  28572 requests in 1.00m, 2.04MB read
  Socket errors: connect 0, read 744, write 1280414, timeout 0
Requests/sec:    475.70
Transfer/sec:     34.84KB
```
Мы видим почти десятикратное падение rps и двадцатикратное увеличение средней задержки, что по-видимому, обусловлено временем на создание процесса обработки.

Если мы не хотим чтобы у __КАЖДОГО__ сообщения был свой обработчик, то можно сделать несколько отдельных сервисов `process`и запустить их в отдельных контейнерах.
Однако, в этом случае придется балансировать нагрузку и распределять, какой сервис будет читать в данный момент:

+ простейший способ - случайный топик, в который отправляется сообщение:
```go
rand.Seed(time.Now().Unix())
num := rand.Intn(3) + 1
err = natsConn.Publish(topic+strconv.Itoa(num), []byte(rate))
```
При этом запущено 3 копии сервиса `process`, каждый читает из своего топика
```yaml
process:
    container_name: process
    environment:
    - TOPIC=rates1
    ...
process2:
    container_name: process2
    environment:
    - TOPIC=rates2
    ...
process3:
    container_name: process3
    environment:
    - TOPIC=rates3
    ...
```

Результат теста:
```bash
wrk -c5 -t5 -d1m -s ./wrk.lua 'http://127.0.0.1:8081'
Running 1m test @ http://127.0.0.1:8081
  5 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.04ms    5.91ms 143.36ms   90.38%
    Req/Sec   396.64    153.24     1.14k    69.54%
  118562 requests in 1.00m, 8.48MB read
Requests/sec:   1972.75
Transfer/sec:    144.49KB
```
Трехкратное падение `rps` по сравнению со стандартным вариантом и четырехкратное увеличение `latency`.

Но это гораздо быстрее, чем единый сервис с отдельным обработчиком каждого сообщения.
