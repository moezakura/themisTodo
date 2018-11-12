export default class TaskAddRequest {
    public name: string
    public deadline: string
    public description: string
    public assign: number
    public projectId: number

    public toJson(): string {
        return JSON.stringify(this)
    }
}