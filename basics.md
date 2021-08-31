Links
- https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html

Создаём сервер как отдельный тип(структуру) в которую включаем все зависимости
Таким образом мы сможем избежать глобальных переменных, что является плохой практикой
```
type server struct {
  db *Database
  logger *Logger
  router *Router
}
```

Запускать сервер лучше в отдельной функции, которая возвращает ошибку, и в `main()` просто вызывать эту функцию, так мы будем дерать `main()` чистым

```
func main() {
  //prepare for run
  err := run()
  if err != nil {
    log.Fatal(err)
  }
}

func run() {
  //run server
}

```

Обработчики(handlers) в данном случае являются методами структуры. 
Это позволяет обеспечить им доступ к зависимостям включённым в структуру сервера 
```
func (s *server) handleSomething() http.HandlerFunc { ... }
```

В качестве обработчиков лучше использовать не просто функции которые которые реализуют интерфейс `http.Handler`, 
например `func handleAction(w http.ResponseWriter, r *http.Request)`,
а функции которые возвращают такой обработчик, например:
```
func (s *server) handleSomething(responseFormat string) http.HandlerFunc {
    //do something
    result := doSomething()
    return func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintf(w, format, result)
    }
}
```
В данном примере мы посредством оборачивания функции-обработчика внешней функцией `handleSomething` получаем следующие возможности:
* можем реализовать механизм middlware, т.к. функции могут вкладываться друг в друга, это позволит к примеру проверять авторизацию, включать дополнительную бизнес логику непосредственно перед выполнением функции-обработчика и т.п, результат работы преварительного когда так же можно использовать в функции-обаботчике. В данном примере это `doSomething()`
* если в обработчике требуются дополнительные параметры, которые будут уникальны только для этого обработчика и которые мы не хотим включать в качестве зависимостей в структуру сервера, мы можем передать их в качестве аргумента в функцию обёртки (`responseFormat`)

TODO

* обработка паник
* graceful shutdown, работа с сигналами системы
* линтеры
* тестирование