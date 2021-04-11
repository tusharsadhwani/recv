# recv

A super simple, super fast way to transfer text and URLs between two devices.

![Demo](https://user-images.githubusercontent.com/43412083/114314685-43dd5b00-9b19-11eb-8c57-bdf3235c32b7.gif)

## Usage

- You can use the website directly: [recv.live](https://recv.live)
- You can install the website as a PWA, and share directly to it to instantly create a room
- You can download and install the recv CLI app from the Releases page:

  - Windows:

    - rename the `recv-windows-XX.exe` file to `recv.exe` and add it to PATH
    - simply type `recv` to create a room
    - type `recv <roomcode>` to join a room, eg. `recv 12345`

  - MacOS/Linux:

    - Extract `recv` from the `recv-PLATFORM-XX.exe` file, and add it to PATH

      For eg. on linux you can do:

      ```bash
      tar -xvzf recv-linux-x86-64.tar.gz
      sudo install ./dist/recv /usr/local/bin
      ```

## Deploy your own

- Server

  ```bash
  cd server
  go build ./cmd/server
  ./server
  ```

  you can specify a port using the `PORT` using the `-p` flag:

  ```bash
  ./server -p 5000
  ```

- CLI

  ```bash
  cd server
  go build -o recv ./cmd/cli
  ./recv 12345
  ```

## Local development

- Server

  > you will need [air](https://github.com/cosmtrek/air) installed

  ```bash
  air -c .air.toml
  ```

  you can specify a port using the `PORT` environment variable:

  ```bash
  PORT=5000 air -c .air.toml
  ```

- CLI

  ```bash
  APP_ENV=dev go run ./cmd/cli
  ```

  you can specify port and room code using environment variables:

  ```bash
  PORT=5000 ROOM=12345 go run ./cmd/cli
  ```
