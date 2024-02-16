# example-api

## Инструкция по запуску:
В файлах конфигураций config.yml и .env Подставить необходимые значения для запуска.

В корневой папке проекта прописать команду:

`$ make up`

Для тестирования некоторых компонентов проекта:

`$ make test`

### Основные End-Point:
### `/ReserveProduct` - Отвечает за создание резервирования какого-либо товара(ов), если товаров не хватает, или его нет, то выдаст ошибку.

Есть три обязательных поля: номер склада, массив товаров для резервирования, количество товара, который мы резервируем. Примеры запросов (cURL):

curl --location 'http://host/ReserveProduct' \
--header 'Content-Type: application/json' \
--data '{
"warehouse_id": 2,
"unique_codes": [
"olkiuj",
"tghyuj"
],
"counts": [
1000,
1200
]
}'

curl --location 'http://host/ReserveProduct' \
--header 'Content-Type: application/json' \
--data '{
"warehouse_id": 1,
"unique_codes": [
"olkiuj"
],
"counts": [
5
]
}'

Примеры ответов:
Статус код 200 и тело, если все продукты успешно зарезервированы:

`{"successful": [{"id": 9,"unique_codes": "olkiuj"}]}`

Статус код 207 и тело, если часть продуктов зарезервировано. i-я ошибка соответствует i-му товару.  
Тело:   
`{"successful": [{"id": 13,"unique_codes": "tghyuj"}],"unsuccessful": ["olkiuj"],"errors": ["not enough product"]}`

Если все товары не были зарезервированы по тем или иным причинам придет ответ с статус кодом 400. Тело ответа может быть разным. Примеры:

`{"successful":[{"id":13,"unique_codes":"tghyuj"}],"unsuccessful":["olkiuj"],"errors":["not enough product"]}` - не хватило продукта на складе, чтобы зарезервировать

или

`{"error": "invalid warehouse id"}`  - если недоступен склад

### `/FreeReservation`  - освобождение склада от резервирования.

На вход приходит id резервирования. Т.к. один и тот же товар может быть зарезервирован на складе много раз, то конкретное резервирование можно найти по id.

curl --location 'http://host/GetAllProducts' \
--header 'Content-Type: application/json' \
--data '{
"warehouse_id": 1
}'

curl --location 'http://host/FreeReservation' \
--header 'Content-Type: application/json' \
--data '{
"id": [1, 3, 81]
}'

Примеры ответа:
200 OK
`{"products_codes":["olkiuj"]}`

207 Multi - Status

`{"successful":[1,3],"unsuccessful":[81],"errors":["non-existent reservation id"]}`

400 Bad request
`{"errors":["non-existent reservation id"],"unsuccessful":[81]}`
or
`{"error":"invalid body"}`

### `GetAllProducts`  - получение всех продуктов на складе

На вход приходит идентификатор склада.

curl --location 'http://host/GetAllProducts' \
--header 'Content-Type: application/json' \
--data '{
"warehouse_id": 1
}'

Примеры ответа:

Если товаров на складе под этим идентификатором не найдено, но запрос выполнен успешно, то вернется ответ: 204 No content (без тела ответа)

Если ошибка в теле запроса

400 Bad request

`{"error":"invalid request body"}`
{"error":"invalid request body"}    


### **Обязательные требования**

· Использование go fmt и goimports

· Следование Effective Go

· Go актуальной версии

· Использование JSON-API. Каждая операция должна быть RPC-like, то есть выполнять определенное законченное действие.   
  
· PostgreSQL или MySQL в качестве хранилища данных

· Наличие команды make up в Makefile, которая: поднимает без ошибок приложение с помощью Docker контейнеров, и готовую инфраструктуру для работы приложения (база данных, миграции, данные для тестирования работы приложения)

· Описание API методов с работающим запросом и ответом в одном из следующих форматов: .http файлы (IDEA) с ответами, curl команды с ответами в README.md, коллекция Postman, построенная на основе swagger / openapi коллекции.


### **Критерии оценки**

    · Работоспособность API        · API выполняет заявленные функции  
        · API предусматривает граничные кейсы  
        · Нахождение и решение потенциальных проблем  
        · Организация и читаемость кода  
        · Обработка ошибок  

### Будет плюсом

    · Покрытие кода unit или функциональными тестами        · Аргументация выбора пакетов в go.mod, приложить отдельным файлом packages.md  

### **Результат**

    · Проект должен быть выложен в публичный репозиторий Github/Gitlab        · В проекте должен присутствовать README и содержать в себе:  
        · Инструкцию по запуску сервиса  
        · Инструкцию по запуску тестов при их наличии  

## Задание

### **#1. API для работы с товарами на складе**

Необходимо спроектировать и реализовать API методы для работы с товарами на одном складе.     
Учесть, что вызов API может быть одновременно из разных систем и они могут работать с одинаковыми товарами.    
Методы API можно расширять доп. параметрами на своё усмотрение

· Спроектировать и реализовать БД для хранения следующих сущностей

    · Склад            
            ·название  
            ·признак доступности  
    · Товар  
            ·название  
            ·размер  
            ·уникальный код  
            ·количество  

**Реализовать методы API:**

· резервирование товара на складе для доставки

· на вход принимает:

· массив уникальных кодов товара

· освобождение резерва товаров

· на вход принимает

· массив уникальных кодов товара

· получение кол-ва оставшихся товаров на складе

· на вход принимает:

· идентификатор склада

**Будет плюсом**

· Реализация логики работы с товарами, которые одновременно могут находиться на нескольких складах