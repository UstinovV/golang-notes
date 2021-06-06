Создаём сервер как отдельный тип(структуру) в которую включаем все зависимости

```
type server struct {
  db *Database
  logger *Logger
  router *Router
}
```
