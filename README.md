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
|                 |                                                 |
|                 | **Pan & Zoom Actions**                          |
| Right mouse     | Pan viewport                                    |
| Mouse scroll    | Zoom viewport                                   |
| C               | Reset viewport                                  |
|                 |                                                 |
|                 | **Crop Actions**                                |
| Left mouse      | Grab and move selection edges                   |
| R               | Reset current selection                         |
|                 |                                                 |
|                 | **Trim Actions**                                |
| A               | Set trimming beginning                          |
| B               | Set trimming end                                |
|                 |                                                 |
|                 | **Seek Actions** (*)                            |
| Up              | Seek one second forward                         |
| Down            | Seek one second backward                        |
| Right           | Seek one frame forward                          |
| Left            | Seek one frame backward                         |
| Shift + B       | Seek to video beginning                         |
| Shift + E       | Seek to video end                               |
|                 |                                                 |
|                 | **Export Actions**                              |
| M               | Toggle mute                                     |
| Enter           | Export current video and move on to next video  |
| S               | Export screenshot and move on to next video     |
| Esc             | Close                                           |

(*) Seeking actions can be multiplied similar to vim actions. For example seeking 
15 frames forward can be accomplished by pressing "1+5+Right", or two seconds
backwards with "2+Down"

## Compilation

Run the following to create binary `cropper`:

```
$ git clone https://github.com/leevi-laitala/cropper.git
$ cd cropper
$ make build
```

