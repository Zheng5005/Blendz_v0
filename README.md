# Blendz

A real-time language exchange platform that connects people worldwide to practice and learn languages together. Blendz enables users to find language partners, chat in real-time, and make video calls — all in one place.

## Features

- **User Authentication** — Sign up, login, and secure session management with JWT and HTTP-only cookies
- **Onboarding Flow** — Profile setup with native language, learning language, bio, and location
- **Friend System** — Send, accept, and manage friend requests with incoming/outgoing request tracking
- **User Discovery** — Get recommended language exchange partners based on your profile
- **Real-time Chat** — Powered by Stream Chat for instant messaging between friends
- **Video Calls** — Integrated Stream Video for face-to-face language practice sessions
- **Notifications** — Track friend requests and social activity
- **Theme Support** — Customizable UI themes with DaisyUI

## Tech Stack

### Frontend (`Client/`)

| Technology | Purpose |
|---|---|
| React 19 + TypeScript | UI framework |
| Vite | Build tool and dev server |
| React Router 7 | Client-side routing |
| TanStack Query | Server state management |
| Zustand | Client state management |
| Tailwind CSS 4 + DaisyUI | Styling and component library |
| Stream Chat React SDK | Real-time chat components |
| Stream Video React SDK | Video call integration |
| Axios | HTTP client |
| Lucide React | Icon library |
| React Hot Toast | Notifications |

### Backend (`Server/`)

| Technology | Purpose |
|---|---|
| Go 1.24 | Server runtime |
| MongoDB (mongo-driver v2) | Database |
| JWT (golang-jwt) | Authentication tokens |
| Stream Chat Go SDK | Chat backend integration |
| godotenv | Environment configuration |

## Project Structure

```
Blendz_v0/
├── Client/                 # React frontend application
│   ├── src/
│   │   ├── components/     # Reusable UI components
│   │   ├── hooks/          # Custom React hooks (auth, logout, etc.)
│   │   ├── lib/            # API client, axios instance, utilities
│   │   ├── pages/          # Route-level page components
│   │   ├── store/          # Zustand stores (theme, etc.)
│   │   └── constants/      # Application constants
│   └── ...
├── Server/                 # Go backend API
│   ├── auth/               # Authentication handlers (signup, login, logout)
│   ├── chat/               # Stream chat token generation
│   ├── db/                 # MongoDB connection and collections
│   ├── middleware/         # CORS and route protection
│   ├── models/             # Data models and database operations
│   ├── stream/             # Stream SDK initialization and helpers
│   ├── utils/              # JWT utilities and password hashing
│   └── main.go             # Application entry point and route definitions
└── TODOS.txt               # Pending improvements
```

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+
- MongoDB (local or Atlas)
- Stream API credentials (chat + video)

### Backend Setup

1. Navigate to the server directory:

```bash
cd Server
```

2. Create a `.env` file:

```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
STREAM_API_KEY=your_stream_api_key
STREAM_API_SECRET=your_stream_api_secret
JWT_SECRET=your_jwt_secret
```

3. Install dependencies and run:

```bash
go mod tidy
go run main.go
```

The server starts at `http://localhost:8080`.

### Frontend Setup

1. Navigate to the client directory:

```bash
cd Client
```

2. Install dependencies:

```bash
npm install
```

3. Create a `.env` file:

```env
VITE_STREAM_API_KEY=your_stream_api_key
```

4. Start the development server:

```bash
npm run dev
```

The application opens at `http://localhost:5173`.

## API Endpoints

### Authentication

| Method | Path | Description |
|---|---|---|
| POST | `/api/auth/signup` | Register a new user |
| POST | `/api/auth/login` | Authenticate and create session |
| POST | `/api/auth/logout` | End current session |
| POST | `/api/auth/onboarding` | Complete user profile setup |

### Users

| Method | Path | Description |
|---|---|---|
| GET | `/api/users` | Get recommended users |
| GET | `/api/users/me` | Get authenticated user info |
| GET | `/api/users/friends` | Get user's friend list |
| POST | `/api/users/friend-request/{id}` | Send a friend request |
| PUT | `/api/users/friend-request/{id}/accept` | Accept a friend request |
| GET | `/api/users/friend-requests` | Get incoming friend requests |
| GET | `/api/users/outgoing-friend-requests` | Get outgoing friend requests |

### Chat

| Method | Path | Description |
|---|---|---|
| GET | `/api/chat/token` | Get Stream chat authentication token |

## Routes

| Path | Description | Auth Required |
|---|---|---|
| `/` | Home — friends list and recommendations | Yes |
| `/login` | Login page | No |
| `/signup` | Registration page | No |
| `/onboarding` | Profile setup flow | Yes (first login) |
| `/notifications` | Friend request notifications | Yes |
| `/chat/:id` | Real-time chat with a friend | Yes |
| `/call/:id` | Video call with a friend | Yes |

## License

This project is open source and available under the MIT License.
