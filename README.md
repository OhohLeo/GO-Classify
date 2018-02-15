# Classify

Go Rest API and Angular 4 Website to classify all kind of datas

## What is the purpose ?

* Collect all my payslips, store them in a specified directory and print them
* When I scan a document, store the data in the specified document
* Classify all my movies/musics searching through APIs

## How is it organised ?

* Collections: where are stored all validated datas
* Imports: where are from the datas (files, emails, scanner, websites ...)
* Exports: where are written/moved the datas (file, database, printer ...)
* Websites: search for more informations about the data to classify

## Getting Started

To launch the server :

```
go run src/main.go
```

To launch the website :

```
cd web
npm start
```

## Current API

```
GET /references
GET /stream

# Collections

POST   /collections
GET    /collections
GET    /collections/:name
GET    /collections/:name/config
PATCH  /collections/:name/config
PATCH  /collections/:name
DELETE /collections/:name

GET    /collections/:name/buffers
DELETE /collections/:name/buffers
GET    /collections/:name/buffers/:id
POST   /collections/:name/buffers/:id/validate
PATCH  /collections/:name/buffers/:id
DELETE /collections/:name/buffers/:id

GET    /collections/:name/items
DELETE /collections/:name/items
GET    /collections/:name/items/:id
PATCH  /collections/:name/items/:id
DELETE /collections/:name/items/:id

# Imports

POST   /imports
GET    /imports
GET    /imports/config
PUT    /imports/:import/start
PUT    /imports/:import/stop
DELETE /imports/:import

# Exports
```

## Authors

* **Léo Martin** - *Initial work*