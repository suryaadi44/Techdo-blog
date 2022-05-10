# TechDo-Blog

### How to run?

1. Install Golang in your machine
2. Install Golang-Migrate for preparing database

    Repository can be found [here](https://github.com/golang-migrate/migrate).
    
3. Run migrate

    Migration can be executed with this command:

    ``` 
    migrate -source file://./db/migrations -database "mysql://root:root@tcp(localhost:3307)/techdo_db" up 
    ```
    Note:
    - `root:root` is your username:password for database.
    - `localhost:3307` is your database address.
    - `techdo_db` is your database name. **Need to be created first!**

3. Change `.env` file value according to your own value.
    Variable description:
    
    | Variable       | Description                                                        |
    | -------------- | ------------------------------------------------------------------ |
    | APP_NAME       | Current app name, not affecting front end, only used for logging   |
    | PORT           | Port which server will listen to, eg: 80                           |
    | DB_ADDRESS     | DBMS address, eg: localhost:3307 / :3307                           |
    | DB_USERNAME    | Database username that have acces to CRUD operation in table       |
    | DB_PASSWORD    | Database password for username                                     |
    | DB_NAME        | Database name where table is stored                                |
    | IMGKIT_PRIVKEY | Private key for IMGKIT account                                     |
    | IMGKIT_PUBKEY  | Public key for IMGKIT account                                      |
    
4. Run app
    App can be started with this command:
    
    ``` 
    go run cmd/main/main.go 
    ```
