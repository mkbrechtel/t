package cli

import (
    "fmt"
    "log"
    
    todo "t5.mkbrechtel.dev/sync/todotxt"
)

// SyncTasks synchronizes tasks between the todo.txt file and the event store
// Follows the t5 sync todo [file] pattern
func SyncTasks(ctx *AppContext) {
    // Sync with the specified todo.txt file
    result, err := todo.SyncWithRepository(ctx.Repository, ctx.Config.TodoFile)
    if err != nil {
        log.Fatalf("Sync failed: %v", err)
    }
    
    // Print results
    fmt.Printf("Sync completed for %s:\n", ctx.Config.TodoFile)
    fmt.Printf("  Added: %d\n", result.Added)
    fmt.Printf("  Updated from todo.txt: %d\n", result.FromTodoTxt)
    fmt.Printf("  Updated to todo.txt: %d\n", result.ToTodoTxt)
    fmt.Printf("  Skipped: %d\n", result.Skipped)
    if result.Conflicts > 0 {
        fmt.Printf("  Conflicts detected: %d\n", result.Conflicts)
    }
}