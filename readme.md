# Toolkit 🧰

A Go module with commonly used tools for various operations.

## Installation 💻

To install the toolkit, use the following command:

```bash
go get -u github.com/javikuma/toolkit
```

## Usage 🎲

Import the package:
```go
import "github.com/javikuma/toolkit"
```

Create a new instance of the Tools struct:
```go
t := &toolkit.Tools{}
```

## Available Tools 🎯

### 1. Random String Generation

Generates a random string of specified length.

```go
randomString := t.RandomString(length int) string
```

### 2. File Upload

Uploads a single file.

```go
uploadedFile, err := t.UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error)
```

Uploads multiple files.

```go
uploadedFiles, err := t.UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error)
```

### 3. Directory Creation

Creates a directory and all necessary parent directories if they don't exist.

```go
err := t.CreateDirIfNotExist(path string) error
```

### 4. Slug Generation

Creates a URL-safe slug from a string.

```go
slug, err := t.Slugify(s string) (string, error)
```

### 5. Static File Download

Downloads a static file and sets appropriate headers for browser handling.

```go
t.DownloadStaticFiles(w http.ResponseWriter, r *http.Request, pathname, fileName, displayName string)
```

### 6. JSON Operations

Reads JSON from request body into a Go data structure.
```go
err := t.ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error
```

Writes JSON response with specified status code and optional headers.
```go
err := t.WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error
```

Sends a JSON error response.
```go
err := t.ErrorJSON(w http.ResponseWriter, err error, status ...int) error
```

### 7. Remote JSON Push

Posts JSON data to a remote URL and returns the response.
```go
response, statusCode, err := t.PushJSONToRemote(uri string, data interface{}, client ...*http.Client) (*http.Response, int, error)
```

## Customization 🧮

The `Tools` struct allows customization of certain parameters:

* `MaxFileSize`: Maximum allowed file size for uploads. **_Type:_** int | **_e.g._** 1024 * 1024 | **_default:_** 1GB
* `AllowedFileTypes`: Slice of allowed file MIME types for uploads. **_Type:_** []string | **_e.g._** []string{"image/jpeg", "image/png", "image/gif"} | **_default:_** all types
* `MaxJSONSize`: Maximum allowed JSON size. **_Type:_** int | **_e.g._** 1024 * 1024 | **_default:_** 1MB
* `AllowUnknownFields`: Whether to allow unknown fields in JSON parsing. **_Type:_** bool | **_e.g._** true | **_default:_** false

## Types 🚧

### UploadedFile

```go
type UploadedFile struct {
    NewFileName      string
    OriginalFileName string
    FileSize         int64
}
```

Represents information about an uploaded file.

### JSONResponse

```go
type JSONResponse struct {
    Error   bool        `json:"error"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

Represents a standardized JSON response structure.

### Error Handling

Most functions return errors that should be checked and handled appropriately in your application.

## Contributing ⭐

If you find any issues or have suggestions for improvement, please feel free to open an issue or create a pull request in my GitHub Repository. We welcome contributions from the community to make this package even better.

### Steps to contribute:

1. Fork the repository.
2. Create a new branch for your contribution.
3. Make your changes and commit them.
4. Push the changes to your fork.
5. Open a pull request with a detailed description of your contribution.
