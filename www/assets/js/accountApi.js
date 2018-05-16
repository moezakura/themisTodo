class AccountApi {
    static NewAccountObject(uuid) {
        return {
            "name": "",
            "displayName": "",
            "uuid": uuid,
            "password": ""
        };
    }

    static Change(accountObject){
        return fetch("/account/update/" + accountObject.uuid, {
            method: 'POST',
            body: JSON.stringify(accountObject),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }
}