# t5 Data Model

The data model for t5 is designed to support the core features of the application, including managing tasks and tracking time. The model is based on a set of entities and relationships that capture the essential information needed to implement these features.

We separate the data model into Entities and Events to get an event-sourced application. The state of the Entities is only derived from Events. The Entities are also the aggregate root of the application.

## Entities

### Project
- **ID**: Internal project ID.
- **Name**: Name of the project.
- **Description**: Description of the project.
- **Tasks**: Tasks associated with the project.
- **Budget**: Budgeted time for the project.

### Task
- **ID**: Internal task ID.
- **Original**: Original raw task text.
- **Todo**: Todo part of task text.
- **Priority**: Priority of the task.
- **Projects**: Projects associated with the task.
- **Contexts**: Contexts associated with the task.
- **AdditionalTags**: Addon tags will be available here.
- **CreatedDate**: Date when the task was created.
- **DueDate**: Date when the task is due.
- **CompletedDate**: Date when the task was completed.
- **Completed**: Indicates if the task is completed.
- **UsedTime**: Time used for the task.


## Events

### TaskUpdate

Update a task with new information.

- **ID**: Internal task ID.
- **Original**: Original raw task text.
- **Todo**: Todo part of task text.
- **Priority**: Priority of the task.
- **Projects**: Projects associated with the task.
- **Contexts**: Contexts associated with the task.
- **AdditionalTags**: Addon tags will be available here.
- **CreatedDate**: Date when the task was created.
- **DueDate**: Date when the task is due.
- **CompletedDate**: Date when the task was completed.
- **Completed**: Indicates if the task is completed.
- **UsedTime**: Time used for the task.
- **Source**: Source of the task.

### TaskStartTime

Start a task.

- **ID**: Internal task start time ID.
- **TaskID**: ID of the task associated with the start time.
- **StartTime**: Start time of the task.

### TaskEndTime

End a task.

- **ID**: Internal task end time ID.
- **TaskID**: ID of the task associated with the end time.
- **EndTime**: End time of the task.

### SetTimeBudgetForProject

Set a time budget for a project.

- **ID**: Internal time budget ID.
- **ProjectID**: ID of the project associated with the time budget.
- **Budget**: Budgeted time for the project.
- **UsedTime**: Time used for the project.
- **Description**: Description of the time budget.
