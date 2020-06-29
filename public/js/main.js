var grid = document.getElementById("grid");

function down_handler(ev) {
    var cell = ev.target;

    if ("" !== cell.innerHTML) {
        return;
    }

    fetch('/play', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(jsonData => {
        cell.innerHTML = jsonData.token
        cell.setAttribute('data-token', jsonData.token)
    })
    .catch(err => {
        console.log(err)
    })
}

grid.onpointerdown = down_handler;