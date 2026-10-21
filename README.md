# 🚀 Go URL Shortener & Click Tracker API

A lightweight, concurrent REST API built with **Go (1.22+)** and containerized with a multi-stage **Dockerfile**.

## ✨ Key Features
- **Zero External Dependencies:** Uses Go's modern `net/http` standard library with method and wildcard routing (`GET /{code}`).
- **Thread-Safe Concurrency:** Uses `sync.RWMutex` to safely handle concurrent read/write requests across Goroutines without race conditions.
- **Multi-Stage Docker Build:** Compiles a static Go binary and packages it into an ultra-lightweight Alpine Linux container (<20MB).

## 🛠️ Quick Start (Docker)
```bash
docker build -t go-shortener .
docker run -p 8080:8080 go-shortener

GET /health — Health check endpoint

    What it does: When someone visits /health, your server replies with {"status": "healthy"}. Cloud services use this to verify your app is running properly.

POST /shorten — Create a new short link

    What it does: Allows a user to send a JSON payload (like {"short_code": "hireme", "original_url": "[https://go.dev](https://go.dev)"}) to your server to instantly create a new custom redirect link.

GET /{code} — Redirects to original URL and increments click counter

    What it does: When someone types http://localhost:8080/hireme into their browser, your server looks up hireme, adds +1 to the click counter safely using your mutex, and redirects them straight to [https://go.dev](https://go.dev).

GET /stats/{code} — Returns live JSON click analytics

    What it does: When you visit http://localhost:8080/stats/hireme, your server outputs the JSON analytics showing how many times that specific link has been clicked.