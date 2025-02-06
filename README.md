# README.md

# Go HTML App

This project is a simple Go web application that displays information about the host where the backend is running. It utilizes Go's HTML templating capabilities to dynamically render host details.

## Project Structure

```
go-html-app
├── main.go         # Entry point of the application
├── go.mod          # Module definition and dependencies
├── templates       # Directory containing HTML templates
│   └── host.html   # Template for displaying host information
└── README.md       # Project documentation
```

## Requirements

- Go 1.16 or later

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd go-html-app
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Running the Application

To run the application, execute the following command in the project directory:

```
go run main.go
```

The application will start an HTTP server on `localhost:8080`. You can access it by navigating to `http://localhost:8080` in your web browser.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.