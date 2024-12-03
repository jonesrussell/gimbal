# Gimbal

Gimbal is a simple, fun game built using the Ebiten game library and the Ebitengine input library in Go. The game is a clone of Gyruss.

## Technologies

This project uses the following technologies:

- **Go**: The programming language used to develop the game. Version: `1.23`
- **Ebitengine**: An open-source game engine for the Go programming language. Version: `v2.6.6`
- **Task**: Modern task runner and build tool that aims to be simpler and easier to use than Make
- **HTML**: Contains the index.html file for the game for WASM

## Running the Game

To run the game, navigate to the directory where this repository is cloned and run the following command:

```bash
task run
```

## Controls

The game is controlled using the left and right arrow keys or the 'A' and 'D' keys. Pressing the left arrow key or 'A' moves the player to the left, and pressing the right arrow key or 'D' moves the player to the right.

## Installation Instructions

Follow these steps to set up the project on your local machine:

1. **Install Go**
    - Visit the [Go downloads page](https://golang.org/dl/) and download the latest version for your operating system
    - Install Go by following the prompts

2. **Install Task**
    - Visit [Task installation guide](https://taskfile.dev/installation/)
    - Or install with:
    ```bash
    go install github.com/go-task/task/v3/cmd/task@latest
    ```

3. **Install Required Dependencies (Ubuntu/Debian)**
    ```bash
    sudo apt-get install libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev libxxf86vm-dev
    ```

4. **Clone the Repository**
    - Navigate to the directory where you want to clone the repository
    - Run `git clone https://github.com/jonesrussell/gimbal.git` to clone the repository

5. **Install Dependencies**
    - Navigate to the cloned repository's directory
    - Run `go mod download` to download the necessary Go modules

You should now have the game set up and ready to run on your local machine. To run the game, use the command `task run`.

## Building the Game

The game can be built for Linux, Windows (Win32), and WebAssembly. To build the game, navigate to the directory where the game is installed and run the following commands:

- For Linux:

  ```bash
  task build:linux
  ```

- For Windows (Win32):

  ```bash
  task build:win32
  ```

- For WebAssembly:

  ```bash
  task build:web
  ```

## Development Roadmap

### Research and Planning
- [ ] Play the original Gyruss game to understand its mechanics, controls, and aesthetics
- [ ] Document the features and behaviors you observe for future reference

### Design
- [ ] Sketch out the game's levels, enemies, power-ups, and other elements
- [ ] Plan the user interface, including menus, score displays, and control schemes

### Development Environment Setup
- [x] Choose a game engine or programming language that suits your needs and skills
- [x] Set up your development environment, including any necessary software and tools

### Gameplay Mechanics
- [x] Implement player movement
- [ ] Implement shooting mechanics
- [ ] Implement enemy behavior
- [ ] Implement level progression
- [ ] Implement Entity Component System

### Collision Detection
- [ ] Implement collision detection between the player's ship and enemy ships
- [ ] Implement collision detection between the player's shots and enemy ships

### Scoring
- [ ] Develop a scoring system
- [ ] Define the points for shooting down enemy ships
- [ ] Implement bonus points for meeting certain conditions

### Player Lives
- [ ] Implement a system to track the player's lives
- [ ] Define the conditions for losing a life
- [ ] Implement game over condition when player has no more lives

### User Interface (UI)
- [ ] Design the UI layout
- [ ] Implement the player's score display
- [ ] Implement the display for the number of lives remaining

### Sound Effects and Music
- [ ] Choose or create sound effects for shooting, explosions, and other game events
- [ ] Choose or create background music

### Intro Screen
- [ ] Design the intro screen layout
- [ ] Implement options to start the game, view high scores, and exit the game

### High Score Screen
- [ ] Design the high score screen layout
- [ ] Implement a system to track and display the highest scores

### End Game
- [ ] Define what happens when the game ends
- [ ] Implement a game over screen
- [ ] Implement a system to return to the intro screen after the game ends

### Testing
- [ ] Playtest the game regularly, both yourself and with others
- [ ] Use feedback to refine and improve the game
- [ ] Test the game on different platforms and systems to ensure compatibility

### Release
- [ ] Prepare the game for release
- [ ] Create promotional materials, set up a website, and submit the game to different platforms
- [ ] Release the game and monitor for feedback

## License

Gimbal is licensed under the MIT License. See the `LICENSE` file for more information.
