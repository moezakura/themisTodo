import ProjectUtils from './projectUtils.js';

class Home {
    constructor() {
        this.todoListElem = document.querySelector("#todoList");
        this.doingListElem = document.querySelector("#doingList");
        if (this.todoListElem === undefined || this.todoListElem == null) return;

        let that = this;
        todoList.forEach(function (value) {
            let item = ProjectUtils.createTaskItem(value.createDate, value.name, value.taskId, value.assignName, value.assign, value.deadline, value.limitDate);
            that.todoListElem.appendChild(item);
        });

        doingList.forEach(function (value) {
            let item = ProjectUtils.createTaskItem(value.createDate, value.name, value.taskId, value.assignName, value.assign, value.deadline, value.limitDate);
            that.doingListElem.appendChild(item);
        });
    }
}

new Home();