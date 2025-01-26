# Go Web Server with Static Files and Form Handling

This is a simple Go web server that serves static files and handles form submissions. The server serves an `index.html` page and a `form.html` page for submitting user data (name and address) through a form.

## Features

- Serve static HTML files (e.g., `index.html`, `form.html`) from the `static` directory.
- Handle `POST` form submissions and display the submitted data.
- Simple Go HTTP server using the `net/http` package.

## Endpoints



```
### `GET /`

- Serves the `index.html` page located in the `static` directory.

### `GET /hello`

- A simple endpoint that returns a "Hello" message.

### `GET /form.html`

- Serves the `form.html` page located in the `static` directory.

### `POST /form`

- Handles form submissions from the `form.html` page.
- Displays the submitted name and address.

