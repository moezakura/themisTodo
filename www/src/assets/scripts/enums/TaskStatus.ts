export const enum TaskStatus {
    TODO = "TODO",
    DOING = "DOING",
    PULL_REQUEST = "PULL_REQUEST",
    DONE = "DONE",
    HIDE = "HIDE",
    OTHER = "OTHER",
}

export default class TaskStatusConvert {
    static toNumber(taskStatus: TaskStatus): number {
        let numberList = {
            "TODO": 0,
            "DOING": 1,
            "PULL_REQUEST": 2,
            "DONE": 3,
            "HIDE": 4,
            "OTHER": 5,
        }

        return numberList[taskStatus.toString()]
    }
}