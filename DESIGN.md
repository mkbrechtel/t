# t software design

goal: simple application for some tasks to aid me in my time and management

## pitch

This tool is supposed to help me with my time management. It is based on the ideas of Getting Things Done and timeboxing. You can select one todo as your current activity and it will track the time spend on it.

## planned features

### todo list manager

- based on todo.txt format
- basic cli for adding, searching, changing and sorting tasks

### activitiy tracker

- based on timeclock format
- select a todo task as the current activity
- use i3 socket protocol to automatically guess what i am currently doing

### i3/sway support

- i3blocks status output
- i3 mode template for easy usage

### time management gadgets

- timeboxing/pomodoro feature

### todo.txt split, sync and merge tool

- split out the tasks of a big todo txt file based on filter criteria, for example all tasks of a project, in a seperate file
- sync todo.txt files between multiple locations
  - file system
  - git repo
  - webdav
  - ssh/sftp remote location
  - t websocket live sync
- intelligently merge multiple todo.txt files

### web app

- REST API
- t websocket live sync
- web GUI
