# t5 features

t5 is supposed to provide a flexible solution for task and todo management with the ability to track time. It also has integrations to various ticket systems. 

## todo list
- [ ] list open todos with `t5 todo` or `t5 list todo`
- [ ] update or add todo with `t5 update todo`
  - [ ] reads from stdin if no further argument is given
  - [ ] adds tasks directly from arguments like `t5 add todo "Buy catfood @supermarket"`
- [ ] task filter flags
  - [ ] filter for todo.txt properties with the todo.txt syntax
  - [ ] regex based filter (starting with /)
  - [ ] combine multiple filters with boolean operators like you can do with GNU find
  - [ ] json query based fitlers
- [ ] task modifiers, for example to add project or context info to tasks
  - [ ] modifiers to add or remove todo.txt properties
  - [ ] regex based modifiers, like with sed
  - [ ] json query based modifiers

## activity and time tracking
- [ ] ability to select a project that is currently being worked on
- [ ] select a task currently being worked upon, starting time tracking
- [ ] action to asign a time budget to a project or task
- [ ] pause the current activity
- [ ] finish current activity as done, stopping time tracking and marking it as done

## sync
- [ ] CalDAV sync
- [ ] todo.txt sync
- [ ] GitHub sync
- [ ] GitLab sync
- [ ] OpenProject sync

## i3 integration
- [ ] i3 bar current activity display with a display of the current 
- [ ] i3 bar current project 
- [ ] i3 bar next events and remaining time on task
