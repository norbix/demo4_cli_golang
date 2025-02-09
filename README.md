# demo4_cli_golang

Demo application containing a CLI solution in pure Golang and the application is a file filter CLI program for OS of your choice.

## Specification

### Overview

- Using systems programming language of your choice (Go / C / C++) write a file filter
  CLI program for OS of your choice. This filter should accept hot folder path and
  backup folder path, backing-up any file that is created or modified in the chosen
  folder.

### Requirements

- create a copy of any file created or modified in the hot folder
- backup files should have the same name of the original file with .bak extension
- if the file name is prefixed with 'delete_' it should be immediately deleted from
  the hot folder and backup folder
- keep a log file of all action taken by your program (file created, altered, backedup
  or deleted)
- log file can be viewed/filtered by you CLI app.
- log file filters should accept filter by [date, filename regex]
- the application must keep/save it's state between reboots

### Bonus

- if the file name is prefixed with 'delete_ISODATETIME_' it should be deleted at the
  specified time
- Use non-blocking IO

