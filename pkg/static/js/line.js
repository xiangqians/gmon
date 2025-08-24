// @author xiangqian
// @date 2025/07/26 22:52

/**
 * 扩展日期格式化函数
 * yyyy/MM/dd HH:mm:ss
 * yyyy/MM/dd HH:mm:ss.SSS
 *
 * @param pattern
 * @returns {string}
 */
Date.prototype.format = function (pattern) {
    let object = {
        "yyyy": this.getFullYear().toString(),
        "yy": this.getFullYear().toString().substring(2),
        "MM": (this.getMonth() + 1).toString().padStart(2, '0'),
        "dd": this.getDate().toString().padStart(2, '0'),
        "HH": this.getHours().toString().padStart(2, '0'),
        "mm": this.getMinutes().toString().padStart(2, '0'),
        "ss": this.getSeconds().toString().padStart(2, '0'),
        "SSS": this.getMilliseconds().toString().padStart(3, '0'),
    };
    return pattern.replace(/yyyy|yy|MM|dd|HH|mm|ss|SSS/g, match => object[match]);
}

/**
 * 格式化字节数
 * @param n    字节数
 * @param prec 小数位数
 * @returns {string}
 */
function formatBytes(n, prec) {
    // 1B  = 8b (1 Byte = 8 bit)
    // 1KB = 1024B
    // 1MB = 1024KB
    // 1GB = 1024MB
    // 1TB = 1024GB

    if (n <= 0) {
        return "0 B"
    }

    let gb = n / (1024 * 1024 * 1024)
    if (gb > 1) {
        // 四舍五入保留小数点后 n 位
        return gb.toFixed(prec) + ' GB'
    }

    let mb = n / (1024 * 1024)
    if (mb > 1) {
        // 四舍五入保留小数点后 n 位
        return mb.toFixed(prec) + ' MB'
    }

    let kb = n / 1024
    if (kb > 1) {
        // 四舍五入保留小数点后 n 位
        return kb.toFixed(prec) + ' KB'
    }

    return n + ' B'
}

/**
 * 格式化毫秒
 * @param millisecond
 * @returns {string}
 */
function formatMillisecond(millisecond) {
    // 1 s = 1000 ms
    // 1 m = 60 s
    // 1 h = 60 m

    if (millisecond <= 0) {
        return "0 ms"
    }

    let hour = millisecond / (60 * 60 * 1000)
    if (hour > 1) {
        return hour.toFixed(2) + ' h'
    }

    let minute = millisecond / (60 * 1000)
    if (minute > 1) {
        return minute.toFixed(2) + ' m'
    }

    let second = millisecond / 1000
    if (second > 1) {
        return second.toFixed(2) + ' s'
    }

    return millisecond + ' ms'
}

/**
 * Y轴维度
 */
const YAxis = {
    // Y轴左侧
    Left: 'yLeft',
    // Y轴右侧
    Right: 'yRight',
};

/**
 * uPlot Line -- 折线图
 * https://github.com/leeoniya/uPlot
 * @param element 要渲染的元素
 * @param title   图表标题
 * @param width   图表宽度，单位：像素
 * @param height  图表高度，单位：像素
 * @param series  n 个系列
 * @constructor
 */
function Line(element, title, width, height, series) {
    // 数据：[ [X轴数据（时间戳（毫秒数）或数值）], [第一个系列数据集], [第二个系列数据集], ..., [第n个系列数据集] ]
    let data = Array.from({length: 1 + series.length}, () => new Array());

    // 配置选项
    let options = {
        // 图表标题
        title: title,
        // 图表宽度，单位：像素
        width: width,
        // 图表高度，单位：像素
        height: height,
        // Y轴刻度
        scales: {},
        // 坐标轴样式
        axes: [
            // X轴样式
            {
                // 格式化值
                values: (u, vals) => vals.map(v => uPlot.fmtDate('{HH}:{mm}')(new Date(v))),
            },
        ],
        // 坐标轴系列
        series: [
            // X轴系列
            {
                // 标签名称
                label: '时间',
                // 是否为时间轴
                // uPlot 要求时间戳必须是 JavaScript 时间戳（毫秒数）
                time: true,
                // 格式化值
                value: (u, v) => v === null ? '--' : uPlot.fmtDate('{YYYY}/{MM}/{DD} {HH}:{mm}:{ss}')(new Date(v)),
            },
        ],
    };

    let index = 0;
    series.forEach(ser => {
        // 数据系列线条颜色
        ser.stroke = Line.strokes[index++];

        if (ser.scale === YAxis.Left && options.scales[YAxis.Left] === undefined) {
            // Y轴左侧刻度
            options.scales[YAxis.Left] = {
                // 格式化值
                values: ser.formats,
            };

            // Y轴左侧样式
            options.axes.push({
                // 指定使用Y轴左侧
                scale: YAxis.Left,
                // 左侧样式
                side: 3,
                // 轴占用空间，单位：像素
                space: 50,
                // 格式化值
                values: ser.formats,
            });
        } else if (ser.scale === YAxis.Right && options.scales[YAxis.Right] === undefined) {
            // Y轴右侧刻度
            options.scales[YAxis.Right] = {
                // 格式化值
                values: ser.formats,
            };

            // Y轴右侧样式
            options.axes.push({
                // 指定使用Y轴右侧
                scale: YAxis.Right,
                // 右侧样式
                side: 1,
                // 轴占用空间，单位：像素
                space: 50,
                // 格式化值
                values: ser.formats,
            });
        }
        ser.value = ser.format;
        options.series.push(ser);
    });

    // 创建图表
    this.chart = new uPlot(options, data, element);

    // 最大保留的数据点数
    // 推荐每个数据点至少占据 1.5~2 个像素
    // 最佳数据点数 = width / 1.5
    this.maxPoints = Math.round(width / 1.5); // Math.round 将数字四舍五入为最接近的整数
    // console.log('maxPoints', this.maxPoints);
}

Line.strokes = [
    '#FF0000', // 红色
    '#0000FF', // 蓝色
    '#00FF00', // 绿色
    '#FFA500', // 橙色
    '#800080', // 紫色
    '#FF00FF', // 洋红色
    '#00FFFF', // 青色
    '#B19CD9', // 薰衣草紫
    '#FFB347', // 蜜橙
    '#6A5ACD', // 石板蓝
];

/**
 * 添加数据点
 * @param time    时间戳（毫秒数）
 * @param values  n 个系列数据
 */
Line.prototype.push = function (time, ...values) {
    // console.log('push', uPlot.fmtDate('{YYYY}/{MM}/{DD} {HH}:{mm}:{ss}')(new Date(time)), values);

    // 获取当前数据
    let data = this.chart.data;

    // 添加数据点
    data[0].push(time);
    let length = values.length;
    for (let i = 0; i < length; i++) {
        data[i + 1].push(values[i]);
    }

    // 只保留最近 maxPoints 个数据点
    if (data[0].length > this.maxPoints) {
        console.log(`shift data[${data[0].length}]`);
        for (let i = 0; i < 1 + length; i++) {
            // 移除数组的第一个元素并返回该元素
            data[i].shift();
        }
    }

    // 更新图表
    this.chart.setData(data);
}
