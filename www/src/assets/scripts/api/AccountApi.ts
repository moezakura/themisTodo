import Account from "../model/Account"
import AccountUpdateResult from "../model/api/AccountUpdateResult"
import User from "@scripts/model/api/user/User"
import ProfileResult from "@scripts/model/api/ProfileResult"
import AccountUpdateImageResult from "@scripts/model/api/AccountUpdateImageResult"
import AccountCreateRequest from "@scripts/model/api/AccountCreateRequest"
import AccountCreateResult from "@scripts/model/api/AccountCreateResult"
import AccountListResult from "@scripts/model/api/AccountListResult"
import Task from "@scripts/model/api/task/Task"
import AccountSearchRequest from "@scripts/model/api/AccountSearchRequest"
import {AccountSearchResult} from "@scripts/model/api/AccountSearchResult"
import BaseApi from "@scripts/api/BaseApi"

export default class AccountApi extends BaseApi {
    static getProfile(): Promise<ProfileResult> {
        return fetch("/api/account/profile", {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
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

    static getList(): Promise<AccountListResult> {
        return fetch("/api/account/list", {
            method: 'GET',
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            const res = new AccountListResult()

            res.success = json["success"]
            res.message = json["success"]

            res.users = []
            for (let i of <Array<any>>json["users"]) {
                let user = new User()

                user.uuid = i["uuid"]
                user.name = i["name"]
                user.displayName = i["displayName"]
                user.iconPath = i["iconPath"]

                res.users.push(user)
            }

            return res
        })
    }

    static change(accountObject: Account): Promise<AccountUpdateResult> {
        return fetch("/api/account/update", {
            method: 'POST',
            body: accountObject.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(response => {
            return response.json()
        }).then(resJson => {
            const res = new AccountUpdateResult()
            res.success = resJson["success"]
            res.message = resJson["message"]

            return res
        })
    }

    static create(createRequest: AccountCreateRequest): Promise<AccountCreateResult> {
        return fetch("/api/account/add", {
            method: 'POST',
            body: createRequest.toJson(),
            credentials: "same-origin",
            headers: this.getHeader(),
        }).then(response => {
            return response.json()
        }).then(resJson => {
            const res = new AccountCreateResult()

            res.success = resJson["success"]
            res.message = resJson["message"]
            res.name = resJson["name"]
            res.password = resJson["password"]

            return res
        })
    }

    static search(searchRequest: AccountSearchRequest): Promise<AccountSearchResult> {
        const queryString = searchRequest.toQueryString()

        return fetch("/api/account/search" + queryString, {
            method: 'GET',
            headers: this.getHeader(),
        }).then(res => {
            return res.json()
        }).then(json => {
            let res:AccountSearchResult = []

            for(const _user of json){
                let user = new User()
                user.name = _user["name"]
                user.uuid = _user["uuid"]
                user.displayName = _user["displayName"]
                user.iconPath = _user["iconPath"]
                res.push(user)
            }

            return res
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
                const progress = e.loaded / e.total
                progressEvent(progress)
            })

            uploadXhr.open("POST", "/api/account/updateIcon")
            uploadXhr.send(uploadData)
        }).then(json => {
            let res = new AccountUpdateImageResult()

            res.success = json["success"]
            res.message = json["message"]
            res.fileId = json["fileId"]

            return res
        })
    }
}