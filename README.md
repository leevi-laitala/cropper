# Cropper

Simple bulk video editor.

Features:
- Cropping
- Trimming
- Muting

## Usage:

```
$ cropper <path (optional)>
```

Once executed, program finds video files and launches a window for simple video 
editor that allows for cropping, trimming and muting.

Cropper attempts to find video files from path provided with cmd arg. If not 
provided, defaults to cwd. The program will not recurse into any directories.

All edits are commited once **enter** is pressed. Original video file is not 
overwritten, but a "_cropped" suffix is appended to the filename.

### Keybinds

| Key             | Action                                          |
|-----------------|-------------------------------------------------|
| Left mouse      | Grab and move selection edges                   |
| Right mouse     | Pan viewport                                    |
| Mouse scroll    | Zoom viewport                                   |
| A               | Set trimming beginning                          |
| B               | Set trimming end                                |
| Shift+B         | Seek to video beginning                         |
| Shift+E         | Seek to video end                               |
| Right           | Seek one frame forward                          |
| Left            | Seek one frame backward                         |
| Up              | Seek one second forward                         |
| Down            | Seek one second backward                        |
| C               | Reset viewport                                  |
| R               | Reset current selection                         |
| Enter           | Export current video and move on to next video  |
| Esc             | Close                                           |

Seeking actions can be multiplied similar to vim actions. For example seeking 
15 frames forward can be accomplished by pressing "1+5+Right", or two seconds
backwards with "2+Down"

## Compilation

Run the following:

```
$ git clone https://github.com/leevi-laitala/cropper.git
$ cd cropper
$ go build
$ ./cropper
```

