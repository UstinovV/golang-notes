<h3>Benchmarks</h3>

Бенчмарки пишутся в `*_test.go` файлах

```golang
func BenchmarkFunctionA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        functionA()
    }
}
```

**Запуск бенчмарков**

`go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out file_test.go`

Где:
* `bench` - список запускаемых бенчмарков
* `benchmem` - отображать потребление памяти и аллокации
* `-cpuprofile=cpu.out -memprofile=mem.out` - сохранение профайлов памяти и ЦП

<h3>Profiling</h3>

Для профилирования используется утилита `go tool pprof`
Чтобы получить исходник для профайлера можно либо использовать бенчмарки, либо напрямую вызывать профайлер в коде
Можно использовать стандартный профайлер `runtime/pprof` или `github.com/pkg/profile`

```golang
func profHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Create("cpuprof.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		w.WriteHeader(200)

		w.Write([]byte(fmt.Sprintf("Hello, %s", r.UserAgent())))
	})
}
```

[Cам pprof](https://github.com/google/pprof)
[Удобный профайлер для Go](https://github.com/pkg/profile)
[Общая статья по диагностике](https://go.dev/doc/diagnostics)
[Perfomance](https://github.com/golang/go/wiki/Performance)
[Статья о pprof в офицальном блоге](https://go.dev/blog/pprof)
[Профилирование и оптимизация веб-приложений на Go](https://habr.com/ru/company/badoo/blog/324682/)
[Специфические вопросы производительности](https://youtu.be/8UESXMJwTpc)
