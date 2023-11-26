# Gimbal

Gimbal is a simple, fun game built using the Ebiten game library and the Ebitengine input library in Go. The game features a player that moves around the circumference of a circle.

## Installation

To install the game, you need to have Go installed on your machine. If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/).

Once you have Go installed, you can install the game by running the following command:

```bash
go get github.com/jonesrussell/gimbal
```

## Running the Game

To run the game, navigate to the directory where the game is installed and run the following command:

```bash
make run
```

## Controls

The game is controlled using the left and right arrow keys or the 'A' and 'D' keys. Pressing the left arrow key or 'A' moves the player to the left, and pressing the right arrow key or 'D' moves the player to the right.

## Building the Game

The game can be built for Linux, Windows (Win32), and WebAssembly. To build the game, navigate to the directory where the game is installed and run the following commands:

- For Linux:

  ```bash
  make build/linux
  ```

- For Windows (Win32):

  ```bash
  make build/win32
  ```

- For WebAssembly:

  ```bash
  make build/web
  ```

## License

Gimbal is licensed under the MIT License. See the `LICENSE` file for more information.
