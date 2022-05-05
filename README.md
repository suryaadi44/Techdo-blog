# TechDo-Blog

### How to run?

1. Install Golang in your machine
2. Install Golang-Migrate for preparing database

    Repository can be found [here](https://github.com/golang-migrate/migrate).
3. Run migrate

    Migration can be executed with this command:

    ``` migrate -source file://./db/migrations -database "mysql://root:root@tcp(localhost:3307)/techdo_db" up ```

    Note:
    - `root:root` is your username:password for database.
    - `localhost:3307` is your database address.
    - `techdo_db` is your database name.

3. Change `.env` file value according to your own value.

    Variable description:
    
    | Variable    | Description                                                        |
    | ----------- | ------------------------------------------------------------------ |
    | APP_NAME    | Current app name, not affecting front end, only used for logging   |
    | ADDRESS     | Address which server will listen to, eg: localhost:80 / :80        |
    | DB_ADDRESS  | DBMS address, eg: localhost:3307 / :3307                           |
    | DB_USERNAME | Database username that have acces to CRUD operation in table       |
    | DB_PASSWORD | Database password for username                                     |
    | DB_NAME     | Database name where table is stored                                |

4. Run app

    App can be started with this command:
    
    ``` go run cmd/main/main.go ```