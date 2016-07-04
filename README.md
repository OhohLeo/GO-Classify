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

POST   /collections/:name/imports
GET    /collections/:name/imports
PUT    /collections/:name/imports/:import/launch
DELETE /collections/:name/imports/:import