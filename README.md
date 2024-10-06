# I. Database Config

1. **Pull Postgres image from Docker Hub:**

   ```sh
   docker pull postgres
   ```

2. **Run Postgres container:**

   ```sh
   docker run --name postgresdb -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
   ```

3. **Connect to Postgres container:**

   ```sh
   docker exec -it postgresdb psql -U postgres
   ```

4. **Create a database:**

   ```sql
   CREATE DATABASE bookdb;
   ```

5. **Navigate to the created database:**

   ```sql
   \c bookdb
   ```

6. **Create book schema:**

   ```sql
   CREATE TABLE books (
       id SERIAL PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       author VARCHAR(255) NOT NULL
   );
   ```

# II. Redis Caching Setup

1. **Pull Redis image from Docker:**

   ```sh
   docker pull redis
   ```

2. **Run Redis container and map port:**

   ```sh
   docker run --name redisdb -p 6379:6379 -d redis
   ```

3. **Execute command in Redis Container:**

   ```sh
   docker exec -it redisdb redis-cli
   ```

4. **Verify Redis is running:**

   ```sh
   ping
   ```

# III. Env Config

Create a [`.env`](command:_github.copilot.openRelativePath?%5B%7B%22scheme%22%3A%22file%22%2C%22authority%22%3A%22%22%2C%22path%22%3A%22%2Fhome%2Fthanhnt%2FProgramming%2Fgo%2Fgin-books-api%2F.env%22%2C%22query%22%3A%22%22%2C%22fragment%22%3A%22%22%7D%2C%22d743820e-d803-4947-8088-bdc1f0215879%22%5D "/home/thanhnt/Programming/go/gin-books-api/.env") file with the following content:

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=bookdb
DB_PORT=5432
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Shanghai
```

# IV. API Request Testing

1. **Create a Book:**

   ```sh
   curl -X POST -H "Content-Type: application/json" -d '{"title":"Golang 101", "author":"John Doe"}' http://localhost:8080/books
   ```

2. **Get All Books:**

   ```sh
   curl -i http://localhost:8080/books
   ```

3. **Get Book by ID:**

   ```sh
   curl -i http://localhost:8080/books/1
   ```

4. **Update Book by ID:**

   ```sh
   curl -X PUT -H "Content-Type: application/json" -d '{"title":"Advanced Golang", "author":"Jane Doe"}' http://localhost:8080/books/1
   ```

5. **Delete Book by ID:**

   ```sh
   curl -X DELETE http://localhost:8080/books/1
   ```