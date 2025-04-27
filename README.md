# WebHookDelivery

**WebHookDelivery** is a robust and scalable service designed to handle webhook deliveries. The system listens for webhook requests, queues them for delivery, retries failed deliveries, logs delivery statuses, and maintains the system clean by automatically removing old logs. It is fully dockerized for easy setup and deployment.

---

## ✨ Key Features
- **Subscription Creation**: Create a subscription by providing a subscription ID, target URL, and event types you want to listen to.
- **Webhook Task Processing**: Only webhooks matching the subscription's event types are queued for delivery.
- **Asynchronous Processing**: Webhook tasks are handled asynchronously using Go goroutines and Redis queues.
- **Retry Logic**: Failed deliveries are retried with an exponential backoff strategy.
- **Delivery Logging**: Every attempt (success or failure) is logged into a PostgreSQL database.
- **Auto Cleaning Logs**: Logs older than 72 hours are automatically deleted to keep the database clean.
- **Fetch Webhook Logs**: Retrieve delivery logs for subscriptions — only logs within the last 72 hours are retained.
- **Event Type Check**: Only webhook payloads with an event type matching the subscription are processed.
- **Signature Verification**: *(Not implemented yet)* — current validation is based only on event type matching.
- **Fully Dockerized**: Docker and Docker Compose setup for easy local development.

---

## 🛠 Tech Stack
- **Backend**: Golang
- **Frontend**: Next.js
- **Database**: PostgreSQL
- **Queue**: Redis
- **Containerization**: Docker, Docker Compose

---

## ⚙️ How the System Works

1. **Subscription Creation**:
   - A user creates a subscription by providing:
     - Subscription ID
     - Target URL
     - Event types they are interested in (e.g., `user.created`, `order.completed`).

2. **Receiving Webhook Requests**:
   - A webhook request must include the event type.
   - If the incoming webhook's event type matches a subscription's registered event types, the webhook is queued for delivery.

3. **Queueing and Delivery**:
   - Webhook payloads are pushed into a Redis queue.
   - Background workers (using Go goroutines) continuously pull from the queue and attempt to deliver the payload to the target URL.

4. **Retry on Failure**:
   - If a delivery fails, the worker retries sending the webhook with exponential backoff until a maximum retry limit is reached.

5. **Logging Deliveries**:
   - Each delivery attempt (successful or failed) is logged into PostgreSQL with details like HTTP status code, timestamps, and number of attempts.

6. **Log Cleanup**:
   - A separate background goroutine regularly checks the logs and deletes logs older than **72 hours** automatically.

7. **Fetching Logs**:
   - Users can query the backend to fetch webhook logs for their subscription.
   - Only logs within the last 72 hours are available.

---
Deployed on same ews ec2 server.

Frontend Link: http://13.51.170.153:3000/subscriptions
Backend Link: http://13.51.170.153:8080
DocumentationLink: https://documenter.getpostman.com/view/34442065/2sB2j1gBvt
In frontend everything can be accessed except the ingest for post.
For that use documentation link.

## 🚀 Running the Project Locally

### Prerequisites
- Docker and Docker Compose installed
  - [Install Docker](https://docs.docker.com/get-docker/)
  - [Install Docker Compose](https://docs.docker.com/compose/install/)

### Clone the Repository
```bash
git clone https://github.com/harshgupta9473/WebHookDelivery.git
cd WebHookDelivery

Create a .env file in the root directory:
DB_HOST=postgres
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=webhookdb
REDIS_HOST=redis
REDIS_PORT=6379
to run locally also load env using godotenv as it is commented for now to run locally without docker.
To run using docker run docker-compose up build -d or docker compose up -d
Will need to change the Host to localhost for running locally.

Access the Application
Frontend: http://localhost:3000

Backend API: http://localhost:8080


