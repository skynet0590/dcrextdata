import { Controller } from 'stimulus'
// import { Date } from '../utils'
import axios from 'axios'

const Dygraph = require('../../../dist/js/dygraphs.min.js')
// const atomsToDCR = 1e-8

// function intComma (amount) {
//   return amount.toLocaleString(undefined, { maximumFractionDigits: 0 })
// }

// function legendFormatter (data) {
//   var html = ''
//   if (data.x == null) {
//     let dashLabels = data.series.reduce((nodes, series) => {
//       return `${nodes} <div class="pr-2">${series.dashHTML} ${series.labelHTML}</div>`
//     }, '')
//     html = `<div class="d-flex flex-wrap justify-content-center align-items-center">
//               <div class="pr-3">${this.getLabels()[0]}: N/A</div>
//               <div class="d-flex flex-wrap">${dashLabels}</div>
//             </div>`
//   } else {
//     data.series.sort((a, b) => a.y > b.y ? -1 : 1)
//     var extraHTML = ''
//     // The circulation chart has an additional legend entry showing percent
//     // difference.
//     if (data.series.length === 2 && data.series[1].label.toLowerCase() === 'coin supply') {
//       let predicted = data.series[0].y
//       let actual = data.series[1].y
//       let change = (((actual - predicted) / predicted) * 100).toFixed(2)
//       extraHTML = `<div class="pr-2">&nbsp;&nbsp;Change: ${change} %</div>`
//     }

//     let yVals = data.series.reduce((nodes, series) => {
//       if (!series.isVisible) return nodes
//       let yVal = series.yHTML
//       switch (series.label.toLowerCase()) {
//         case 'coin supply':
//           yVal = intComma(series.y) + ' DCR'
//           break

//         case 'hashrate':
//           yVal = formatHashRate(series.y)
//           break
//       }
//       return `${nodes} <div class="pr-2">${series.dashHTML} ${series.labelHTML}: ${yVal}</div>`
//     }, '')

//     html = `<div class="d-flex flex-wrap justify-content-center align-items-center">
//                 <div class="pr-3">${this.getLabels()[0]}: ${data.xHTML}</div>
//                 <div class="d-flex flex-wrap"> ${yVals}</div>
//             </div>${extraHTML}`
//   }

//   dompurify.sanitize(html)
//   return html
// }

// // function zipXYZData (gData, isHeightAxis, isDayBinned, yCoefficient, zCoefficient, windowS) {
// //   windowS = windowS || 1
// //   yCoefficient = yCoefficient || 1
// //   zCoefficient = zCoefficient || 1
// //   return map(gData.x, (n, i) => {
// //     var xAxisVal
// //     if (isHeightAxis && isDayBinned) {
// //       xAxisVal = n
// //     } else if (isHeightAxis) {
// //       xAxisVal = i * windowS
// //     } else {
// //       xAxisVal = new Date(n * 1000)
// //     }
// //     return [xAxisVal, gData.y[i] * yCoefficient, gData.z[i] * zCoefficient]
// //   })
// // }

export default class extends Controller {
  static get targets () {
    return [
      'chartWrapper',
      'labels',
      'chartsView',
      'axisOption'
    ]
  }

  connect () {
    // this.drawInitialGraph()
    // var windowSize = parseInt(this.data.get('windowSize'))
    this.plotGraph()
  }

  // drawInitialGraph () {
  //   var options = {
  //     axes: { y: { axisLabelWidth: 70 }, y2: { axisLabelWidth: 70 } },
  //     labels: ['Date', 'Pools rate', 'Pools hashrate'],
  //     digitsAfterDecimal: 8,
  //     showRangeSelector: true,
  //     rangeSelectorPlotFillColor: '#8997A5',
  //     rangeSelectorAlpha: 0.4,
  //     rangeSelectorHeight: 40,
  //     drawPoints: true,
  //     pointSize: 0.25,
  //     legend: 'always',
  //     labelsSeparateLines: true,
  //     labelsDiv: this.labelsTarget,
  //     legendFormatter: legendFormatter,
  //     highlightCircleSize: 4,
  //     ylabel: 'Pool Rate',
  //     y2label: 'Pools hashrate',
  //     labelsUTC: true
  //   }

  //   this.chartsView = new Dygraph(
  //     this.chartsViewTarget,
  //     [[1, 1, 10], [2, 5, 110]],
  //     options
  //   )
  // }

  plotGraph () {
    // var d = []
    // var gOptions = {
    //   zoomCallback: null,
    //   drawCallback: null,
    //   axes: {},
    //   visibility: null,
    //   y2label: null,
    //   stepPlot: false
    // }

    // var isHeightAxis = this.selectedAxis() === 'height'
    // var xlabel = isHeightAxis ? 'Block Height' : 'Date'
    // var isDayBinned = this.selectedBin() === 'day'

    this.nextPage = 1

    var options = {
      axes: { y: { axisLabelWidth: 70 }, y2: { axisLabelWidth: 70 } },
      labels: ['Date', 'Network Difficulty', 'pool hash'],
      digitsAfterDecimal: 2,
      showRangeSelector: true,
      rangeSelectorPlotFillColor: '#8997A5',
      rangeSelectorAlpha: 0.4,
      rangeSelectorHeight: 40,
      drawPoints: true,
      pointSize: 0.25,
      legend: 'always',
      labelsSeparateLines: true,
      highlightCircleSize: 4,
      ylabel: 'hash',
      y2label: 'diff',
      labelsUTC: true
    }

    const _this = this
    axios.get(`/getChartPowData?page=${this.nextPage}`)
      .then(function (response) {
        console.log(response.data)
        let result = response.data

        var data = []
        var dataSet = []
        result.powData.forEach(pow => {
          data.push(new Date(pow.Time))
          data.push(pow.PoolHashrate)
          data.push(pow.NetworkDifficulty)

          dataSet.push(data)
          data = []
        })
        console.log('...java Script Array... \n' + JSON.stringify(dataSet))
        _this.chartsView = new Dygraph(
          _this.chartsViewTarget,
          dataSet, options
        )
      }).catch(function (e) {
        console.log(e)
      })
  }
}
