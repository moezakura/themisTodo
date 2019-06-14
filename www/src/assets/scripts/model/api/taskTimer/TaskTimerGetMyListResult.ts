import BaseApiResult from "@scripts/model/api/BaseApiResult"
import TaskTimer from "@scripts/model/api/taskTimer/TaskTimer"

export default class TaskTimerGetResult extends BaseApiResult {
    public list: Array<TaskTimer>
}