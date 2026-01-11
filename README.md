# Middleware Project
## M1 Informatique

## Run

Tidy / download modules :
```
go mod tidy
```
Build & run :
```
cd <program>/cmd
```
then
```
go run main.go
```


## API Key
La clé API qui permet l'envoi de mail ([mail.go](Config/internal/services/mail.go)) n'est pas fourni, il faut donc rajouter votre propre fichier.env, ou modifier le code pour qu'il soit codé en interne.

## Documentation

Documentation is visible in **api** directory:
- ([here for Config](Config/api/swagger.json))
- ([here for TimeTable](TimeTable/api/swagger.json))


## Authors
Jacques KOZIK
Marc MORCOS
