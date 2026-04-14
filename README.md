# Medical Scheduling Platform (gRPC Migration)

This project is the second assignment for the Advanced Programming 2 course. We have migrated the Medical Scheduling Platform from a REST-based architecture to gRPC. The project consists of two microservices that communicate using the gRPC protocol and Protocol Buffers.

## Project Structure
- **Doctor Service**: Manages doctor profiles. It runs on port `:50051`.
- **Appointment Service**: Manages medical appointments. It runs on port `:50052`.
- **Inter-service communication**: The Appointment Service acts as a gRPC client. It calls the Doctor Service to check if a doctor exists before creating an appointment.

## How to Regenerate Proto Stubs
If you change the `.proto` files, use these commands from the root directory to update the Go code:

**Doctor Service:**
```bash
protoc --go_out=. --go-opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    doc/proto/doctor.proto
```

**Appointment Service:**
```bash
protoc --go_out=. --go-opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    appointment/proto/appointment.proto
```

## How to Run
1. **Start the Doctor Service**:
   ```bash
   cd doc
   go run cmd/main.go
   ```
2. **Start the Appointment Service**:
   ```bash
   cd appointment
   go run cmd/main.go
   ```

## Testing the Project

**1. Create a Doctor:**
```bash
grpcurl -plaintext -import-path ./doc/proto -proto doctor.proto \
    -d '{"name": "Gregory House", "specialty": "Diagnostics", "email": "house@med.com"}' \
    localhost:50051 doctor.DoctorService/CreateDoctor
```

**2. Create an Appointment:**
```bash
grpcurl -plaintext -import-path ./appointment/proto -proto appointment.proto \
    -d '{"patient_id": "P-001", "doctor_id": "PASTE_DOCTOR_ID_HERE", "date": "2026-05-10"}' \
    localhost:50052 appointment.AppointmentService/CreateAppointment
```

## Error Handling and Resilience
- **NotFound**: If you try to create an appointment with a doctor ID that does not exist, the service returns a `codes.NotFound` error.
- **Unavailable**: If the Doctor Service is offline, the Appointment Service returns a `codes.Unavailable` error. This proves that the system handles service failures gracefully without crashing.

## REST vs gRPC Trade-offs
Main differences:

1. **Protocol and Speed**: REST uses HTTP/1.1 and JSON (text format), which is slower to parse. gRPC uses HTTP/2 and Protocol Buffers (binary format), which is much faster and uses less network bandwidth.
2. **Strict Contracts**: In REST, documentation (like Swagger) is optional. In gRPC, the `.proto` file is a strict contract. You cannot send data that doesn't follow the schema, which reduces bugs.
3. **Multiplexing**: gRPC (via HTTP/2) allows sending multiple requests over a single connection at the same time. REST (via HTTP/1.1) usually handles one request per connection, making it less efficient for many microservices.