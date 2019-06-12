import {TaskStatus} from "@scripts/enums/TaskStatus"

export default class TaskSearchRequest {
    public taskId: number
    public projectId: number
    public status: number
    public assign: number
    public creator: number

    public toQueryString(): string {
        let queryString = "?"
        for (let key in this) {
            // noinspection JSUnfilteredForInLoop
            const value: any = this[key]
            if (typeof value === "object" || typeof value === "function") {
                continue
            }
            const _value: number = value;
            queryString += key + "=" + encodeURIComponent(_value.toString()) + "&"
        }
        queryString = queryString.slice(0, -1)

        return queryString
    }


}