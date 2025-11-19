# Go Volume Manager

This project is a Go-based web service that provides an HTTP API for dynamically creating and managing Docker volumes. It listens for requests to create, delete and check the health of the service.

## Project Structure
- `main.go`: The main application entry point, responsible for setting up the web server and routing.
- `handlers/`: Contains the HTTP handlers for the different API endpoints.
- `volume/`: Contains the logic for creating and deleting volumes.

The service listens on port `8203`.

## Endpoints

### Health Check
- **Endpoint:** `/health`
- **Method:** `GET`
- **Description:** Checks the health of the service.
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** `{"status": "ok"}`

### Create Volume
- **Endpoint:** `/create-volume`
- **Method:** `POST`
- **Description:** Creates a new Docker volume.
- **Payload:**
  ```json
  {
    "Name": "my-test-volume",
    "DriverOpts": {
      "size": "5G"
    }
  }
  ```
- **Success Response:**
    - **Code:** 200 OK
    - **Content:** `{"status": "success", "name": "my-test-volume"}`

### Delete Volume
- **Endpoint:** `/delete-volume`
- **Method:** `POST`
- **Description:** Deletes a Docker volume.
- **Payload:**
    ```json
    {
      "Name": "my-test-volume"
    }
    ```
- **Success Response:**
    - **Code:** 200 OK
    - **Content:** `{"status": "success", "name": "my-test-volume"}`

## Building and Running

### Dependencies
- Go 1.17 or later
- Docker
- `fallocate`, `mkfs.ext4`, `mount`, `umount` command-line utilities
- `sudo` access is required for the service to execute system commands.

### Build
To build the application, run the following command:
```bash
go build -o volume-manager .
```

### Run
To run the server, execute the built binary with sudo privileges:
```bash
sudo ./volume-manager
```
The server will start and listen on port `8203`.

## Example Usage

### Create a volume
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "Name": "my-test-volume",
  "DriverOpts": {
    "size": "5G"
  }
}' http://localhost:8203/create-volume
```

### Delete a volume
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "Name": "my-test-volume"
}' http://localhost:8203/delete-volume
```

### Check health
```bash
curl http://localhost:8203/health
```
