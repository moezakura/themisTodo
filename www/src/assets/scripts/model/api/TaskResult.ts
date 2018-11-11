import BaseApiResult from "@scripts/model/api/BaseApiResult"
import Task from "@scripts/model/api/task/Task"

export default class TaskResult extends BaseApiResult {
    public task: Task
}