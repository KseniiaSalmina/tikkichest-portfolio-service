# Portfolio service
Сервис для хранения и управления пользовательскими портфолио с крафтами. Сообщения обо всех изменениях кладёт в кафку (ключ - айди пользователя, он же профайл айди). 

В качестве хранилища используется PostgreSQL. Решение с MongoDB находится в разработке.

## Объекты

Портфолио имеет следующий вид:

    ID          int      `json:"portfolio_id"`
	ProfileID   int      `json:"profile_id"`
	Name        string   `json:"name"`
	Category    Category `json:"category"`
	Description string   `json:"description"`
	Crafts      []Craft  `json:"crafts"`

Категория:

    ID   int    `json:"category_id"`
    Name string `json:"category_name"`

Крафт:

    ID          int       `json:"craft_id"`
    Name        string    `json:"craft_name"`
    Tags        []Tag     `json:"tags"`
    Description string    `json:"craft_description"`
    Contents    []Content `json:"contents"`

Тэг:

    ID   int    `json:"tag_id"`
    Name string `json:"tag_name"`

Контент:

    ID          int    `json:"content_id"`
    Description string `json:"content_description"`
    Data        []byte `json:"data"`

## API
Сервис работает с форматом JSON.

Доступные методы:

    GET /profiles/{profileID}/portfolios - возвращает сокращенные версии (без крафтов) портфолио (на выбор: всех, отобранных по айди профайла или отобранных по айди категории)
	GET /profiles/{profileID}/portfolios/{id} - возвращает сокращённую версию портфолио по его айди
	POST /profiles/{profileID}/portfolios - создаёт новое портфолио 
	PATCH /profiles/{profileID}/portfolios/{id} - редактирует портфолио по его айди
	DELETE /profiles/{profileID}/portfolios/{id} - удаляет портфолио по его айди

	POST /categories - создаёт категорию
	DELETE /categories/{id} - удаляет ктегорию
	GET /categories - выдаёт все категории с их айди

	GET /profiles/{profileID}/portfolios/{id}/crafts - возвращает крафты для выбранного портфолио 
	GET /profiles/{profileID}/portfolios/{id}/crafts/{craftID} - возвращает крафт по его айди
	POST /profiles/{profileID}/portfolios/{id}/crafts - создаёт крафт

	POST /profiles/{profileID}/portfolios/{id}/crafts/{craftID}/tags/{tagID} - добавляет тэг к крафту
	DELETE /profiles/{profileID}/portfolios/{id}/crafts/{craftID}/tags/{tagID} - удаляет тэг крафта

	PATCH /profiles/{profileID}/portfolios/{id}/crafts/{craftID} - редактирует крафт
	DELETE /profiles/{profileID}/portfolios/{id}/crafts/{craftID} - удаляет крафт

	GET /tags/{id}/crafts - возвращает крафты по выбранному тэгу
	GET /tags - возвращает все тэги
	POST /tags - создаёт новый тэг
	DELETE /tags/{id} - удаляет тэг

	POST /profiles/{profileID}/portfolios/{id}/crafts/{craftID}/contents - создаёт контент
	DELETE /profiles/{profileID}/portfolios/{id}/crafts/{craftID}/contents/{contentID} - удаляет контент
	PATCH /profiles/{profileID}/portfolios/{id}/crafts/{craftID}/contents/{contentID} - редактирует контент

Методы, в теории возвращающие больше одного объекта, на самом деле вернут объект следующего вида:

	{Object}    []{Object}  `json:"{objects}"` // objects - это portfolios, categories, crafts или tags
	PageNo      int         `json:"page_number"`
	Limit       int         `json:"limit"`
	PagesAmount int         `json:"pages_amount"`

## Kafka
Сервис после каждого обновления отправляет в кафку сообщение с айди пользователя в качестве ключа и объектом JSON в качестве значения:

    Object   Object `json:"object"` 
    ObjectID int    `json:"object_id"`
    Change   Change `json:"change"` 

Список доступных значений поля Object:

    Portfolio Object = "portfolio"
	Craft     Object = "craft"
	Content   Object = "content"

Список доступных значений поля Change:

    CreateObj Change = "created"
	UpdateObj Change = "changed"
	DeleteObj Change = "deleted"

## Переменные окружения

Сервис умеет считывать переменные из файла .env в директории исполняемого файла (в корне проекта).

В примерах указаны дефолтные значения. Если программа не сможет считать пользовательские env, то возьмет их.

Переменные сервера:

    SERVER_LISTEN=:8088
    SERVER_READ_TIMEOUT=5s
    SERVER_WRITE_TIMEOUT=5s
    SERVER_IDLE_TIMEOUT=30s

Переменные Postgres:

    PG_USER=
	PG_PASSWORD=
	PG_HOST=localhost
	PG_PORT=5432
	PG_DATABASE=

Переменные Kafka:

    KAFKA_HOST=localhost
	KAFKA_PORT=9092
	KAFKA_TOPIC=