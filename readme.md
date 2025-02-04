# Lead Stream Service

Lead Stream Service is a web service for managing schemas and uploading files associated with those schemas. It uses MongoDB for data storage and Huma for API documentation and routing.

## Project Structure

```
.github/
└── workflows/
    └── go.yml
cmd/
├── main.go
└── run_tests.sh
internal/
├── api/
│   ├── handlers/
│   │   ├── error_handler.go
│   │   ├── file_handler.go
│   │   └── schema_handler.go
│   └── router.go
├── configuration/
│   └── config.go
├── domain/
│   ├── errors.go
│   ├── file.go
│   └── schema.go
├── infrastructure/
│   └── mongo_connection.go
├── integration/
│   ├── resources/
│   │   └── file/
│   │       ├── test_file_handler_fail.csv
│   │       ├── test_file_handler_fail_2.csv
│   │       ├── test_file_handler_fail_3.csv
│   │       ├── test_file_handler_fail_4.csv
│   │       └── test_file_handler_success.csv
│   ├── file_integration_test.go
│   ├── schema_integration_test.go
│   └── server_test.go
├── repositories/
│   ├── lead_repository.go
│   └── schema_repository.go
├── services/
│   ├── file_service.go
│   ├── mocks_service_test.go
│   ├── schema_service.go
│   └── schema_service_test.go
└── tools/
    └── directory.go
config.yaml
go.mod
go.sum
```

## Technologies Used

- **Go**: The primary programming language used for the service.
- **MongoDB**: The database used for storing schemas and leads.
- **Huma**: A framework for building and documenting APIs.
- **Echo**: A high-performance, extensible, minimalist web framework for Go.
- **Testify**: A toolkit with common assertions and mocks that plays nicely with the standard library.
- **YAML**: Used for configuration files.

## Getting Started

### Prerequisites

- Go 1.23
- MongoDB

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/vitortenor/lead-stream-service.git
    cd lead-stream-service
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Set up your MongoDB instance and update the `config.yaml` file with your MongoDB URI and database details.

### Configuration

The configuration file `config.yaml` should look like this:

```yaml
server:
  api:
    name: "Lead Stream Service"
    version: "1.0.0"
  host: "localhost"
  port: 8080
database:
  uri: "mongodb://localhost:27017"
  name: "lead_stream_db"
  collection:
    schemas: "schemas"
    leads: "leads"
```

### Running the Service

To start the service, run:

```sh
go run cmd/main.go
```

The service will be available at `http://localhost:8080`.

### Running Tests

To run the tests, use the following command:

```sh
./cmd/run_tests.sh
```

## API Endpoints

### Schemas

- **Create Schema**
  - **URL:** `/schema`
  - **Method:** `POST`
  - **Description:** Create a new schema with the given fields.

### Files

- **Upload File**
  - **URL:** `/schema/{schemaId}/file`
  - **Method:** `POST`
  - **Description:** Upload a file to the given schema.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the MIT License.
