# t5 sync

The t5 sync works by having sync stores that you can configure with multiple formats. Usually the goal is it to make actionable tasks that are present in those stores available to the users and not so much to sync the whole store state between multiple systems, so we concentrate on exchaning task information and uploading time tracking information to those stores where requested and possible.

The default store is a t5 append log store, with this store you can track the state of the different store backends and are able to detect incoming changes and therefore detect possible conflicts that might occur.

## Sync Groups

The sync system uses the concept of sync groups to control which tasks are synchronized between systems. Only tasks that belong to a specific sync group will be synchronized with providers in that group. This prevents unwanted synchronization and keeps task data compartmentalized.

Each sync group:
- Has a unique identifier
- Can contain multiple sync providers
- Only affects tasks explicitly assigned to that group
- Can be configured with different synchronization rules
- Allows for targeted synchronization of specific task sets

Tasks can be assigned to sync groups using the `sync:group_name` tag.

## Directions

The sync has two directions when it comes to the data flow, the sync store can be a source or a destination or both. Not all storage backends support all directions. 

## Filters and Modifiers

You can use filters and modifiers for both directions of the task and information flow: inbound (when data flows into t5) and outbound (when data flows to the provider). This allows you to select which tasks to sync and apply certain properties when syncing in either direction. For example, you might have a certain ticket system for a specific project and can sync tasks from this system into your todo list (inbound), with all tasks being modified so they get added to the project accordingly. Similarly, when sending tasks back to the provider (outbound), you can filter and modify the data to meet the requirements of the external system.

## Example sync configuration file (config.yaml)

```yaml
# Example configuration file
sync_groups:
  - name: personal
    todo_txt_files:
      file: ~/Documents/todo.txt
      enabled: true
      direction: both
      ensure_properties: true

  - name: work
    todo_txt_files:
      file: ~/work/todo.txt
      enabled: true
      direction: both
      ensure_properties: true
    github_sync:
      repo: username/repository
      enabled: true
      direction: import
      token: ${GITHUB_TOKEN}
      labels: ["todo", "task"]
      import_filter:
        project: work
        maxissues: 50
      import_modifier:
        addproject: github
        addcontext: online
        add_sync_group: work
    gitlab_sync:
      project: my-project
      enabled: false  # disabled provider
      direction: export
      url: https://gitlab.company.com
      token: ${GITLAB_TOKEN}
      export_filter:
        project: gitlab
        not_completed: true
        has_tag: sync:work  # only sync tasks with the work sync group tag
```
