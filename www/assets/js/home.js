var todoListElem = document.querySelector("#todoList"),
    doingListElem = document.querySelector("#doingList");

todoList.forEach(function (value) {
    let item = ProjectUtils.createTaskItem(value.createDate, value.name, value.taskId, value.assignName, value.assign, value.deadline, value.limitDate);
    todoListElem.appendChild(item);
});

doingList.forEach(function (value) {
    let item = ProjectUtils.createTaskItem(value.createDate, value.name, value.taskId, value.assignName, value.assign, value.deadline, value.limitDate);
    doingListElem.appendChild(item);
});
