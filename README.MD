# Шаблон микросервиса

# Компиляция
Важно добавить тэг go_tarantool_ssl_disable
```
go build -tags go_tarantool_ssl_disable -ldflags="-X 'main/build.Time=$(date)'" cmd\main.go
```

# Добавление logger в другой проект
```
go get -u github.com/san035/basicApiGo@latest
```

или
```
go get -u github.com/san035/basicApiGo@v1.0.6
```

или
```
go mod tidy
```

## git добавление новой версии
git tag v1.0.73
git push --tags