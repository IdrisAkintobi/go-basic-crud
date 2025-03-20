# go-basic-crud

A basic CRUD application written in Go.

## Overview

This project demonstrates how to build a simple CRUD API for user management and authentication using Go.

## Prerequisites

- [Go](https://golang.org/dl/) 1.16 or later
- [Git](https://git-scm.com/)
- [PostgreSQL](https://www.postgresql.org/) (for database setup)
- [Goose](https://github.com/pressly/goose) (for database migrations)

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/your-username/go-basic-crud.git
   cd go-basic-crud
   ```

2. Create a `.env` using the `.env.example` file:

   ```sh
   cp .env.example .env
   ```

   Update the `.env` file with your PostgreSQL database connection details.

3. Install dependencies:

   ```sh
   go mod tidy
   ```

## Usage

### Migrate the Database

Run the following command to create the database and apply migrations:

```sh
make db.migrate.up
```

This command will create a new database and apply all migrations defined in the `migrations` directory.

Down migrations can be applied using:

```sh
make db.migrate.down
```

### Run Tests

Run tests using the `Makefile`:

```sh
make test
```

This will prepare the test database and run all tests located in the `./tests/` directory.

### Build and Run the Application

You can build and run the application using the following commands:

```sh
go build -o app .
./app
```

After running the application, you can use your favorite HTTP client (e.g., [Postman](https://www.postman.com)) to interact with the API endpoints for Create, Read, Update, and Delete operations.

## API Endpoints

- `POST /` - Create a new user

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
