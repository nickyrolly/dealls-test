# Dating App API Documentation

## Overview
This document provides comprehensive documentation for the Dating App API, which enables users to create profiles, match with other users, and manage their dating preferences.

## Authentication
The API uses JWT (JSON Web Token) for authentication. All protected endpoints require a valid JWT token in the Authorization header.

### Authentication Endpoints

#### POST /api/v1/auth/signup
Create a new user account.

**Request Body:**
```json
{
  "email": "string",
  "password": "string",
  "name": "string"
}
```

**Response:**
```json
{
  "token": "string",
  "user": {
    "id": "uuid",
    "email": "string",
    "name": "string"
  }
}
```

#### POST /api/v1/auth/login
Login to an existing account.

**Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "string",
  "user": {
    "id": "uuid",
    "email": "string",
    "name": "string"
  }
}
```

## Profile Management

### Profile Endpoints

#### GET /api/v1/profile
Get the current user's profile.

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "string",
  "bio": "string",
  "age": "integer",
  "gender": "string",
  "location": "string",
  "height": "integer",
  "weight": "integer",
  "occupation": "string",
  "photos": [
    {
      "id": "uuid",
      "url": "string",
      "is_primary": "boolean"
    }
  ]
}
```

#### PUT /api/v1/profile
Update the current user's profile.

**Request Body:**
```json
{
  "name": "string",
  "bio": "string",
  "age": "integer",
  "gender": "string",
  "location": "string",
  "height": "integer",
  "weight": "integer",
  "occupation": "string"
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "string",
  "bio": "string",
  "age": "integer",
  "gender": "string",
  "location": "string",
  "height": "integer",
  "weight": "integer",
  "occupation": "string"
}
```

## Matching System

### Match Endpoints

#### GET /api/v1/matches
Get all matches for the current user.

**Response:**
```json
[
  {
    "id": "uuid",
    "user1_id": "uuid",
    "user2_id": "uuid",
    "matched_at": "datetime",
    "status": "string"
  }
]
```

#### POST /api/v1/profile/{id}/like
Like another user's profile.

**Response:**
```json
{
  "match": {
    "id": "uuid",
    "user1_id": "uuid",
    "user2_id": "uuid",
    "matched_at": "datetime",
    "status": "string"
  }
}
```

#### POST /api/v1/profile/{id}/pass
Pass on another user's profile.

**Response:**
```json
{
  "success": true
}
```

## Error Responses
All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "string"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```
