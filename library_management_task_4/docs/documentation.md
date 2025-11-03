# Library Management System (Go)

## Overview
This console-based system manages books and members using Go’s structs, interfaces, maps, and slices.

### Features
- Add / Remove books  
- Borrow / Return books  
- List available and borrowed books  

### Architecture
- **models/** — Struct definitions  
- **services/** — Business logic & interface implementation  
- **controllers/** — Console interaction logic  
- **main.go** — Entry point  

### Data Structures
- `map[int]Book` — Stores all books  
- `map[int]Member` — Stores all members  
- `[]Book` — Tracks member’s borrowed books  

### Usage
Run:
```bash
go run main.go
