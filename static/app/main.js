const get = (path, cb) => {
    fetch(path)
        .then(res => {
            return res.json()
        })
        .then(data => {
            cb(data);
        })
}

const populateDeviceList = () => {
    get('/deviceCount', data => {
        const deviceList = document.getElementById('device-list');
        deviceList.innerHTML = '';
        const pre = document.createElement('pre');
        data.List = data.List.filter(device => device !== "")
        pre.append(`${data.List.length} devices on network\n`)
        pre.append('\n')
        data.List.forEach(device => {
            pre.append(`${device}\n`);
        });
        deviceList.append(pre);
    })
}

populateDeviceList();
const poll = () => {
    setTimeout(() => {
        populateDeviceList();
        poll();
    }, 3000)
}
poll();
