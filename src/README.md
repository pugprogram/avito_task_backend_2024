# Тестовое задание для стажёра Backend

Этот репозиторий содержит реализацию тестового задания на позицию стажёра-бэкендера в Avito. Микросервис написан на Go 1.22.2, СУБД - PostgreSQL.

## Тестовое задание

## Выполнено 
+ Реализована вся основная логика.
+ Реализован расширенный процесс согласования.
+ Реализован просмотр отзывов на прошлые предложения.
+ Реализована возможность оставления отзывов на предложения.
+ Добавлена возможность отката для тендера и предложения.
+ Добавлен файл конфигурации [линтера](https://github.com/golangci/golangci-lint-action). 


## Запуск

Для запуска требуется выполнить следующие шаги:

1. Скачать исходный код и перейти в директорию с проектом:
```text
git clone https://git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725725133-team-77957/zadanie-6105
cd zadanie-6105
```

2. Настроить переменные окружения. В корне проекта есть [.env](.env) файл со значениями переменных окружения по умолчанию.
```text
    SERVER_ADDRESS = <адрес_и_порт_который_будет_слушать_HTTP_сервер>
    SERVER_HOST = <порт_который_будет_слушать_HTTP_сервер>
    POSTGRES_USERNAME = <имя_пользователя_для_подключения_к_PostgreSQL>
    POSTGRES_PASSWORD = <пароль_для_подключения_к_PostgreSQL>
    POSTGRES_HOST = <xост_для_подключения_к_PostgreSQL>
    POSTGRES_PORT = <порт_для_подключения_к_PostgreSQL>
    POSTGRES_DATABASE = <имя_базы_данных_PostgreSQL>
    POSTGRES_CONN = <URL-строка_для_подключения_к_PostgreSQL>
    POSTGRES_JDBC_URL = <JDBC-строка_для_подключения_к_PostgreSQL>

```

3. Для развертывания среды:
    + Запустить контейнеры с базой данных, сервисом и редисом
   ```
   docker-compose -f docker-compose.yml up --build 
   ```

4. Если сервис не запускается, то, пожалуйста, откатитесь до [версии до добавления линтера](https://github.com/pugprogram/avito_task_backend_2024/commit/461c130c4e2c0a390c047a7aca5563481aba2ab2).

# Детали реализации

## Структура проекта
```
zadanie-6105/
|-src/
|  |-cmd/
|  |- service/      точка входа в сервис  
|- interlna/
|  |-api/           API   
|  |-database/      взаимодействие с базой данных
|  |-repository/    слой с бизнес логикой

```

## Кодогенерация

Для автоматической генерации обертки хэндлеров, позволяющей избавиться от ручного парсинга параметров, была использована библиотека [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen). Сгенерированный файл находится в [internal/api/handlers/handlers.go](src/internal/api/handlers/handlers.go).

## Вопросы по ТЗ и их решения

Не все сценарии, описанные в логике задания можно было воспринимать однозначно. Поэтомы были ниже представлены возникшие вопросы и принятые решения.


### Вопрос

Организация состоит только из ответсвенных за организацию или еще из обычных работников?

### Ответ

Было принято решение, что организация состоит из обычных работников и ответсвенных за организацию. Ответственные работники - тоже работники. Была создана вспомогательная таблица в базе данных organization_employee, реализацию которой можно посмотреть в файле migration.sql. 

### Вопрос

Кто может редактировать предложение?

### Ответ

Было принято решение редактировать предложение со всеми статусами может только автор предложения и организация, в которой работает автор.

### Вопрос

Если ли ограничение на username? В api ограничение не было указано, но в базе данных стоял тип Varchar(50).

### Ответ

Поставила ограничение длины 50.

### Вопрос

Какие тендеры могут быть получены при использовании метода Get Tenders, в котором не передается поле username для проверки доступа?

### Ответ

Будем считать, что в ручке Get Tenders, можно вернуть только тендеры со сатусом Published.

### Вопрос

Как относится к переданным пустым значениям в методе Edit Tender и Edit Bid?

### Ответ
Было принято решение при передачи параметров нулевой длины, оставлять параметры старого тендера.

### Вопрос

Кто может взаимодействовать с предложениями при разных статусах?

### Ответ

Было принято решение, что при статусах Created, CANCELED предложение могут просматривать и взаимодейстовать с ним только автор предложения или ответственные за организацию, в которой работает автор.

При статусе Published просматривать и взаимодейстовать с предложением могут автор, ответственные за организацию, а которой работает автор, а также ответственные в организации тендера, на который было добавлено данное предложение.

### Вопрос

Что делать с дендером при достижении кворума подтверждений в предложении?

### Ответ

Было принято решение, что при достижении кворума, тендер закрывается.


## Хотелось, но не успелось

* Пробросить логирование внутрь сервиса
* Написать unit и интеграционные тесты
