# ForumBook

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Usage](#usage)
- [Authentication](#authentication)
- [Post Management](#post-management)
- [Likes and Dislikes](#likes-and-dislikes)
- [Filtering](#filtering)
- [Image Upload](#image-upload)
- [Error Handling](#error-handling)
- [Docker](#docker)
- [Best Practices](#best-practices)
- [Learning Outcomes](#learning-outcomes)

## Overview
ForumBook is a web-based forum that enables users to communicate through posts and comments, categorize content, and engage with posts via likes and dislikes. The project is built using Go and SQLite, ensuring efficient data management and authentication through cookies.

## Features
- User registration and authentication (with optional password encryption).
- Session management using cookies.
- Post creation, categorization, and commenting.
- Like and dislike system for posts and comments.
- Filtering posts by categories, user-created posts, and liked posts.
- Image upload support for JPEG, PNG, and GIF (max 20MB).
- Docker containerization for easy deployment.

## Technologies Used
- **Backend:** Go
- **Database:** SQLite
- **Authentication:** bcrypt (for password hashing), cookies
- **UUID:** For unique identifiers
- **Docker:** For containerization

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/anlazaar/my_forum.git
   cd forumbook
   ```
2. Build and run using Docker:
   ```sh
   docker-compose up --build
   ```
3. Alternatively, run manually:
   ```sh
   go run cmd/main.go
   ```

## Usage
Once the server is running:
- Navigate to `http://localhost:8080` in your browser.
- Register an account and log in to start posting.
- Explore the forum, filter posts, like/dislike content, and comment.

## Authentication
- Users must register with a unique email, username, and password.
- Only registered users can create posts, comment, and interact.
- Sessions are managed using cookies with an expiration date.
- Optional: Passwords can be hashed using bcrypt.

## Post Management
- Users can create posts with text and images (JPEG, PNG, GIF up to 20MB).
- Posts can be assigned one or more categories.
- Both registered and guest users can view posts.

## Likes and Dislikes
- Registered users can like or dislike posts and comments.
- The number of likes and dislikes is visible to all users.

## Filtering
Users can filter posts by:
- **Categories:** View posts based on topics.
- **Created Posts:** View only their own posts.
- **Liked Posts:** View posts they have liked.

## Image Upload
- Registered users can attach images to posts.
- Supported formats: JPEG, PNG, GIF.
- Images exceeding 20MB will be rejected with an error message.

## Error Handling
- Proper HTTP status codes are returned for errors.
- User-friendly error messages for authentication failures, invalid input, and file size limits.
- Technical errors are logged for debugging.

## Docker
- The project is containerized for easy deployment.
- Run the forum with a single command using `docker-compose up`.
- Database persistence is ensured through volume management.

## Best Practices
- Proper error handling and validation.
- Clean code structure with separation of concerns.
- Secure session and password management.
- Unit testing for core functionalities.

## Learning Outcomes
Through this project, you will gain knowledge in:
- Web development fundamentals (HTML, HTTP, Sessions, Cookies).
- Database management with SQLite and SQL queries.
- Authentication and security best practices.
- Image processing and handling in web applications.
- Containerization using Docker.

---

This project serves as an excellent opportunity to apply and enhance your skills in full-stack Go web development!
