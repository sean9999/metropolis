const es = new EventSource('/events');
const sse = document.getElementById('sse');
const btnReset = document.getElementById('reset');

btnReset.addEventListener("click", () => {
    sse.childNodes.forEach(el => {
        el.remove();
    });
});

es.addEventListener("message", ev => {
    console.log(ev);
    let li = document.createElement("li");

    let line = `${atob(ev.data)}\t${ev.timeStamp}\t${ev.lastEventId}`;
    li.innerText = line;

    //li.innerText = (atob(ev.data) + "\t" + `ts:${ev.timeStamp}`
    sse.appendChild(li);
});
es.addEventListener("error", console.error);
