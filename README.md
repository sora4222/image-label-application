# Image Training Set Generator
This is intended to make a GUI that is capable of allowing the user to
sort a directory of images into a set of labelled images by putting them into a directory.
This will allow quick labelling of a dataset for the use of an ML pipeline and to learn
how to write a GUI application in Go.

## Usage
This program can be run with the following command:
```bash
image-label-applicaation <directory> <label-name1>=<key1> <label-name2>=<key2> ...
```
Where `<directory>` is the directory that contains the images to be labelled and the key arguments are the
keyboard character to press to label images with the corresponding label name. In case a person mislabels an image,
the backspace key can be used to undo the last label (not yet implemented).