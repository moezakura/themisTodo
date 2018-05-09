var taskLists = [
    document.querySelector("#todo>.taskList"),
    document.querySelector("#doing>.taskList"),
    document.querySelector("#pr>.taskList"),
    document.querySelector("#done>.taskList")
];

for (var i = 0; i < taskLists.length; i++) {
    Sortable.create(taskLists[i], {
        group: {
            name: "shares",
        },
        animation: 100
    });
}
