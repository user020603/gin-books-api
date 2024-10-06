# I. Database Config

1. Pull Postgres image from Docker Hub:

   ```sh
   docker pull postgres
   ```

2. Run Postgres container:

   ```sh
   docker run --name postgresdb -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
   ```

3. Connect to Postgres container:

   ```sh
   docker exec -it postgresdb psql -U postgres
   ```

4. Create a database:

   ```sql
   CREATE DATABASE bookdb;
   ```

5. Navigate to the created database:

   ```sql
   \c bookdb
   ```

6. Create book schema:
   ```sql
   CREATE TABLE books (
       id SERIAL PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       author VARCHAR(255) NOT NULL
   );
   ```

# II. Env Config

Create a `.env` file with the following content:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=bookdb
DB_PORT=5432
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Shanghai
```

# III. API Request Testing

1. Create a Book:

   ```sh
   curl -X POST -H "Content-Type: application/json" -d '{"title":"Golang 101", "author":"John Doe"}' http://localhost:8080/books
   ```

2. Get All Books:

   ```sh
   curl http://localhost:8080/books
   ```

3. Get Book by ID:

   ```sh
   curl http://localhost:8080/books/1
   ```

4. Update Book by ID:

   ```sh
   curl -X PUT -H "Content-Type: application/json" -d '{"title":"Advanced Golang", "author":"Jane Doe"}' http://localhost:8080/books/1
   ```

5. Delete Book by ID:
   ```sh
   curl -X DELETE http://localhost:8080/books/1
   ```
