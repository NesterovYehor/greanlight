# Greanlight

Greanlight is a Go-based application designed to manage API interactions, providing a clean, scalable structure for handling web services. It leverages modular design principles, making it suitable for both small and large-scale projects.

## Features

- **API**: Greanlight delivers a structured and easy-to-use API for managing resources.
- **Modular Design**: Organized into separate packages for improved scalability and maintenance.
- **Database Support**: Includes migrations for PostgreSQL to handle schema changes.
- **Configuration**: Offers easy configuration via environment variables for flexible deployment.

## Requirements

- Go 1.x or higher
- PostgreSQL database
- Make (for build automation)

## Installation

1. Clone the repository:
   `git clone https://github.com/NesterovYehor/greanlight.git` and navigate to the project directory.
   
2. Install dependencies using Go modules: `go mod tidy`.

3. Set up your PostgreSQL database and update the environment variables for database connection.

4. Run the migrations to set up the database schema using the command `make migrate-up`.

5. Start the server with `make run`.

## Usage

The API can be accessed via `http://localhost:4000` by default. You can interact with the available endpoints using tools like `curl`, Postman, or through the browser.

### Example Endpoints:

- **POST /users**: Create a new user.
- **GET /users/{id}**: Retrieve user details by ID.
- **POST /login**: Authenticate a user and return a JWT token.

## Directory Structure

- **/cmd/**: Main application entry point.
- **/internal/**: Contains core application logic (services, handlers, etc.).
- **/migrations/**: Database schema migrations.
- **/pkg/**: Shared libraries and utilities.

## Database Migrations

Database migrations are managed using a simple system to ensure the schema is always up to date. Use the following commands:

- `make migrate-up`: Apply all pending migrations.
- `make migrate-down`: Rollback the last migration.

## Configuration

Configuration is managed via environment variables. The following variables need to be set:

- `DB_DSN`: PostgreSQL connection string.
- `PORT`: The port on which the server runs.
- `JWT_SECRET`: Secret key used for JWT token generation.

You can create an `.env` file in the project root for easy management of these variables.

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
