# Simple Version Control System (SVCS)

This project is a simple version control system (SVCS) implemented in Go. It provides basic functionalities such as configuring a username, adding files to an index, committing changes, viewing commit logs, and checking out to specific commits.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
  - [Config](#config)
  - [Add](#add)
  - [Commit](#commit)
  - [Log](#log)
  - [Checkout](#checkout)
- [License](#license)

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/svcs.git
    cd svcs
    ```

2. Build the project:
    ```sh
    go build -o svcs
    ```

3. Ensure the binary is executable:
    ```sh
    chmod +x svcs
    ```

## Usage

### Config

Set or get the username for the version control system.

- Set username:
    ```sh
    ./svcs config <username>
    ```

- Get username:
    ```sh
    ./svcs config
    ```

### Add

Add a file to the index.

- Add a file:
    ```sh
    ./svcs add <file>
    ```

- View tracked files:
    ```sh
    ./svcs add
    ```

### Commit

Save changes to the repository with a commit message.

- Commit changes:
    ```sh
    ./svcs commit <message>
    ```

### Log

Show the commit logs.

- View commit logs:
    ```sh
    ./svcs log
    ```

### Checkout

Restore a file to a previous commit.

- Checkout to a commit:
    ```sh
    ./svcs checkout <commit_hash>
    ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
