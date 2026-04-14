# AP2 Assignment 1 – Medical Scheduling Platform

## What This Project Does

This project has two services that work together:

- **Doctor Service** – stores and manages doctor information. Runs on port `8080`.
- **Appointment Service** – stores and manages appointments. Runs on port `8081`.

When someone creates an appointment, the Appointment Service checks with the Doctor Service that the doctor actually exists. If the Doctor Service is down, the appointment is not created.

---

## Architecture

```
┌─────────────────────────────┐         ┌──────────────────────────────────┐
│       Doctor Service        │         │       Appointment Service         │
│         :8080               │         │            :8081                  │
│                             │         │                                   │
│  [Handler]                  │ ◄──HTTP─┤  [Handler]                        │
│  [Use Case]                 │  GET    │  [Use Case]                       │
│  [Repository]               │ /doctors│  [Repository]  [DoctorClient]     │
│  (in-memory)                │  /{id}  │  (in-memory)   (HTTP client)      │
└─────────────────────────────┘         └──────────────────────────────────┘
```

Each service has its own data. They only talk to each other through REST API calls.

---

## What Each Service Does

### Doctor Service
- Creates, stores, and returns doctor profiles.
- Does not know anything about appointments.
- Rules: `full_name` is required, `email` is required and must be unique.

### Appointment Service
- Creates, stores, and returns appointments.
- Before creating an appointment, it calls the Doctor Service to check that the doctor exists.
- Rules: `title` is required, `doctor_id` must point to a real doctor, status can only be `new`, `in_progress`, or `done`. You cannot change status from `done` back to `new`.

---

## Folder Structure

```
service/
├── cmd/
│   └── main.go              ← starts the server
├── internal/
│   ├── model/               ← data types (Doctor, Appointment)
│   ├── usecase/             ← business logic
│   ├── repository/          ← stores data in memory
│   ├── client/              ← (appointment service only) calls Doctor Service
│   ├── transport/http/      ← HTTP handlers
│   └── app/                 ← connects all layers together
└── go.mod
```

Dependencies go in one direction: `handler → use case → repository/client → model`.
Handlers do not contain business logic. Use cases do not import Gin or anything HTTP-related.

---

## How Services Communicate

When a client sends `POST /appointments`, the Appointment Service does this:

```
POST /appointments
  └─► GET http://localhost:8080/doctors/{doctor_id}
        ├─ 200 OK   → create the appointment
        ├─ 404      → return 422: doctor does not exist
        └─ error    → return 503: doctor service unavailable
```

The HTTP client has a 5-second timeout so the service does not wait forever.

---

## How to Run

You need Go 1.22 or newer.

**Terminal 1 – Doctor Service:**
```bash
cd doctor-service
go mod tidy
go run ./cmd
```

**Terminal 2 – Appointment Service:**
```bash
cd appointment-service
go mod tidy
DOCTOR_SERVICE_URL=http://127.0.0.1:8080 go run ./cmd
```

---

## Why There Is No Shared Database

Each service owns its own data. The Appointment Service cannot read the Doctor Service's database directly — it can only ask through the API.

This is important because:
- If the Doctor Service changes its database structure, the Appointment Service still works.
- Each service can be updated or restarted independently.
- A shared database would make this a **distributed monolith** — two separate processes but still tightly connected, which removes most of the benefits of microservices.

---

## Failure Scenario

If the Doctor Service is not running when someone tries to create an appointment:

1. The HTTP client gets a connection error.
2. The use case receives `ErrDoctorServiceUnavailable`.
3. The handler returns `503 Service Unavailable` with a clear error message.
4. The error is logged.
5. The appointment is **not** created.

**In a production system**, you would also add:

| Pattern | Purpose |
|---|---|
| **Timeout** | Already added (5 seconds). Prevents waiting forever. |
| **Retry with backoff** | Try again a few times before giving up, in case the error is temporary. |
| **Circuit breaker** | Stop sending requests to a service that keeps failing. Try again after some time. Useful libraries: `gobreaker`, `hystrix-go`. |

---

## API Examples

### Doctor Service

**Create a doctor**
```
POST http://localhost:8080/doctors
Content-Type: application/json

{
  "full_name": "Dr. Aisha Seitkali",
  "specialization": "Cardiology",
  "email": "a.seitkali@clinic.kz"
}
```

**Get one doctor**
```
GET http://localhost:8080/doctors/{id}
```

**Get all doctors**
```
GET http://localhost:8080/doctors
```

### Appointment Service

**Create an appointment**
```
POST http://localhost:8081/appointments
Content-Type: application/json

{
  "title": "Initial cardiac consultation",
  "description": "Patient referred for palpitations",
  "doctor_id": "<id from doctor service>"
}
```

**Get one appointment**
```
GET http://localhost:8081/appointments/{id}
```

**Get all appointments**
```
GET http://localhost:8081/appointments
```

**Update appointment status**
```
PATCH http://localhost:8081/appointments/{id}/status
Content-Type: application/json

{
  "status": "in_progress"
}
```

Valid values: `new`, `in_progress`, `done`. Changing from `done` to `new` returns `400 Bad Request`.
