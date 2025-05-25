# Image Previewer Microservice

![Microservice](https://img.shields.io/badge/Type-Microservice-blue) 
![Go](https://img.shields.io/badge/Go-1.20-brightgreen) 
![Libvips](https://img.shields.io/badge/Libvips-required-orange)

Микросервис для динамической обработки изображений с поддержкой кэширования и различных стратегий ресайза.

## ✨ Основные возможности

- **Ресайз изображений** с поддержкой различных стратегий:
  - `fill` - заполнение области с обрезкой
  - (TODO) `fit` - вписание в область без обрезки
- **Поддержка форматов**: JPEG, PNG, WebP
- **Кэширование результатов** обработки (LRU-кэш)
- **Работа с удаленными источниками** изображений
- (TODO) **Гибкая конфигурация** через JSON-файл
- **Логирование** операций с настраиваемым уровнем детализации

## 🛠 Установка и настройка

### Предварительные требования

```bash
sudo apt-get install -y libvips-dev
```

### Конфигурация

Файл `configs/config.json`:

```json
{
  "host": ":8080",
  "timeout": 30,
  "cache_capacity": 1000,
  "max_body_size": 10485760,
  "logger": {
    "level": "debug",
    "output": "./logs/previewer.log"
  }
}
```

| Параметр         | Описание                          | По умолчанию |
|------------------|-----------------------------------|--------------|
| host             | Адрес и порт для запуска сервера  | :8080        |
| timeout          | Таймаут операций (сек)            | 30           |
| cacheСapacity   | Размер кэша (кол-во элементов)    | 1000         |
| maxBodySize    | Макс. размер обрабатываемого изображения (байт) | 10MB     |
| logger.level     | Уровень логирования (debug, info, warn, error) | debug |
| logger.output    | Файл для записи логов             | ./logs/previewer.log |

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
make unit-test                                          # запуск unit-тестов
make integration-test SRC_HOST=source-host.local:8081   # запуск интеграционных тестов
make lint                                               # проверка кодстайла
```

## 📌 Примеры использования

### Базовые запросы

1. **Заполнение области с обрезкой**:
   ```
   http://my-resizer.local/fill/600/600/source.site/image.jpg
   ```

2. (TODO) **Вписание в область без обрезки**:
   ```
   http://my-resizer.local/fit/300/300/source.site/image.png
   ```

### Параметры URL:

- Первый сегмент: стратегия ресайза (`fill`, `fit (TODO)`)
- Далее: ширина и высота результата
- Последний сегмент: URL исходного изображения

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
