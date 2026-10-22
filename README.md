# HTTP Server

This is an HTTP server that implements a Twitter clone API, written in Go. It allows users to create accounts, log in, and create and manage short posts called chirps.

## About

This project is part of a backend development course on [Boot.dev](https://www.boot.dev/).

The goal of this project is to learn backend development by building an HTTP server, working with a database, and implementing authentication.

## Features

- User registration and login
- Password hashing
- JWT authentication
- Refresh token management
- Create, retrieve, and delete chirps
- Filter chirps by author and sort them by creation date
- Profanity filtering
- User account updates
- Webhook integration for account upgrades
- API health checks
- File server metrics

## Requirements

- Go
- PostgreSQL
- Git

## Installation

1. Clone the repository.

   ```bash
   git clone https://github.com/faisalalhatimi/http-server.git
   cd http-server
   ```

2. Create a `.env` file in the project root and configure the following environment variables.

   ```env
   DB_URL=your_postgres_connection_string
   PLATFORM=dev
   JWS_SECRET=your_jwt_secret
   POLKA_KEY=your_polka_api_key
   ```

   Replace the example values with your PostgreSQL connection string and your own secrets.

3. Create a PostgreSQL database and apply the SQL schema files in the `sql/schema` directory in numerical order.

4. Install the Go dependencies.

   ```bash
   go mod download
   ```

## Usage

Start the server by running:

```bash
go run .
```

The server listens on port `8080` by default.

### API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/healthz` | Check server health |
| POST | `/api/users` | Create a user |
| PUT | `/api/users` | Update user details |
| POST | `/api/login` | Log in |
| POST | `/api/refresh` | Refresh an access token |
| POST | `/api/revoke` | Revoke a refresh token |
| POST | `/api/chirps` | Create a chirp |
| GET | `/api/chirps` | Retrieve chirps |
| GET | `/api/chirps/{chirpID}` | Retrieve a chirp by ID |
| DELETE | `/api/chirps/{chirpID}` | Delete a chirp |
| POST | `/api/polka/webhooks` | Handle account upgrade webhooks |
| GET | `/admin/metrics` | View file server metrics |
| POST | `/admin/reset` | Reset development data |

### Authentication

Endpoints that require access tokens expect a bearer token in the `Authorization` header.

```http
Authorization: Bearer your_access_token
```

### Query Parameters

The chirps endpoint supports filtering by author and sorting by creation date.

- `author_id`: Filter chirps by author UUID.
- `sort`: Sort chirps using `asc` or `desc`.

For example:

```http
GET /api/chirps?sort=desc
```

### Development

The `/admin/reset` endpoint is restricted to the development environment. Set `PLATFORM=dev` to enable it.
