1. Структура проекта
2. Сервер как структура
  * Ресурсы и зависимости
  * Обработчики, middleware
3. Запуск, остановка сервера 
  * обработка сигналов системы
  * работа с контекстом

<h4>Структура проекта</h4>

Общепринятой архитекруры не существует. Архитектура должна подбираться исходя из задач/ситуации

* Плоская архитектура. Подходит для маленьких проектов
* Слоистая/модульная. Малые/средние проекты. [Пример](https://github.com/oralordos/separation/blob/master/main.go)
* Standart project layout. Для средних и крупных проектов - неофицально принятая в сообществе [архитектура](https://github.com/golang-standards/project-layout).

---

<h4>Приложение как структура</h4>

**Экземпляр приложения:** 

Создаём приложение как отдельный тип(структуру) в которую включаем все зависимости. Таким образом мы сможем избежать глобальных переменных, что является плохой практикой.

Если это веб-приложение в зависимости можно включить непосредственно сервер. Это позволит сконфигурировать сервер под наши нужны

```golang
type Application struct {
	server *http.Server
	logger *zap.Logger
	db     *sql.DB
	errors chan error
}
```

**Обработчики, middleware**

В качестве обработчиков лучше использовать не просто функции которые которые реализуют интерфейс `http.Handler`, 
например `func handleAction(w http.ResponseWriter, r *http.Request)`, а функции которые возвращают такой обработчик.

* Обработчики как функции / stateless

```golang
func mainHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get("panic")
		if p != "" {
			panic("panic received")
		}
		w.WriteHeader(200)
		w.Write([]byte("Hello, world"))
	})
}
```

* Структуры как обработчик / stateful

```golang
//handlers.go
type CounterHandler struct {
	counter int
}

func (ct *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ct.counter++
	fmt.Fprint(w, "Counter:", ct.counter)
}
```

```golang
//myapp.go
mux := http.NewServeMux()
counter := &CounterHandler{0}

mux.Handle("/counter", counter)
```

---

**Запуск, остановка сервиса**

* Сервис запускается в отдельной горутине. Коммуникация сервиса с основной горутиной может осуществляется посредством канала с ошибками

```golang
//main.go
app.Run(ctx)
e := <-app.Notify():
fmt.Printf("Server error: %s \n", e.Error())
```

```golang
//app.go
func (app *Application) Run(ctx context.Context) {
	go func() {
		app.errors <- app.server.ListenAndServe()
		close(app.errors)
	}()
}

func (app *Application) Notify() <-chan error {
	return app.errors
}
```

* Обработка системных сигналов

```golang
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
s := <-interrupt
fmt.Printf("Received a signal: %s \n", s.String())
```

* Освобождение ресурсов

Перед завершением основной программы необходимо убедиться что все запущенные горутины завершились. Этого можно добиться посредством передачи `context.Context`.

В случае веб-сервиса в стандартной библиотеке присутствует метод `http.Shutdown(ctx context.Context)`

```golang
func (app *Application) Shutdown() {
	fmt.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.logger.Sync()
	if err != nil {
		log.Fatal(fmt.Errorf("error flushing logger buffer: %s", err.Error()))
	}

	if err = app.server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is down...")
}
```

**Востановление после паники**
Важно! Допускать паники не стоит. Однако если есть места где возможно её возникновение - необходимо правильно восстановить работу сервиса.

В случае веб-сервиса обработку паники можно включить в middleware

```golang
//handlers.go
func mainHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get("panic")
		if p != "" {
			panic("panic received")
		}
		w.WriteHeader(200)
		w.Write([]byte("Hello, world"))
	})
}
```

```golang
//myapp.go
mux := http.NewServeMux()
counter := &CounterHandler{0}

mux.Handle("/", RecoverMiddleware(mainHandler()))
mux.Handle("/counter", RecoverMiddleware(counter))
```

---

Полезные ссылки
- [Пример архитектуры go-сервиса](https://github.com/rtbpanda/go-application-template)
- [Production-ready сервис на Go](https://github.com/PetStores/go-simple/tree/base), [Видео с пояснением](https://youtu.be/yxE5zxTOeUI?t=1822)
- [Как я пишу сервисы спустя 8 лет](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html) - статья с полезными заметками по написанию go-сервисов
- [Golang tutorial](https://tutorialedge.net/golang/) - туториал с большим количеством рецептов и примерами кода 

TODO
* линтеры
* тестирование
* профилирование