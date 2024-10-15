package todo

import (
	"log"
	"fmt"
	"time"
	todo "github.com/1set/todotxt"
	"t/utils"
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
		id := utils.NewUUID()
		// take uuid if it is there and remove it
		if _, hasUuid := taskList[i].AdditionalTags["uuid"]; hasUuid {
			id = utils.DecodeUUID(taskList[i].AdditionalTags["uuid"])
			delete(taskList[i].AdditionalTags, "uuid")
		}
		// set id
		if _, hasId := taskList[i].AdditionalTags["id"]; !hasId {
			taskList[i].AdditionalTags["id"] = utils.EncodeUUID(id)
		}
    }
	return taskList
}

func PrintTaskList(taskList todo.TaskList){
	for _, t := range taskList {
		fmt.Println(t.ID,t)
	}
}
