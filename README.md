# Movie Streaming Web 🎬

A full-featured movie streaming backend API built with Go, Gin framework, and MongoDB. This application provides user authentication, movie management, and personalized movie recommendations based on user preferences.

## 📋 Features

- **User Authentication & Authorization**
  - User registration and login with JWT tokens
  - Role-based access control (ADMIN/USER)
  - Secure password hashing
  - Token refresh mechanism

- **Movie Management**
  - Browse all movies
  - Get detailed movie information
  - Add new movies (Admin only)
  - Update movie reviews (Admin only)
  - Movie genre categorization
  - Movie ranking system (1-10 scale)

- **Personalized Recommendations**
  - AI-powered movie recommendations based on user's favorite genres
  - Integrated with LangChain for intelligent suggestions

- **CORS Support**
  - Configurable allowed origins
  - Support for cross-origin requests

## 🛠️ Tech Stack

- **Language:** Go 1.24
- **Framework:** Gin Web Framework
- **Database:** MongoDB
- **Authentication:** JWT (JSON Web Tokens)
- **AI/ML:** LangChain Go
- **Environment Management:** godotenv

### Key Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `go.mongodb.org/mongo-driver/v2` - MongoDB driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `github.com/tmc/langchaingo` - LangChain integration
- `golang.org/x/crypto` - Password hashing

## 📁 Project Structure

```
Movie-Streaming-web/
├── Server/
│   └── MovieServer/
│       ├── controllers/       # Request handlers
│       │   ├── movie_controller.go
│       │   └── user_controller.go
│       ├── database/          # Database connection
│       │   └── database_connection.go
│       ├── middleware/        # Authentication middleware
│       │   └── auth_middleware.go
│       ├── models/            # Data models
│       │   ├── movie_models.go
│       │   └── user_model.go
│       ├── routes/            # Route definitions
│       │   ├── protected_routes.go
│       │   └── unprotected_routes.go
│       ├── utils/             # Utility functions
│       │   └── token_utils.go
│       ├── main.go            # Application entry point
│       ├── go.mod             # Go module dependencies
│       └── go.sum             # Dependency checksums
├── API_TEST.http              # HTTP client test file
└── README.md                  # Project documentation
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- MongoDB instance (local or cloud)
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Rudraksh121a/Movie-Streaming-web.git
   cd Movie-Streaming-web/Server/MovieServer
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the `Server/MovieServer` directory:
   ```env
   MONGODB_URI=your_mongodb_connection_string
   SECRET_KEY=your_jwt_secret_key
   ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

## 🔑 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | Required |
| `SECRET_KEY` | JWT secret key for token generation | Required |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins | `http://localhost:5173` |

## 📡 API Endpoints

### Public Endpoints (No Authentication Required)

#### Get All Movies
```http
GET /movies
```
Returns a list of all available movies.

#### User Registration
```http
POST /register
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "SecurePass123",
  "role": "USER",
  "favourite_genres": [
    {
      "genre_id": 1,
      "genre_name": "Comedy"
    }
  ]
}
```

#### User Login
```http
POST /login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```
Returns JWT tokens for authentication.

### Protected Endpoints (Require Authentication)

Include the JWT token in the Authorization header:
```http
Authorization: Bearer <your_jwt_token>
```

#### Get Movie Details
```http
GET /movie/:imdb_id
```
Get detailed information about a specific movie.

#### Add Movie (Admin Only)
```http
POST /addmovie
Content-Type: application/json

{
  "imdb_id": "tt0102034",
  "title": "Movie Title",
  "poster_path": "https://image.tmdb.org/t/p/original/poster.jpg",
  "youtube_id": "videoID",
  "genre": [
    {
      "genre_id": 6,
      "genre_name": "Sci-fi"
    }
  ],
  "admin_review": "Great movie!",
  "ranking": {
    "ranking_value": 8,
    "ranking_name": "excellent"
  }
}
```

#### Get Recommended Movies
```http
GET /recommendedmovies
```
Get personalized movie recommendations based on user's favorite genres.

#### Update Movie Review (Admin Only)
```http
PATCH /updatereview/:imdb_id
Content-Type: application/json

{
  "admin_review": "Updated review content"
}
```

#### Health Check
```http
GET /health
```
Returns server health status.

## 🧪 Testing

Use the included `API_TEST.http` file with REST Client extensions (VS Code, IntelliJ) to test the API endpoints.

## 🔒 Security

- Passwords are hashed using bcrypt
- JWT tokens for secure authentication
- Role-based access control
- CORS configuration for allowed origins
- Input validation on all endpoints

## 📝 Data Models

### User Model
- User ID, Name, Email
- Password (hashed)
- Role (ADMIN/USER)
- Favorite genres
- JWT tokens

### Movie Model
- IMDB ID
- Title
- Poster path (URL)
- YouTube trailer ID
- Genres
- Admin review
- Ranking (1-10 scale)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 👤 Author

**Rudraksh121a**
- GitHub: [@Rudraksh121a](https://github.com/Rudraksh121a)

## 🙏 Acknowledgments

- Gin Web Framework for the excellent HTTP framework
- MongoDB for the database solution
- LangChain for AI-powered recommendations
