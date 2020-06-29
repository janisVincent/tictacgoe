var grid = document.getElementById("grid");

function down_handler(ev) {
    var cell = ev.target;

    if ("" !== cell.innerHTML) {
        return;
    }

    fetch('/play', {
        method: 'POST',
        body: JSON.stringify({ "cell": parseInt(cell.dataset['id']) })
    })
    .then(response => response.json())
    .then(jsonData => {
        cell.innerHTML = jsonData.Symbol
    })
    .catch(err => {
        console.log(err)
    })
}

grid.onpointerdown = down_handler;