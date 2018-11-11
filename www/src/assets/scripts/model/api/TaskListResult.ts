import BaseApiResult from "@scripts/model/api/BaseApiResult"
import Task from "@scripts/model/api/task/Task"

export default class TaskListResult extends BaseApiResult {
    public task: Array<Task>
}