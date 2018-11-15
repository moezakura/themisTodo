import Account from "../model/Account"
import AccountUpdateResult from "../model/api/AccountUpdateResult"
import User from "@scripts/model/api/user/User"
import ProfileResult from "@scripts/model/ProfileResult"

export default class AccountApi {
    static GetProfile(): Promise<ProfileResult> {
        return fetch("/api/account/profile", {
            method: 'GET',
            credentials: "same-origin"
        }).then(res => {
            return res.json()
        }).then(json => {
            const res = new ProfileResult()

            res.success = json["success"]
            res.message = json["success"]

            let user = new User()
            const resUser = json["user"]
            if (resUser != null) {
                user.uuid = resUser["uuid"]
                user.name = resUser["name"]
                user.displayName = resUser["displayName"]
            }
            res.user = user

            return res
        })
    }

    static Change(accountObject: Account): Promise<AccountUpdateResult> {
        return fetch("/account/update", {
            method: 'POST',
            body: accountObject.toJson(),
            credentials: "same-origin"
        }).then(response => {
            return response.json()
        }).then(resJson => {
            const resObj = new AccountUpdateResult()
            resObj.success = resJson.success
            resObj.message = resJson.message

            return resObj
        })
    }
}