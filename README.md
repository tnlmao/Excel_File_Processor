# Excel File Management Service

## Overview

This project provides a backend service for managing Excel data through three main functionalities: uploading, viewing, and editing records. It processes Excel files, stores the data in a MySQL database, and uses Redis for caching.

### Key Features

- **Upload**: Upload an Excel file, which is processed to save data in a MySQL table and cached in Redis for 5 minutes.
- **View**: Retrieve records from the Redis cache first; if not available, fallback to the MySQL database.
- **Edit**: Update specific records in both the MySQL database and Redis cache.

## Technologies Used

- **Go**: The programming language used for backend development.
- **MySQL**: For persistent data storage.
- **Redis**: For caching data temporarily.
- **Gin**: A web framework for Go used to handle HTTP requests.

## API Endpoints

### 1. Upload

- **Method**: `POST`
- **Endpoint**: `/upload`
- **Description**: Uploads an Excel file, processes it, and saves its contents to the MySQL database and Redis cache.
- **Request**: An Excel file in the form-data.

### 2. View

- **Method**: `GET`
- **Endpoint**: `/view`
- **Description**: Retrieves records, checking the Redis cache first; if the data is not in the cache, it retrieves from the MySQL database.

### 3. Edit

- **Method**: `PATCH`
- **Endpoint**: `/edit`
- **Description**: Updates a specific record in the database and Redis.
- **Request Body**: JSON object with the record details, including the ID and fields to be updated.

## Installation

1. **Clone the repository**:
   git clone https://github.com/tnlmao/Excel_File_Processor.git
   cd Go_Assignment/

2. **Install Dependencies**:
   go mod tidy

3. **Setup Database**:
    Create a MySQL database and configure the connection in your code.

4. **Start the server**
    go run main.go


## Usage

1. Upload an Excel file: Use Postman or a similar tool to send a POST request to /upload with the Excel file.

2. View data: Send a GET request to /view to see the data

3. Edit a record: Send a PUT request to /edit with the JSON payload containing the updated record.