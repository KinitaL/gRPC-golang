Приложение содержит в себе gRPC сервер и консольный клиент для него. 

Процесс запуска:
1.Настроить файлы mainenv.env и database.env. Изменить файлы docker-compose.yml в соотвествии с используемыми портами, 
указанными в .env файлах.

2.Запустить PostgreSQL:
	docker-compose up -d (в папке postgreSQL)
	Далее, как вариант:
		docker-compose exec database bash
		Внутри контейнера с помощью psql подключиться к БД: psql --host=database --username=youruser --dbname=yourdb
		...Использование интересующих команд...
	p.s: папка database-data содержит тестовую БД, используемую при разработке. Работать можно как с ней, так и с новой БД.
3.Запустить Redis:
	В папке redis в терминале docker-compose up -d
	Также можно подключится к redis через приложение another redis desktop manager для отслеживания записи данных

4.Запустить сервер: 
	В корне проекта go run main.go
5.Запустить клиент:
	В папке client в терминале: go run main.go
	Следовать консольным инструкциям
	


