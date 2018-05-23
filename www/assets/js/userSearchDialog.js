class UserSearchDialog {
    constructor(target, config) {
        this.selectUserIndex = -1;
        this.selectMode = false;
        this.oldSearchInput = "";
        this.searchResult = [];
        this.target = target;
        this.userSelectUserList = UserSearchDialog.createSearchBox();
        this.sendEvent = config.forceSubmit;
        this.singleEnter = config.singleEnter === true;
        this.isIn = config.isIn === true;
        let that = this;

        target.addEventListener("keyup", function (e) {
            that.userSelectInputKeyUp(e, that);
        }, true);
        target.addEventListener("keydown", function (e) {
            that.userSelectInputKeyDown(e, that);
        }, true);
    }

    static createSearchBox() {
        let searchBox = document.createElement("ul");
        searchBox.id = "usersSearchedList-" + (new Date()).getTime();
        searchBox.classList.add("usersList");
        searchBox.classList.add("userSearchDialog");
        document.body.appendChild(searchBox);

        return searchBox;
    }

    get() {
        return getUuid(this);
    }

    userSelectInputKeyUp(e, that) {
        let inputText = that.target.value;
        if (inputText.length < 1) {
            that.userSelectUserList.style.display = "none";
            that.selectUserIndex = -1;
            return;
        }

        if (e.keyCode === 38 || e.keyCode === 40 || e.keyCode === 13) return;

        if (inputText !== that.oldSearchInput) {
            that.selectUserIndex = -1;
        }

        let queryObject = {
            "name": inputText,
            "displayName": inputText,
            "project": projectId,
            "isInProject": that.isIn,
            "max": 20
        };
        console.log(queryObject);
        let queryString = "?";
        for (let key in queryObject)
            queryString += key + "=" + encodeURIComponent(queryObject[key]) + "&";
        queryString = queryString.slice(0, -1);

        fetch("/account/search" + queryString, {
            method: 'GET'
        }).then(function (response) {
            return response.json();
        }).then(function (json) {
            that.searchResult = json;
            let targetPos = that.target.getBoundingClientRect();
            that.userSelectUserList.style.display = "block";
            that.userSelectUserList.style.left = (targetPos.left + window.pageXOffset) + "px";
            that.userSelectUserList.style.top = (targetPos.top + targetPos.height + window.pageYOffset) + "px";
            that.userSelectUserList.style.width = targetPos.width + "px";
            that.userSelectUserList.scrollTop = 0;
            that.selectMode = false;

            that.userSelectUserList.innerHTML = "";
            json.forEach(function (value) {
                let elem = ProjectUtils.createUserListLine(value.uuid, value.name, value.displayName);
                that.userSelectUserList.appendChild(elem);
            });
        });
    }

    userSelectInputKeyDown(e, that) {
        if (e.keyCode === 38) {
            //up key
            that.selectUserIndex--;
            if (that.selectUserIndex < 0) that.selectUserIndex = that.searchResult.length - 1;
        } else if (e.keyCode === 40) {
            //down key
            that.selectUserIndex++;
            if (that.selectUserIndex > that.searchResult.length - 1) that.selectUserIndex = 0;
        } else if (e.keyCode === 13) {
            this.projectMemberAddSubmit(e, that);
            return;
        } else return;

        if (that.selectUserIndex < 0 || that.selectUserIndex < 0 || that.selectUserIndex > that.searchResult.length - 1) return;
        that.selectMode = true;

        let all = that.userSelectUserList.querySelectorAll("li");
        all.forEach(function (value) {
            value.classList.remove("select");
        });

        let targetId = "#searchResult" + that.searchResult[that.selectUserIndex].uuid;
        let target = that.userSelectUserList.querySelector(targetId);
        let targetPos = target.getBoundingClientRect();
        let userSelectUserListPos = that.userSelectUserList.getBoundingClientRect();
        that.userSelectUserList.scrollTop = targetPos.y + that.userSelectUserList.scrollTop -
            userSelectUserListPos.height + targetPos.height - userSelectUserListPos.y;
        target.classList.add("select");
    }

    projectMemberAddSubmit(e, that) {
        e.preventDefault();
        console.log(this.singleEnter);
        if(!that.selectMode || this.singleEnter){
            let sendUuid = that.getUuid(that);
            that.sendEvent(sendUuid);
            that.userSelectUserList.style.display = "none";
        }
        else if (that.selectMode) {
            if (!(that.selectUserIndex < 0 || that.selectUserIndex < 0 || that.selectUserIndex > that.searchResult.length - 1))
                that.target.value = that.searchResult[that.selectUserIndex].name;
            that.userSelectUserList.style.display = "none";
            that.selectMode = false;
        }
    }

    getUuid(that) {
        let sendUuid = -1;
        if (that.selectUserIndex < 0 || that.selectUserIndex < 0 || that.selectUserIndex > that.searchResult.length - 1) {

        } else sendUuid = that.searchResult[that.selectUserIndex].uuid;

        return sendUuid;
    }

}