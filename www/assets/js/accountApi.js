export default class AccountApi {
    static NewAccountObject() {
        return {
            "name": "",
            "displayName": "",
            "password": "",
            "currentPassword": "",
        };
    }

    static Change(accountObject) {
        return fetch("/account/update", {
            method: 'POST',
            body: JSON.stringify(accountObject),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }
}