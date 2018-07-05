export default class TaskApi {
    static NewTaskObject() {
        return {
            "status": -1,
            "name": "",
            "assign": -1,
            "creatorName": "",
            "id": -1,
            "deadline": "",
            "description": ""
        };
    }

    static GetTaskStatuses() {
        return {
            STATUS_TODO: 0,
            STATUS_DOING: 1,
            STATUS_PULL_REQUEST: 2,
            STATUS_DONE: 3,
            STATUS_HIDE: 4,
        };
    }

    static stringToIntStatus(str) {
        switch (str) {
            case "todo":
                return 0;
            case "doing":
                return 1;
            case "pr":
                return 2;
            case "done":
                return 3;
            case "hide":
                return 4;
        }
    }

    static GetTaskFromCreateDate(createDate) {
        return fetch("/tasks/view/" + createDate, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }

    static GetSearch(taskId, projectId) {
        return fetch(`/tasks/search?taskId=${taskId}&projectId=${projectId}`, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }

    static Create(taskJson) {
        return fetch("/tasks/create", {
            method: 'POST',
            body: JSON.stringify(taskJson),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }

    static Update(createDate, taskApi) {
        return fetch("/tasks/update/" + createDate, {
            method: 'POST',
            body: JSON.stringify(taskApi),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }

    static Delete(createDate) {
        return fetch("/tasks/delete/" + createDate, {
            method: 'POST',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }

    static Get(createDate) {

    }

    static Search(projectId, status, name, description) {
        let query = {
            projectId: projectId,
            status: status,
            name: name === undefined ? "" : name,
            description: description === undefined ? "" : description
        };
        let queryString = "";
        Object.keys(query).forEach(key => {
            let value = query[key];
            queryString += (queryString === "" ? "" : "&") + key + "=" + encodeURIComponent(value);
        });

        return fetch("/tasks/searches?" + queryString, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }
}