# Backend - Little Great Programmers

## Build
1. Agregar un archivo `.env` simil a `.env_example`.

2. Levantar el backend con `docker compose`

```bash
docker compose up -d --build
```
3. Ver logs con

```bash
docker compose logs -f
```
4. Frenar la aplicación
```bash
docker compose down
```

## Routes

Para ver la documentación de la API, acceder a `localhost:8001/swagger/index.html`.

## Contribuir

Una vez que se actualiza la documentación en un controller, se debe ejecutar 

```bash
swag init
```