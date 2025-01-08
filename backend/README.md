
# Moody Weather Backend

## Overview

The Moody Weather backend is a containerized Go application that provides the backend functionality for the Moody Weather project. This application handles HTTP requests, integrates with third-party weather APIs, and communicates with ChatGPT to generate fun and engaging responses based on the weather.

## Prerequisites

- Docker installed on your system.

## Building and Running the Application

1. Clone the repository:

   ```bash
   git clone https://github.com/jw81/moody-weather.git
   cd moody-weather/backend
   ```

2. Build the Docker image:

   ```bash
   docker build -t moody-weather-backend .
   ```

3. Run the Docker container:

   ```bash
   docker run -p 8080:8080 moody-weather-backend
   ```

4. Access the application:

   Open your browser and navigate to `http://localhost:8080`. You should see the message:
   
   ```
   Hello, world! This is the weather app backend.
   ```

## Next Steps

- Implement API calls to fetch weather data.
- Add integration with ChatGPT for customized responses.