# microservices-course-olezhek28
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bahmN/4a62244c443cebb5d8b94cfa746fc7f3/raw/coverage.json)

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0](https://olezhek28.courses/microservices) от [Олега Козырева](http://t.me/olezhek28go).

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
