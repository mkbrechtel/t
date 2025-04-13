# t5 sync

The t5 sync works by having sync stores that you can configure with multiple formats. Usually the goal is it to make actionable tasks that are present in those stores available to the users and not so much to sync the whole store state between multiple systems, so we concentrate on exchaning task information and uploading time tracking information to those stores where requested and possible.

The default store is a t5 append log store, with this store you can track the state of the different store backends and are able to detect incoming changes and therefore detect possible conflicts that might occur. 

The sync has two directions when it comes to the data flow, the sync store can be a source or a destination or both. Not all storage backends support all directions. 

You can use filters and modifiers for both directions of the task and information flow: inbound (when data flows into t5) and outbound (when data flows to the provider). This allows you to select which tasks to sync and apply certain properties when syncing in either direction. For example, you might have a certain ticket system for a specific project and can sync tasks from this system into your todo list (inbound), with all tasks being modified so they get added to the project accordingly. Similarly, when sending tasks back to the provider (outbound), you can filter and modify the data to meet the requirements of the external system.
