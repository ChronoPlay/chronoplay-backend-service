# ChronoPlay Backend Service

A Go-based REST API for a virtual economy and card trading platform that combines collectible card games with social interaction and economic simulation.

## 📚 Documentation

- **[📋 Product Documentation](./Docs/ProductDoc.md)** - Business overview, features, target audience, and product roadmap
- **[⚙️ Developer Documentation](./Docs/DevelopersDoc.md)** - Technical implementation, API endpoints, and development setup 

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- MongoDB
- Environment variables configured

### Installation
```bash
# Clone the repository
git clone https://github.com/ChronoPlay/chronoplay-backend-service.git
cd chronoplay-backend-service

# Install dependencies
go mod download

# Set up environment variables (create .env file)
cp .env.example .env  # Edit with your configuration

# Run the service
go run main.go
```

## 🎮 What is ChronoPlay?

ChronoPlay is a virtual economy platform where players:
- 🃏 **Collect and trade unique cards** with different rarities
- 💰 **Manage virtual currency** for transactions and survival
- 🤝 **Build social networks** to enable trading
- ⏰ **Face daily survival tax** that keeps the economy active
- 🔄 **Execute complex exchanges** combining cards and cash

## ✨ Key Features

### ✅ Implemented
- User authentication with email verification
- Friend networks and social connections
- Card collection and ownership tracking
- Cash and card trading between users
- Complex multi-asset exchange system
- Automated survival tax system (cron jobs)
- Real-time notifications
- Transaction history and audit trails

### 📋 Planned
- Loan system for community support
- Card marketplace and auctions
- Leaderboards and achievements
- Real-time chat integration
- Mobile app development

## 🛠 Tech Stack

- **Backend**: Go with Gin framework
- **Database**: MongoDB
- **Authentication**: JWT tokens
- **Scheduling**: Cron jobs for background tasks
- **Email**: Integration for user verification
- **Security**: CORS, request validation, middleware

## 📖 API Overview

The service provides RESTful endpoints for:

- **Authentication** (`/auth/*`) - Signup, login, verification
- **User Management** (`/user/*`) - Profiles, friends, social features
- **Card System** (`/card/*`) - Card creation and retrieval
- **Transactions** (`/transaction/*`) - Trading, exchanges, history
- **Notifications** (`/notification/*`) - User notifications and updates

> For detailed API documentation, see [Developer Documentation](./Docs/DevelopersDoc.md)

## 🔧 Development

### Project Structure
```
├── controllers/        # HTTP request handlers
├── services/          # Business logic layer
├── model/            # Data models and repositories
├── dto/              # Data transfer objects
├── routes/           # API route definitions
├── middlewares/      # Authentication and middleware
├── crons/            # Background tasks
├── utils/            # Utility functions
├── config/           # Configuration management
└── database/         # Database connection
```

### Environment Variables
```env
MONGO_DB_NAME=your_db_name
MONGO_URI=your_mongo_connection_string
JWT_SECRET=your_jwt_secret
EMAIL_CONFIG=your_email_settings
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Product Documentation](./Docs/ProductDoc.md) - Business and product overview
- [Developer Documentation](./Docs/DevelopersDoc.md) - Technical implementation guide
- [Progress Documentation](./Docs/ProgressDoc.md) - Development progress tracking

---

Built with ❤️ by the ChronoPlay team