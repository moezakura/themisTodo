import TaskApi from "./taskApi"
import TaskDetail from "./taskDetail";

class TaskBoard {
    constructor() {
        this.taskBoardLists = [
            document.querySelector("#todo>.taskList"),
            document.querySelector("#doing>.taskList"),
            document.querySelector("#pr>.taskList"),
            document.querySelector("#done>.taskList")
        ];

        for (let i = 0; i < this.taskBoardLists.length; i++) {
            if (this.taskBoardLists[i] === undefined || this.taskBoardLists[i] == null)
                continue;

            let statusId = i;
            Sortable.create(this.taskBoardLists[i], {
                group: "shares",
                onEnd: function (evt) {
                    let statusStr = evt.item.parentNode.dataset.status;
                    let targetCreateId = evt.item.dataset.id;
                    let task = TaskApi.NewTaskObject();

                    task.status = TaskApi.stringToIntStatus(statusStr);

                    TaskApi.Update(targetCreateId, task);
                },
                animation: 100
            });
        }

        if (document.location.hash !== "") {
            TaskDetail.loadAndShowFromTaskId(
                document.location.hash.replace("#", ""),
                document.location.pathname.split("/").pop() // get projectId
            );
        }

        window.addEventListener("hashchange", function() {
            if (document.location.hash !== "") {
                TaskDetail.loadAndShowFromTaskId(
                    document.location.hash.replace("#", ""),
                    document.location.pathname.split("/").pop() // get projectId
                );
            }
        }, false);
    }
}

document.taskBoard = new TaskBoard();