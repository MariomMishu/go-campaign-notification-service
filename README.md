# Go Campaign Notification Service

## Overview

This Go-based service manages campaign notifications using background task queues and robust authentication mechanisms to ensure secure and efficient delivery of promotional emails.


## User Roles and Permissions

- **Admin (Role ID: 1)**  
  Has full authority to approve campaigns and promotional offers created by Managers or themselves.

- **Manager (Role ID: 2)**  
  Can create campaigns or promotional offers but requires Admin approval before execution.

- **Customer (Role ID: 3)**  
  Receives approved campaign and promotional emails.

---

## Campaign Workflow

1. **Campaign Creation:**  
   Users with Manager or Admin roles can create new campaigns or promotions.

2. **Campaign Approval:**  
   Admin users review and approve submitted campaigns.

3. **Email Notification:**  
   Upon approval, promotional emails are sent asynchronously to customers via background task queues.

---

## Key Features

- Role-based access control to ensure proper authorization.  
- Asynchronous email dispatch using background workers for scalable performance.  
- Clear separation of responsibilities between Managers and Admins.


## Project Structure
go-campaign-notification-service/
├── cmd/ # Application entry points
├── config/ # Configuration logic
├── conn/ # Database connection setup
├── domain/ # Repository and service interface definitions
├── handlers(controllers) # HTTP request handlers (controllers)
├── middlewares/ # Middleware functions for request processing
├── models/ # Data models and structures
├── repositories/ # Database interaction logic
├── routes/ # API route definitions
├── server/ # Server setup and initialization
├── services/ # Core business logic
├── types/ # Shared type definitions
├── utils/ # Shared utility functions
├── dependency # DB queries and Postman collections
├── go.mod # Go module definition
├── go.sum # Dependency checksum file
└── README.md # Project documentation

---

## Tech Stack

- **Language:** Golang  
- **Database:** MySQL  
- **Cache:** Redis  
- **Queue:** Asynq  


## Libraries Used

- [Cobra](https://github.com/spf13/cobra) - Framework for building CLI applications  
- [GORM](https://gorm.io/) - ORM library for interacting with relational databases  
- [Echo](https://echo.labstack.com/) - Web framework for building RESTful APIs and HTTP routing  
- [Ozzo-Validation](https://github.com/go-ozzo/ozzo-validation) - Input validation library  
- [Asynq](https://github.com/hibiken/asynq) - Background task and distributed queue library  

---

## Authentication Overview

### User Signup

Supports secure user registration flow.

### User Login & Logout

Manages user session lifecycle securely.

### Token Management

- **Access Tokens:**  
  - JWT-based authentication  
  - Short-lived tokens for secure access

- **Refresh Tokens:**  
  - Session persistence  
  - Long-lived tokens for maintaining sessions

### Session Management

- Uses Redis for fast in-memory storage  
- Efficiently handles frequent reads/writes  
- Manages token validation and session expiry  

### UUID Mapping

- Tokens are mapped with token’s unique identifiers (UUID)  
- Provides improved security and control over user sessions  

## Getting Started command 

 go run main.go serve



