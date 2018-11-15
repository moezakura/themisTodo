import BaseApiResult from "@scripts/model/api/BaseApiResult"

export default class AccountCreateResult extends BaseApiResult {
    public name: string
    public password: string
}