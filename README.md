# Go HTML App

This project demonstrates a basic Golang web application that uses:
- **Echo** as the web framework.
- **HTMX** as lightweight library that enhances HTML with AJAX, CSS transitions, and server-sent events, enabling dynamic, interactive web applications without heavy JavaScript.
- **Templ** for HTML templating.
- **TailwindCSS** for styling.

## Features

- Simple server setup using Echo.
- HTMX to update data via AJAX call
- Dynamic HTML rendering with templ.
- Modern UI design with TailwindCSS.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/)
- Node.js (for managing TailwindCSS)
- [TailwindCSS](https://tailwindcss.com/docs/installation)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/username/go-html-app.git
   cd go-html-app
   ```

2. Install Golang dependencies:

   ```bash
   go mod tidy
   ```

3. Install Node dependencies for TailwindCSS:

   ```bash
   npm install
   ```

### Running the Application

#### Build and run the Go application:

   ```make live```


## Project Structure

```
.
├── main.go          # Main application entry point.
├── templates/       # Contains HTML templates.
├── static/          # Static files including TailwindCSS output.
└── README.md        # Project documentation.
```

## Technologies Used

- [Echo](https://echo.labstack.com/) - Fast and minimalist web framework.
- [Templ](https://github.com/yourusername/templ) - Template engine for HTML rendering.
- [TailwindCSS](https://tailwindcss.com/) - Utility-first CSS framework.

## License

This project is licensed under the MIT License.

Happy coding!