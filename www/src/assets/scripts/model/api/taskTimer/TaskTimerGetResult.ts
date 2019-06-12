import BaseApiResult from "@scripts/model/api/BaseApiResult"

export default class TaskTimerGetResult extends BaseApiResult {
    public start: boolean
    public LastStartTime: number
    public LastEndTime: number
    public TotalTime: number
    public TodayTime: number
}