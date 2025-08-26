# ChronoPlay Backend Service - Developer Documentation

## Overview
ChronoPlay Backend Service is a Go-based REST API for a card trading and virtual economy game. The service provides user management, card trading, cash transactions, exchange mechanisms, and automated background tasks.

## Tech Stack
- **Language**: Go
- **Framework**: Gin (HTTP web framework)
- **Database**: MongoDB
- **Authentication**: JWT tokens
- **Scheduling**: Cron jobs
- **Email**: Integration for user verification

## Project Structure
```
├── controllers/        # HTTP request handlers
├── services/          # Business logic layer
├── model/            # Data models and repository interfaces
├── dto/              # Data transfer objects
├── routes/           # API route definitions
├── middlewares/      # Authentication and custom middlewares
├── crons/            # Background tasks and scheduled jobs
├── utils/            # Utility functions (email, validation, etc.)
├── config/           # Configuration management
└── database/         # Database connection setup
```

## Features Implemented ✅

### 1. User Authentication & Management
- **User Registration** (`POST /auth/signup`) ✅
- **Email Verification** (`GET /auth/verify`) ✅
- **User Login** (`POST /auth/login`) ✅
- **Get Current User** (`GET /user/user`) ✅
- **Get User by ID** (`GET /user/get_user`) ✅
- **Friend Management**:
  - Add Friend (`PATCH /user/add_friend`) ✅
  - Get Friends (`GET /user/get_friends`) ✅
  - Remove Friend (`PATCH /user/remove_friend`) ✅
- **Admin Functions**:
  - Activate All Users (`PATCH /auth/activate_all_users`) ✅

### 2. Card Management
- **Add New Card** (`POST /card/add`) ✅
- **Get Card Details** (`GET /card/get_card`) ✅
- Cards include: Number, Name, Description, Rarity, Image, Quantity

### 3. Transaction System
- **Cash Transfers** (`POST /transaction/transfer_cash`) ✅
- **Card Transfers** (`POST /transaction/transfer_cards`) ✅
- **Exchange System**:
  - Create Exchange Request (`POST /transaction/exchange`) ✅
  - Get Possible Exchanges (`GET /transaction/get_possible_exchange`) ✅
  - Execute Exchange (`POST /transaction/execute_exchange`) ✅
- **Transaction History** (`GET /transaction/get_transactions`) ✅

### 4. Notification System
- **Get Notifications** (`GET /notification/get_notifications`) ✅
- **Mark as Read** (`PATCH /notification/mark_as_read`) ✅

### 5. Background Tasks (Cron Jobs)
- **Survival Tax System** ✅
  - Runs daily at midnight
  - Deducts 50 cash units from all active users
  - Deactivates users with insufficient funds
  - Sends notifications to affected users

### 6. Security & Middleware
- **JWT Authentication** ✅
- **CORS Support** ✅
- **Request Validation** ✅
- **Custom Context Middleware** ✅

## API Endpoints

### Authentication Routes (`/auth/`)
```
POST   /auth/signup              # User registration
GET    /auth/verify              # Email verification
POST   /auth/login               # User login
PATCH  /auth/activate_all_users  # Admin: activate all users
```

### User Routes (`/user/`) - Protected
```
GET    /user/user               # Get current user profile
GET    /user/get_user           # Get user by ID
PATCH  /user/add_friend         # Add friend
GET    /user/get_friends        # Get friends list
PATCH  /user/remove_friend      # Remove friend
```

### Card Routes (`/card/`) - Protected
```
POST   /card/add                # Add new card
GET    /card/get_card           # Get card details
```

### Transaction Routes (`/transaction/`) - Protected
```
POST   /transaction/transfer_cash        # Transfer cash between users
POST   /transaction/transfer_cards       # Transfer cards between users
POST   /transaction/exchange             # Create exchange request
GET    /transaction/get_transactions     # Get transaction history
GET    /transaction/get_possible_exchange # Get possible exchanges
POST   /transaction/execute_exchange     # Execute exchange
```

### Notification Routes (`/notification/`) - Protected
```
GET    /notification/get_notifications   # Get user notifications
PATCH  /notification/mark_as_read        # Mark notifications as read
```

## Data Models

### User
- UserID, Name, Email, Username, Password
- Cash balance, Phone number
- Cards owned (with quantities)
- Friend list, User type
- Activation status

### Card
- Card number, Name, Description
- Rarity level, Image URL
- Quantity tracking per user

### Transaction
- Transaction GUID, Type (cash/card/exchange)
- Participants (given_by, given_to)
- Items transferred (cash amounts, cards)
- Timestamp, Status

### Notification
- User ID, Title, Message
- Read status, Timestamp

## Features Not Yet Implemented ❌

### Loan System
The loan controller exists but has no implemented methods:
- Create loan between users
- Get loans owed by user
- Get loans owed to user
- Pay loan installments

### User Management
- User logout functionality
- Password reset mechanism
- User profile updates

### Card Features
- Show all available cards
- Get users holding specific cards
- Update card information
- Card marketplace

### Advanced Features
- Real-time notifications (WebSocket)
- Card trading marketplace
- Auction system
- Leaderboards
- Game statistics

## Getting Started

### Prerequisites
- Go 1.19+
- MongoDB
- Environment variables configured

### Running the Service
```bash
# Install dependencies
go mod download

# Run the service
go run main.go
```

### Environment Variables Required
```
MONGO_DB_NAME=your_db_name
MONGO_URI=your_mongo_connection_string
JWT_SECRET=your_jwt_secret
EMAIL_CONFIG=your_email_settings
```

## Development Notes
- All protected routes require JWT authentication
- Survival tax cron job can be enabled/disabled
- Email verification is required for user activation
- Friend system supports bidirectional relationships
- Exchange system supports complex multi-item trades 
