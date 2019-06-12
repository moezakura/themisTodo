export default class AccountSearchRequest {
    public name: string
    public displayName: string
    public project: number
    public isInProject: boolean
    public max: number

    public toQueryString(): string {
        let queryString = "?"
        for (let key in this) {
            // noinspection JSUnfilteredForInLoop
            const value: any = this[key]
            if (typeof value === "object" || typeof value === "function") {
                continue
            }
            const _value: number = value
            queryString += key + "=" + encodeURIComponent(_value.toString()) + "&"
        }
        queryString = queryString.slice(0, -1)

        return queryString
    }
}