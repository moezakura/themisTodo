var taskBoardLists = [
    document.querySelector("#todo>.taskList"),
    document.querySelector("#doing>.taskList"),
    document.querySelector("#pr>.taskList"),
    document.querySelector("#done>.taskList")
];

for (var i = 0; i < taskBoardLists.length; i++) {
    var statusId = i;
    Sortable.create(taskBoardLists[i], {
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
