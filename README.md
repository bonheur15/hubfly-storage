# Go Volume Manager

This project is a Go-based web service that provides an HTTP API for dynamically creating and managing Docker volumes. It listens for requests to create, delete and check the health of the service.

## Project Structure
- `cmd/hubfly-storage/main.go`: The main application entry point, responsible for setting up the web server and routing.
- `handlers/`: Contains the HTTP handlers for the different API endpoints.
- `volume/`: Contains the logic for creating and deleting volumes.

The service listens on port `8203`.

## Features
- **Dynamic Volume Creation**: Create Docker volumes with a specified size.
- **Automatic Cleanup**: Automatically removes the `lost+found` directory upon volume creation.
- **Volume Management**: Delete and list volumes.
- **Volume Statistics**: Get detailed statistics for each volume.
- **Web-based File Access**: Generate temporary URLs for accessing volume data through File Browser.

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
- **Description:** Creates a new Docker volume. The `lost+found` directory is automatically removed upon creation.
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

### Get Volume Stats
- **Endpoint:** `/volume-stats`
- **Method:** `POST`
- **Description:** Gets statistics for a Docker volume.
- **Payload:**
    ```json
    {
      "Name": "my-test-volume"
    }
    ```
- **Success Response:**
    - **Code:** 200 OK
    - **Content:**
      ```json
      {
        "name": "my-test-volume",
        "size": "4.9 GB",
        "used": "8.0 KB",
        "available": "4.7 GB",
        "usage": "1%",
        "mount_path": "/var/lib/docker/volumes/my-test-volume/_data"
      }
      ```

### Get All Volumes
- **Endpoint:** `/dev/volumes`
- **Method:** `GET`
- **Description:** Gets all volumes.
- **Success Response:**
    - **Code:** 200 OK
    - **Content:**
      ```json
      [
        {
          "name": "my-test-volume",
          "size": "4.9 GB",
          "used": "8.0 KB",
          "available": "4.7 GB",
          "usage": "1%",
          "mount_path": "/var/lib/docker/volumes/my-test-volume/_data"
        }
      ]
      ```

### Create URL Volume
- **Endpoint:** `/url-volume/create`
- **Method:** `POST`
- **Description:** Creates a temporary user in File Browser and returns a temporary login URL.
- **Payload:**
  ```json
  {
    "name": "my-test-volume"
  }
  ```
- **Success Response:**
    - **Code:** 200 OK
    - **Content:** `{"url": "http://localhost:8080/login?ott=..."}`


## Building and Running

### Dependencies
- Go 1.17 or later
- Docker
- `fallocate`, `mkfs.ext4`, `mount`, `umount`, `df` command-line utilities
- `sudo` access is required for the service to execute system commands.

### Running the application
To build and run the application, you can use the `run.sh` script. This script will first build the application and then run it with sudo privileges.

```bash
./run.sh
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

### Get volume stats
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "Name": "my-test-volume"
}' http://localhost:8203/volume-stats
```

### Get all volumes
```bash
curl http://localhost:8203/dev/volumes
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
