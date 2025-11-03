# Library Management System — Concurrency Documentation

## Overview
This version extends the original console library system with concurrent reservation support using Goroutines, Channels and a Mutex.

## Concurrency Components
- **sync.Mutex (`mu`)**: Protects access to shared maps and book state (`Books`, `Members`, `reservations`, `reserveTimers`).
- **Channels & Worker**: `concurrency.StartReservationWorker` returns a channel. Each `ReservationRequest` sent on that channel is handled concurrently by the worker which forwards requests to `Library.ReserveBook`.
- **Goroutines**:
  - The worker spins up a goroutine per incoming reservation request to call `ReserveBook`.
  - `ReserveBook` spawns a goroutine to asynchronously attempt to convert a reservation into a borrow operation (simulating asynchronous processing).
  - A `time.AfterFunc` per reservation automatically cancels a reservation if not borrowed within **5 seconds**.
- **Timers**: `reserveTimers` stores `*time.Timer` for each reservation so the timer can be stopped if borrowing happens before timeout.

## Reservation Workflow
1. Client sends `ReservationRequest` (bookID, memberID) on the worker channel.
2. Worker calls `Library.ReserveBook`.
3. `ReserveBook`:
   - Checks and marks book as `Reserved`.
   - Starts a 5-second timer to auto-unreserve.
   - Starts an async borrow attempt (after a short simulated delay).
4. If `BorrowBook` succeeds before timer fires, timer is stopped and reservation cleared.
5. If timer fires first, reservation is cancelled and book reverts to `Available`.

## How to Test
- Use the console option **"Simulate Concurrent Reservations"** and specify multiple member IDs for the same book. Observe that at most one succeeds in borrowing; others will be rejected or time out.
- The console prints messages from async borrow attempts and auto-cancellation.

## Notes
- All shared state mutations are protected by the mutex to prevent data races.
- This sample is in-memory (no persistent storage). For persistence, add file/DB storage and guard access appropriately.

