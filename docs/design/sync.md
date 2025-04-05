# t5 sync

The t5 sync works by having sync stores that you can configure with multiple formats. Usually the goal is it to make actionable tasks that are present in those stores available to the users and not so much to sync the whole store state between multiple systems, so we concentrate on exchaning task information and uploading time tracking information to those stores where requested and possible.

The default store is a t5 append log store, with this store you can track the state of the different store backends and are able to detect incoming changes and therefore detect possible conflicts that might occur. 

The sync has two directions when it comes to the data flow, the sync store can be a source or a destination or both. Not all storage backends support all directions. 

You can use filters and modifiers for each direction of the task and information flow, so that you can select which tasks to sync and that certain properties are also applied for them when syncing. For example you might have a certain ticket system for a specific project and can sync tasks out of this system into your todo list, with all tasks being modified so they get added the project accordingly, so you filter them.
