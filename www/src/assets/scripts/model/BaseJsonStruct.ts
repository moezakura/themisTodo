export default class BaseJsonStruct {
    public toJson(): string {
        return JSON.stringify(this)
    }
}