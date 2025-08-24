// @author xiangqian
// @date 2025/07/27 15:17

function getElement(instance, name) {
    let id = `${instance.addr},${name}`;
    return document.getElementById(id);
}

document.addEventListener('DOMContentLoaded', function () {
    let line = null;
    let indexes = null;
    let eventSource = new EventSource('/event');
    eventSource.onmessage = (e) => {
        let data = JSON.parse(e.data);

        let apps = data.apps;
        // console.log('apps', apps);
        for (let app of apps) {
            for (let instance of app.instances) {
                let statusElement = getElement(instance, 'status');
                statusElement.textContent = instance.status;
                if (instance.status === 'UP') {
                    statusElement.className = 'status status-ok';
                } else if (instance.status === 'DOWN') {
                    statusElement.className = 'status status-error';
                } else {
                    statusElement.className = 'status status-unknown';
                }

                let timeElement = getElement(instance, 'time');
                timeElement.textContent = instance.time;

                let durationElement = getElement(instance, 'duration');
                durationElement.textContent = instance.duration;
            }
        }

        let sample = data.sample;
        // console.log('sample', sample);
        if(sample == null){
            return;
        }


        if (line == null) {
            indexes = new Map();
            let index = 0;
            let series = new Array();
            for (let name in sample.value) {
                let ser = null;
                let arr = name.split(',');
                let addr = arr[0];
                let metric = arr[1];
                if (metric === 'cpu_usage') {
                    ser = {
                        label: 'CPU',
                        scale: YAxis.Left,
                        format: (u, v) => v === null ? '--' : v.toFixed(2) + '%',
                        formats: (u, vals) => vals.map(v => v.toFixed(0) + '%'),
                    };
                } else if (metric === 'mem_used_percent') {
                    ser = {
                        label: 'MEM',
                        scale: YAxis.Left,
                        format: (u, v) => v === null ? '--' : v.toFixed(2) + '%',
                        formats: (u, vals) => vals.map(v => v.toFixed(0) + '%'),
                    };
                } else if (metric === 'mem_used_bytes') {
                    ser = {
                        label: 'MEM',
                        scale: YAxis.Right,
                        format: (u, v) => v === null ? '--' : formatBytes(v, 2),
                        formats: (u, vals) => vals.map(v => formatBytes(v, 0)),
                    };
                }

                if (ser != null) {
                    let label = null;
                    for (let app of apps) {
                        for (let instance of app.instances) {
                            if (instance.addr === addr) {
                                let arr = addr.split(':');
                                if (arr.length === 2) {
                                    label = `${app.name}-${arr[1]} ${ser.label}`;
                                } else {
                                    label = `${app.name} ${ser.label}`;
                                }
                                break;
                            }
                        }
                        if (label != null) {
                            break;
                        }
                    }
                    ser.label = label;
                    series.push(ser);
                    indexes.set(name, index++);
                }
            }
            // console.log('series', series);
            line = new Line(document.getElementById('chart'), 'CPU/MEM', 800, 400, series);
        }

        let timestamp = sample.timestamp;

        let value = sample.value;
        let values = Array.from({length: indexes.size}, () => 0);
        for (let name in value) {
            let index = indexes.get(name);
            if (index !== undefined) {
                values[index] = value[name];
            }
        }

        line.push(timestamp, ...values);
    };
});

