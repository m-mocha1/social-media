# SocialSphere

A Facebook-like social network built with a Go & SQLite backend and a Vue.js frontend, containerized with Docker.



## Features

- **User Accounts & Profiles**  
  - Registration & login (email, password, name, DOB, avatar, nickname, “about me”)  
  - Public/private profiles toggle  
- **Followers**  
  - Send/accept/decline follow requests for private profiles  
  - Auto-follow for public profiles  
- **Posts & Comments**  
  - Create text posts with optional JPEG/PNG/GIF attachments  
  - Privacy levels: public, followers-only, custom private  
- **Groups & Events**  
  - Create/search/join groups by invite or request  
  - Group-only posts & comments  
  - Group events with “Going/Not Going” RSVPs  
- **Chats & Notifications**  
  - One-to-one & group chat via WebSockets (private messages & emojis)  
  - Real-time notifications for follows, invites, requests, events, etc.

---

## Tech Stack

- **Backend:** Go  
  - HTTP server & middleware  
  - Gorilla WebSocket  
  - Sessions & cookies for auth  
  - SQLite (via [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3))  
  - Migrations with [`golang-migrate`](https://github.com/golang-migrate/migrate)  
- **Frontend:** Vue 3 + Vite  
  - Composition API, Vue Router, Pinia (state)  
- **Containerization:** Docker & Docker Compose

---

## Prerequisites

- **Go** ≥ 1.20  
- **Node.js & npm** (or Yarn)  
- **Docker & Docker Compose** (for containerized setup)

---

## Installation & Setup

### 1. Clone the Repository

git clone https://github.com/your-username/social-media.git
cd socialsphere
2. Backend Setup
Navigate to backend folder


cd backend
Install Go dependencies


go mod download
Run database migrations
Ensure your $PWD contains pkg/db/migrations/sqlite:


go run server.go migrate
Start the server
By default listens on :8080:


go run server.go
3. Frontend Setup
Navigate to frontend folder


cd ../frontend
Install npm dependencies


npm install
Run dev server


npm run dev
Opens at http://localhost:5173 by default.

Docker
You can also run both services via Docker Compose:


docker-compose up --build
Frontend: served on http://localhost:3000

Backend: API at http://localhost:8080

Ports are configurable in docker-compose.yml.

Usage
Register a new account.

Log in and explore:
![Screenshot 2025-05-16 165033](https://github.com/user-attachments/assets/6ea2e51b-7a25-4c84-a679-ac823247fb90)

Update your profile (public/private).
![Uploading Screenshot 2025-05-16 165135.png…]()

Follow other users or accept follow requests.

Create posts, comment, upload images/GIFs.
![Screenshot 2025-05-16 165117](https://github.com/user-attachments/assets/270282b3-8a83-49d8-beeb-cb8348f9ce6c)

Join or create groups and events.
![Screenshot 2025-05-16 165125](https://github.com/user-attachments/assets/e5cfb99b-549a-48e6-b486-f3ef280a5945)

Chat privately or in groups; see real-time notifications.

Authentication & Security
Passwords hashed with bcrypt.

Sessions stored server-side; session ID in an HTTP-only cookie.

CSRF protection via double-submit cookie.

WebSockets & Real-time
Gorilla WebSocket hub in backend/pkg/ws/ handles:

Private chat channels

Group chat rooms

Real-time notifications







