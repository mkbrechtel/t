package todo

import (
	"log"
	"fmt"
	"time"
	todo "github.com/1set/todotxt"
	uuid "github.com/gofrs/uuid/v5"
)

func ReadTodoFile(todoFile string) {
	if taskList, err := todo.LoadFromPath(todoFile); err != nil {
		log.Fatal(err)
	} else {
		EnsureProperTasks(taskList)
		PrintTaskList(taskList)
		if err = taskList.WriteToPath(todoFile); err != nil {
			log.Fatal(err)
		}
	}
}


func EnsureProperTasks(taskList todo.TaskList)(todo.TaskList) {
	for i, t := range taskList {
		// make sure task has created date
		if(!t.HasCreatedDate()){
			taskList[i].CreatedDate = time.Now()
		}
		// make sure task has a AdditionalTags map
		if (taskList[i].AdditionalTags == nil) {
			taskList[i].AdditionalTags = make(map[string]string)
		}
		// make sure task has uuid
		if _, hasUuid := taskList[i].AdditionalTags["uuid"]; !hasUuid {
			taskList[i].AdditionalTags["uuid"] = uuid.Must(uuid.NewV7()).String()
		}
    }
	return taskList
}

func PrintTaskList(taskList todo.TaskList){
	for _, t := range taskList {
		fmt.Println(t.ID,t)
	}
}
