# Cropper

Simple bulk video editor which has only crop feature implemented.

Video corpping is done by ffmpeg, for which the program will provide geometry
of the area you wish to crop.

## Usage:

Execute the program in directory with video files you wish to crop. A window
will open which allows you to select an area from the video that you wish to
keep.

By pressing **enter** the video will be cropped with current selection into a
new file and next video, if such exists, will be prompted for cropping.
Existing  files won't be overwritten, new files are created with '_cropped' 
appended to the original filename.

At the moment the program won't accept any flags or command line arguments.
Once executed, it will look for video files from current directory.

### Keybinds

| Key             | Action                                          |
|-----------------|-------------------------------------------------|
| Left mouse      | Grab and move selection edges                   |
| Right mouse     | Pan viewport                                    |
| Mouse scroll    | Zoom viewport                                   |
| C               | Reset viewport                                  |
| R               | Reset current selection                         |
| Enter           | Export current crop and move on to next video   |
| Esc             | Close                                           |

