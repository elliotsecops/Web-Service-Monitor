
# Web Service Monitor

Web Service Monitor is a Go-based tool designed to periodically check the availability and response times of critical web services. It sends alerts if any service is down or if the response time exceeds a specified threshold.

## Features

- **Periodic Monitoring:** Checks the status of web services at regular intervals.
- **Retry Mechanism:** Retries failed requests to ensure transient issues don't trigger false alerts.
- **Response Time Metrics:** Measures and logs the response times of each service.
- **Threshold Alerts:** Sends alerts if the response time exceeds a configurable threshold.
- **Email Notifications:** Sends email alerts for service downtime and slow responses.
- **Secure Configuration:** Uses environment variables for sensitive credentials.

## Prerequisites

- Go (1.16 or higher)
- SMTP server credentials (for email notifications)

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/elliotsecops/web-service-monitor.git
   cd web-service-monitor
   ```

2. Set up environment variables for email credentials:
   ```sh
   export EMAIL_FROM="your-email@example.com"
   export EMAIL_PASSWORD="your-password"
   export SMTP_HOST="smtp.gmail.com"
   export SMTP_PORT="587"
   ```

3. Create a `config.json` file with the configuration:
   ```json
   {
     "websites": [
       "https://www.google.com",
       "https://www.github.com",
       "https://www.example.com"
     ],
     "interval": 60,
     "maxRetries": 3,
     "responseTimeMax": 2000
   }
   ```

4. Build and run the application:
   ```sh
   go build
   ./web-service-monitor
   ```

## Configuration

- **websites:** List of URLs to monitor.
- **interval:** Interval between checks in seconds.
- **maxRetries:** Number of retries before marking a site as down.
- **responseTimeMax:** Maximum acceptable response time in milliseconds.

## Logging

Logs are written to `monitor.log` in the current directory.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.

