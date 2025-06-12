# News App

A modern news application built with Go, MongoDB, and HTMX.

## Features

- Create, read, update, and delete news posts
- Real-time updates using HTMX
- Responsive design with Tailwind CSS
- Pagination and search functionality
- Modal-based interactions
- Clean architecture with separation of concerns

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/news-app.git
cd news-app
```

2. Start the application:
```bash
docker compose up -d --build
```

3. Open http://localhost:8080 in your browser

That's it! The application is now running with MongoDB database.

## Architecture

The application follows a clean architecture pattern with the following layers:

- **Domain**: Core business logic and entities
- **Repository**: Data access layer (MongoDB implementation)
- **Service**: Business logic layer
- **Handler**: HTTP request handling and response generation
- **Server**: Application configuration and setup

## Prerequisites

- Go 1.21 or higher
- MongoDB 6.0 or higher
- Docker and Docker Compose (for containerized deployment)

## Development

### Project Structure

```
.
├── cmd/
│   └── server/         # Application entry point
├── internal/
│   ├── domain/         # Domain models and interfaces
│   ├── handlers/       # HTTP request handlers
│   ├── repository/     # Data access implementations
│   ├── server/         # Server configuration
│   └── services/       # Business logic
├── pkg/
│   ├── config/         # Configuration management
│   ├── logger/         # Logging setup
│   └── mongo/          # MongoDB client
├── templates/          # HTML templates
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Running the Application

```bash
make run
```

### Docker Commands

```bash
# Build and start containers
make docker-up

# Stop containers
make docker-down

# Rebuild containers
make docker-build
```

## API Endpoints

### Routes

- `GET /`: Main page with posts list
- `GET /posts/new`: Post creation form
- `POST /posts`: Create new post
- `GET /posts/{id}`: View post details
- `GET /posts/{id}/edit`: Edit post form
- `GET /posts/{id}/delete`: Delete post confirmation
- `PUT /posts/{id}`: Update post
- `DELETE /posts/{id}`: Delete post

## HTMX Integration

The application uses HTMX for dynamic content updates without writing JavaScript. Key features:

- Real-time form submissions
- Modal-based interactions
- Dynamic content loading
- Pagination without page reloads
- Search functionality

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 