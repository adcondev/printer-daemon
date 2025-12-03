# Learning: Printer Daemon

This document serves as a technical summary of the `printer-daemon` project, highlighting key technologies, architectural decisions, and skills demonstrated. It is intended for personal reference and CV updates.

## Project Overview
The **Printer Daemon** is a specialized WebSocket server designed to act as a bridge between web-based Point of Sale (POS) applications and physical ticket printers. It solves the limitation of browser-based printing by allowing direct hardware control via a local WebSocket connection. The system accepts JSON payloads to configure printers, define ticket templates, and execute print jobs using ESC/POS commands.

## Tech Stack and Key Technologies
*   **Language:** Go (Golang)
*   **Communication:** WebSockets (Real-time bidirectional communication)
*   **Hardware Integration:** ESC/POS Protocol (via custom library)
*   **Operating System:** Windows (Target deployment, utilizing Windows Print Spooler)
*   **Build Automation:** Taskfile (go-task)
*   **Data Serialization:** JSON

## Notable Libraries
*   **`github.com/gorilla/websocket`**: Used to implement the WebSocket server. It was chosen for its stability, performance, and standard compliance, enabling reliable real-time communication with web clients.
*   **`github.com/adcondev/pos-printer`**: A custom/internal library used to abstract the complexities of the ESC/POS protocol. It handles the low-level byte commands for text formatting, cutting paper, and opening cash drawers.
*   **`encoding/json`**: Standard Go library used extensively for parsing complex configuration and ticket data structures received from clients.

## Major Achievements and Skills Demonstrated

### Backend Development & Architecture
*   **Designed a Robust WebSocket Server:** Implemented a concurrent WebSocket server capable of handling multiple client connections (with thread-safe client management) and graceful shutdown.
*   **Hardware Abstraction Layer:** Created a clean service layer (`TicketPrinter`) that decouples the printing logic from the transport layer (WebSockets), making the code testable and maintainable.
*   **Custom Protocol Design:** Defined a structured JSON-based protocol for client-server communication, supporting distinct message types (`config`, `template`, `ticket`) to manage the printing lifecycle.

### Hardware Integration & Optimization
*   **Solved Windows Print Spooling Issues:** Implemented a "Force Flush" mechanism by strategically closing and reopening printer connections to bypass Windows print spooler buffering, ensuring tickets print immediately rather than getting stuck in the queue.
*   **Dynamic Template System:** Developed a flexible templating engine that allows clients to customize ticket layouts (e.g., hiding specific fields, changing headers/footers) without modifying the backend code.
*   **Adaptive Layouts:** Implemented logic to automatically adjust text formatting and column widths based on paper size (58mm vs. 80mm).

### DevOps & Tooling
*   **Build Automation:** Integrated `Taskfile` to streamline development workflows, including running tests, building the binary, and starting the server.

## Skills Gained/Reinforced
*   **Go Programming:** Advanced usage of Goroutines, Channels, Mutexes, and Struct tags.
*   **System Programming:** Interaction with OS-level resources (printers) and signal handling (SIGTERM/SIGINT).
*   **API Design:** Designing event-driven APIs over WebSockets.
*   **Protocol Implementation:** Working with raw byte streams and ESC/POS commands.
*   **Problem Solving:** Debugging and resolving hardware-specific constraints (Windows printing subsystem).
