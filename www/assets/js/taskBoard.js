import TaskApi from "./taskApi"
import TaskDetail from "./taskDetail";
import ProjectUtils from "./projectUtils";

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

                    TaskApi.Update(targetCreateId, task).then(function (json) {
                        TaskBoard.loadTask(targetCreateId);
                    });
                },
                animation: 100
            });
        }

        if (this.taskBoardLists[0] === undefined || this.taskBoardLists[0] == null)
            return;

        if (document.location.hash !== "") {
            TaskDetail.loadAndShowFromTaskId(
                document.location.hash.replace("#", ""),
                projectId
            );
        }

        window.addEventListener("hashchange", function () {
            if (document.location.hash !== "") {
                TaskDetail.loadAndShowFromTaskId(
                    document.location.hash.replace("#", ""),
                    projectId
                );
            }
        }, false);
    }

    /***
     *
     * @param createDate {Number}
     */
    static loadTask(createDate) {
        return TaskApi.GetTaskFromCreateDate(createDate).then(function (json) {
            if (json.success) {
                let task = json.task;
                ProjectUtils.taskboadOnTaskUpdate(task.createDate, task);
            } else location.reload();
        });
    }
}

document.taskBoard = new TaskBoard();