import Account from "../model/Account"
import AccountUpdateResult from "../model/AccountUpdateResult"

export default class AccountApi {
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
        });
    }
}