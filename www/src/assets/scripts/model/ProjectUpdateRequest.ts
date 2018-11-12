export default class ProjectUpdateRequest {
    public name: string
    public description: string

    public toJson(): string {
        return JSON.stringify(this)
    }
}