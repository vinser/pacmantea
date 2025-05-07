# PacManTea

PacManTea is a terminal-based game inspired by the classic Pac-Man, with customizable themes and styles, including Greek, Hebrew, and modern emoji-based designs.

## Features

- **Customizable Ghost Styles**: Choose from Latin, Greek, Hebrew, and more.
- **Dynamic Maze Rendering**: Supports decorative patterns and custom layouts.
- **Difficulty Levels**: Easy, Medium, and Hard modes with adjustable ghost speeds and timers.
- **Themed Mazes**: Includes Greek-inspired ornamental mazes and other creative designs.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/vinser/pacmantea.git
   cd pacmantea
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the game:
   ```bash
   go run .
   ```

## Configuration

You can customize the game by editing the `config.yml` file:
- Change ppacman ang ghost styles (badges) 
- Adjust difficulty settings like ghost speed and revival timers.
- Add new levels with unique maze layouts.
To create default `config.yml` in config folder run app with `-config` flag.

For examples, see the [example configuration](https://github.com/vinser/pacmantea/blob/master/config-example.yml).

## Credits

PacManTea is built using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework. Bubble Tea is a powerful, flexible, and fun library for building terminal applications in Go. Special thanks to the Charmbracelet team for their amazing work!

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

Enjoy!