# Image Previewer Microservice

[![Go Report Card](https://goreportcard.com/badge/github.com/IKolyas/thumbnailer)](https://goreportcard.com/report/github.com/IKolyas/thumbnailer)

Микросервис для динамической обработки изображений с поддержкой кэширования и различных стратегий ресайза.

## ✨ Основные возможности

- **Ресайз изображений** с поддержкой различных стратегий:
  - `fill` - заполнение области с обрезкой
  - (TODO) `fit` - вписание в область без обрезки
- **Поддержка форматов**: JPEG, PNG, WebP
- **Кэширование результатов** обработки (LRU-кэш)
- **Работа с удаленными источниками** изображений
- **Гибкая конфигурация** через JSON-файл
- **Логирование** операций с настраиваемым уровнем детализации

## 🛠 Установка и настройка

### Предварительные требования

> при локальном запуске на системе Linux потребуется установка зависимости

```bash
sudo apt-get install -y libvips-dev
```

### Конфигурация

Файл `configs/config.json`:

```json
{
  "host": ":8080",
  "timeout": 30,
  "cacheCapacity": 1000,
  "maxBodySize": 10485760,
  "storageDir": "./storage/",
  "logger": {
    "level": "debug",
    "output": "./logs/previewer.log"
  }
}
```

| Параметр         | Описание                                       | По умолчанию         |
|------------------|------------------------------------------------|----------------------|
| host             | Адрес и порт для запуска сервера               | :8080                |
| timeout          | Таймаут операций (сек)                         | 30                   |
| cacheСapacity    | Размер кэша (кол-во элементов)                 | 1000                 |
| maxBodySize      | Макс. размер обрабатываемого изображения (байт)| 10MB                 |
| storageDir       | Директория хранения файлов кеша                | ./storage/           |
| logger.level     | Уровень логирования (debug, info, warn, error) | debug                |
| logger.output    | Файл для записи логов                          | ./logs/previewer.log |

## 🚀 Запуск сервиса

### Сборка и запуск

```bash
make build    # только сборка
make run      # сборка и запуск
```

### Docker

```bash
make docker-build    # сборка Docker-образа
make docker-run      # запуск контейнера
make docker-stop     # остановка контейнера
```

## 🧪 Тестирование

```bash
make unit-test          # запуск unit-тестов
make integration-test   # запуск интеграционных тестов
make lint               # проверка кодстайла
```

## 📌 Примеры использования

### Базовые запросы

1. **Заполнение области с обрезкой**:
   ```
   http://my-resizer.local/fill/600/600/https://source.site/image.jpg
   ```

2. (TODO) **Вписание в область без обрезки**:
   ```
   http://my-resizer.local/fit/300/300/https://source.site/image.png
   ```

### Параметры URL:

- Первый сегмент: стратегия ресайза (`fill`, `fit (TODO)`)
- Далее: ширина и высота результата
- Последний сегмент: URL исходного изображения (source.site/image.png | http://source.site/image.png | https://source.site/image.png)

## 📊 Логирование

Логи сохраняются в файл `./logs/previewer.log` с указанным уровнем детализации.

Уровни логирования:
- `debug` - максимальная детализация
- `info` - основная информация о работе
- `warn` - только предупреждения и ошибки
- `error` - только критические ошибки

## � Очистка

```bash
make clean  # удаление артефактов сборки
```
