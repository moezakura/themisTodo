import Account from "../model/Account"
import AccountUpdateResult from "../model/api/AccountUpdateResult"
import User from "@scripts/model/api/user/User"
import ProfileResult from "@scripts/model/ProfileResult"
import AccountUpdateImageResult from "@scripts/model/api/AccountUpdateImageResult"

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
                user.iconPath = resUser["iconPath"]
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

    static uploadImage(uploadData: FormData, progressEvent: (progress: number) => void): Promise<AccountUpdateImageResult> {
        return new Promise((resolve, reject) => {
            let uploadXhr = new XMLHttpRequest()
            let upload = uploadXhr.upload
            if (upload) {
                uploadXhr.addEventListener("load", e => {
                    let json = JSON.parse(uploadXhr.responseText)
                    resolve(json)
                })
            }

            upload.addEventListener("error", e => {
                reject("some error")
            })
            upload.addEventListener("abort", e => {
                reject("abort")
            })

            upload.addEventListener("progress", e => {
                console.log(e.loaded + " / " + e.total)
                const progress = e.loaded / e.total
                progressEvent(progress)
            })

            uploadXhr.open("POST", "/api/account/updateIcon")
            uploadXhr.send(uploadData)
        }).then(json => {
            let res = new AccountUpdateImageResult();

            res.success = json["success"]
            res.message = json["message"]
            res.fileId = json["fileId"]

            return res
        })
    }
}