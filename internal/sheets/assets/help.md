# JSON viewer in your terminal.

USAGE
  sheets <file.json>
      Launch sheets interactively with a JSON file.

OPTIONS
  -h, --help
      Show this help message.

  -v, --version
      Show the current version.

  -w N, --width N
      Set the column cell width (default: 40).

EXAMPLES
  jsonsheets data.json
  jsonsheets config.json

NAVIGATION AND MAINTENANCE
jsonsheets follow vi style for management cursor and editing, e.g.
- use `h j k l` or arrow keys to move cursor
- use i to edit current cell, or I key to edit current cell from start
- use c to clear and edit current cell
- use o/O to insert row below/above
- use dd key to delete current row
- use u/U key to undo/redo
- use / to search forward and ? to search backend
  - use n/N key to find next/previous match
- and use :w to save file
- use :q to quit or :q! to quit without saving
- use :wq or :x to save and quit
