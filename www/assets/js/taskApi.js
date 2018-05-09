class TaskApi {
    static NewTaskObject() {
        return {
            "status": -1,
            "name": "",
            "creator": -1,
            "creatorName": "",
            "id": -1,
            "deadline": "",
            "description": ""
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
        }
    }

    static GetTaskFromCreateDate(createDate) {
        return fetch("/task/" + createDate, {
            method: 'GET',
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }

    static Update(createDate, taskApi) {
        return fetch("/task/" + createDate + "/update", {
            method: 'POST',
            body: JSON.stringify(taskApi),
            credentials: "same-origin"
        }).then(function (response) {
            return response.json();
        });
    }
}