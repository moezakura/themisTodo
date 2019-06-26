import User from "@scripts/model/api/user/User"

export default class UserWithCheck extends User {
    public check: boolean

    constructor(user: User) {
        super()
        this.uuid = user.uuid
        this.displayName = user.displayName
        this.iconPath = user.iconPath
        this.name = user.name
    }
}