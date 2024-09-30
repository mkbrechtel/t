# t0 software design

goal: simple application for some tasks to aid me in my time and management

## language selection
- Go
    - todotxt lib: https://github.com/1set/todotxt
- Python
    - todotxt lib: https://vonshednob.cc/pytodotxt/

## pitch
This tool is supposed to help me with my time management. It is based on the ideas of Getting Things Done and timeboxing. You can select one activity you are currently doing and 

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
    - t0 websocket live sync
- intelligently merge multiple todo.txt files

### web app
- REST API
- t0 websocket live sync
- web GUI
