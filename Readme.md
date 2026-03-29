
# Errandify - Backend API

Errandify is a modern task management system designed to streamline the workflow between Admins (task assigners) and Employees (task executors). This repository contains the backend API built with Golang.

## Tech Stack
* **Language:** Go (Golang)
* **Framework:** Gin Web Framework
* **ORM:** GORM
* **Database:** PostgreSQL
* **Deployment:** Docker & Docker Compose

---

## Getting Started

Follow these instructions to set up and run the project on your local machine.

### Prerequisites
* [Go](https://golang.org/doc/install) (v1.18 or newer)
* [Docker Desktop](https://www.docker.com/products/docker-desktop/) or Docker Engine
* [Postman](https://www.postman.com/) (Optional, for API testing)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/hanifahmadzakir/Errandifyv1.git](https://github.com/hanifahmadzakir/Errandifyv1.git)
    cd Errandifyv1
    ```

2.  **Environment Variables setup:**
    Create a `.env` file in the root directory and add your configuration:
    ```env
    PORT=8080
    GIN_MODE=debug

    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=admin
    DB_PASSWORD=admin
    DB_NAME=mydb

    JWT_SECRET=your_super_secret_key
    ```

3.  **Start the Database via Docker:**
    ```bash
    docker-compose up -d
    ```

4.  **Install dependencies and run the server:**
    ```bash
    go mod tidy
    go run main.go
    ```
    The server should now be running on `http://localhost:8080`.

---

## API Documentation

Below is the list of available endpoints. 

### 1. Users

| Method | Endpoint | Description | Body Type |
| :--- | :--- | :--- | :--- |
| `POST` | `/users` | Add a new user | JSON |
| `DELETE` | `/users/:id` | Delete user by ID | None |
| `GET` | `/users/Employee` | Get list of all employees | None |
| `POST` | `/users/login` | Login for authentication | JSON |

**Login Payload Example:**
```json
{    
    "email": "owner@go.id",
    "password": "123456"
}