# Steam 2FA Code Generator

This Go application generates a Steam 2FA (Two-Factor Authentication) code using the server time from the Steam API. It can be used for Jenkins projects for automated builds and deployments to Steamworks via Steam CLI, and it does not use external dependencies other than the Steam API for security reasons

## Requirements

- Go (version 1.16 or later)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/NeilMooreQ/steam-2fa-generator.git
   cd steam-2fa-generator
   ```
2. Build
    ```bash
    go build -o steam-2fa-generator.exe main.go
    ```
3. Usage
    ```bash
    ./steam-2fa-generator.exe --secret_key W5G6U5QJNY2XG4RK2W5B7W5TQFQ6V4N3
    ```
   Replace `W5G6U5QJNY2XG4RK2W5B7W5TQFQ6V4N3` with your actual shared secret key.

