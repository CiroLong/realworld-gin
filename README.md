# ![RealWorld Example App](logo.png)

> ### [YOUR_FRAMEWORK] codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with **[YOUR_FRAMEWORK]** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **[YOUR_FRAMEWORK]** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

This is a RealWorld backend API implementation using **Go + Gin + GORM + MySQL**. The architecture follows a clean layered design:

- **API Layer**: Handles HTTP requests and responses
- **Service Layer**: Implements core business logic
- **Repository Layer**: Manages data access and persistence
- **Model Layer**: Defines data structures and entities
- **Middleware**: Provides authentication and other cross-cutting concerns
- **Utilities**: Offers common functionalities like JWT, password hashing, etc.

# Getting started

## Prerequisites

- Docker and Docker Compose installed
- Git

## Option 1: Run with Docker (Recommended)

### Step 1: Clone the repository

```bash
git clone https://github.com/CiroLong/realworld-gin.git
cd realworld-gin
```

### Step 2: Start the services

```bash
docker-compose up -d
```

This will:
- Build the Go application
- Start a MySQL database container
- Start the API server container
- Set up network connections between containers

### Step 3: Verify the service is running

```bash
docker-compose ps
```

### Step 4: Test the API

Open your browser or use a tool like Postman to test:

```
http://localhost:8000/api/tags
```

## Option 2: Run locally

### Step 1: Install dependencies

```bash
go mod download
```

### Step 2: Configure the database

Set up a MySQL database and update the `config/config.yaml` file with your database credentials.

### Step 3: Run the application

```bash
go run cmd/server/main.go
```

## Configuration

### Service Port

The API server runs on port **8000** by default. To change this:

1. **In Docker Compose**:
   Edit `docker-compose.yml`:
   ```yaml
   ports:
     - "8080:8000"  # Change 8080 to your desired host port
   ```

2. **In Environment Variables**:
   Edit `docker-compose.yml`:
   ```yaml
   environment:
     - APP_SERVER_ADDR=0.0.0.0:8000  # Change 8000 to your desired container port
   ```

3. **In Configuration File**:
   Edit `config/config.yaml`:
   ```yaml
   server:
     addr: 0.0.0.0:8000  # Change 8000 to your desired port
   ```

### Other Configuration

You can modify other settings through environment variables in `docker-compose.yml`:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| `APP_SERVER_ADDR` | Server address and port | `0.0.0.0:8000` |
| `APP_DATABASE_DSN` | Database connection string | `realworld:realworld@tcp(mysql:3306)/realworld?charset=utf8mb4&parseTime=True&loc=Local` |
| `APP_JWT_SECRET` | JWT signing secret | `your-secret-key-change-in-production` |
| `APP_JWT_EXPIRE_TIME` | JWT token expiration | `24h` |

### Environment Variable Priority

Environment variables take precedence over the configuration file. This allows you to:
- Use `config.yaml` for default settings
- Override specific settings with environment variables
- Keep sensitive information out of the configuration file

## API Documentation

The API follows the RealWorld specification. For detailed API documentation, see:

- [RealWorld API Spec](https://realworld-docs.netlify.app/docs/specs/backend-specs/endpoints)

### Key Endpoints

- **Authentication**: `/api/users`, `/api/users/login`, `/api/user`
- **Profiles**: `/api/profiles/:username`, `/api/profiles/:username/follow`
- **Articles**: `/api/articles`, `/api/articles/feed`, `/api/articles/:slug`
- **Comments**: `/api/articles/:slug/comments`
- **Favorites**: `/api/articles/:slug/favorite`
- **Tags**: `/api/tags`

## Logs and Debugging

To view application logs:

```bash
docker-compose logs app
```

To view database logs:

```bash
docker-compose logs mysql
```

## Stopping the Services

```bash
docker-compose down
```

To remove volumes (for a fresh start):

```bash
docker-compose down -v
```

