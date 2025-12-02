# mone

**mo**nitor + **ne**twork = **mone**

A lightweight network monitoring desktop application built with Go and Fyne.

## Features

- Network device selection and monitoring
- System tray integration
- Cross-platform desktop support
- Clean and intuitive UI

## Tech Stack

- **Go** - Backend and core functionality
- **Fyne** - Cross-platform GUI framework
- **gopacket** - Network packet capture

## Project Structure

```
cmd/
  ├── services/    # Background services
  ├── types/       # Type definitions
  ├── ui/          # UI components and screens
  └── utils/       # Helper utilities and constants
```

## Running the Application

```bash
go run main.go
``` 