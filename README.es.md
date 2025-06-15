# Chirpy


## Esp

### Descripción

Chirpy es un programa servidor de backend en Go que simula aspectos de twitter.
Este proyecto fue diseñado para aprender y practicar conceptos y nociones de
servidores http.
El programa maneja pedidos y respuestas http en formato de JSON.
No hay frontend, el programa solo maneja la lógica del backend.
Los endpoints http y funcionalidad son explicados en mas detalle debajo.
El programa usa PostgreSQL para la administración de la base de datos, sqlc
para generar código Go a partir de SQL y pressly/Goose para realizar las
migraciones de la base de datos.

### Requerimientos 

Para correr este programa debes tener instaladas las últimas versiones
de Go y PostgreSQL en tu sistema.

### Preparación

- Clona el repositorio e instala las dependencias:

```bash
git clone https://github.com/iegpeppino/chirpy
cd path/to/chirpy
go install ...
```

- Crea una base de datos PostgreSQL llamada chirpy 

- Almacena el URL de la base de datos, una secret key, y la Polka ApiKey (es ficticia) en variables del sistema.

- Realiza las migraciones de db:

```bash
cd sql/schema
goose postgres <connection_url> up
```

- Regresa al directorio base, compila y corre el programa:

```bash
cd ../..
go build -o out && ./out
```

### Funciones

Puedes crear, ingresar y actualizar usuarios usando un email y contraseña.
Las contraseñas pasan por el proceso de hash y el usuario recibe un token para
acceder a los demás endpoints. El token se refresca, y también puede ser revocado.
Hay un webhook para verificar la cuenta como 'chirpy_red', a través de una Api 
ficticia llamada Polka.
Los usuarios logeados pueden crear chirps (twitts) y eliminarlos (si son los dueños del mismo).
Se pueden enumerar todos los chirps de la base de datos o de un usuario en particular.
También se puede buscar un chirp específico por ID.
La mayoría de los pedidos requieren una validación del token JWT que se realiza antes de
procesar la respuesta.

### Uso

#### Endpoints Disponibles

* Crear Usuario
    - Método: "POST"
    - Endpoint: /api/users
    - Params: 
        ```json
        {
            email: "user@email.com"
            password: "mypasswordsnot123"
        }
        ```

* Logear Usuario
    - Método: "POST"
    - Endpoint: /api/login
    - Params: 
        ```json
        {
            email: "user@email.com"
            password: "mypasswordsnot123"
        }
        ```

* Actualizar email y contraseña de Usuario
    - Method: "PUT"
    - Endpoint: /api/users
    - Params:
        ```json
        {
            email: "user@email.com"
            password: "mypasswordsnot123"
        }
        ```

* Borrar todos los usuarios de la db
    * Advertencia: La plataforma del cliente debe ser "dev"
    - Método: "POST"
    - Endpoint: /admin/reset

* Crear chirp (twitt)
    - Método: "POST"
    - Endpoint: /api/chirps
    - Params:
        ```json
        {
            body: "hello there!"
        }
        ```

* Enumerar todos los chirps
    Nota: Por defecto, enumera todos los chirps de la db, pero se puede
    pasar un "user_id" y "order" (asc o desc) como parámetros de query url para 
    enumerar todos los chirps de un usuario selecto y ordenarlos por fecha.
    - Método: "GET"
    - Endpoint: /api/chirps

* Buscar chirp por ID
    - Method: "GET"
    - Endpoint /api/chiprs/{chirpID}

* Borrar chirp por ID
    - METHOD: "DELETE"
    - Endpoint: /api/chirps/{chirpID}

* Verificar usuario a "chirpy_red_status"
    * Advertencia: el parametro de evento debe ser "user.upgraded" para que funcione
    - Método: "POST"
    - Endpoint: /api/polka/webhooks
    - Params:
        ```json
        {
            event: "some.event"
            data: {
                user_id: "asfj0123jf1093jf"
            }
        }
        ```

* Refrescar Token JWT de Acceso
    - Método: "POST"
    - Endpoint: /api/refresh
    - Params:
        ```json
        {
            token: "userJWTtokenstring"
        }
        ```

* Revocar Token de Acceso
    - Método: "POST"
    - Endpoint: /api/revoke

* Otros endpoints menos críticos:
    - "GET /api/healthz" : Responde con un código 200 para indicar que el server funciona
    - "GET /admin/metrics" : Regresa el número de visitas al server desde que inició