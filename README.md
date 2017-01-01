# GO-Classify
Classify Collections Tool


GET /references

## Collections

POST   /collections
GET    /collections
GET    /collections/:name
PATCH  /collections/:name
DELETE /collections/:name

## Imports

POST   /imports
GET    /imports
PUT    /imports/:import/start
PUT    /imports/:import/stop
DELETE /imports/:import

### Websites

### Exports


                                    websites
                                       |
	               imports  filter by collection
                          |            |
collections = item [ import data - website data ]
              -> Validate ?
