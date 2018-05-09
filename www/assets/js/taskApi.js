class TaskApi {
    static GetTaskFromCreateDate(createDate) {
        let params = {
            "time": createDate
        };
        let esc = encodeURIComponent;
        let query = Object.keys(params).map(k => esc(k) + '=' + esc(params[k])).join('&');

        return fetch("/task?" + query, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }
}