(function () {
  'use strict'

  const leftServer = document.getElementById('left-server')
  const rightServer = document.getElementById('right-server')
  const subtrahend = document.getElementById('subtrahend')

  const makeColor = () => {
    const octet = () => Math.floor(Math.random() * 16).toString(16)
    return '#' + octet() + octet() + octet()
  }

  const chart = new Chart(document.getElementById('chart'), {
    type: 'bar',
    data: {
      datasets: [{
        backgroundColor: [makeColor(), makeColor()],
        data: [0, 0],
        label: 'weight'
      }],
      labels: ['', '']
    },
    options: {
      legend: { display: false },
      scales: {
        yAxes: [{
          ticks: {
            beginAtZero: true,
            suggestedMax: 256
          }
        }]
      }
    }
  })

  const get = () => {
    const xhr = new XMLHttpRequest()

    xhr.addEventListener('load', () => {
      const data = JSON.parse(xhr.responseText)

      chart.data.datasets[0].data = [256 - data.subtrahend, data.subtrahend]
      chart.update()

      if (leftServer !== document.activeElement) leftServer.value = data.servers[0]
      if (rightServer !== document.activeElement) rightServer.value = data.servers[1]
      if (subtrahend !== document.activeElement) subtrahend.value = data.subtrahend

      leftServer.disabled = data.subtrahend !== 256
      rightServer.disabled = data.subtrahend !== 0
    })

    xhr.open('GET', '/crossfader')
    xhr.send()
  }
  get()
  setInterval(get, 1000)

  const put = () => {
    const xhr = new XMLHttpRequest()
    xhr.open('PUT', '/crossfader')
    xhr.send(JSON.stringify({
      servers: [leftServer.value, rightServer.value],
      subtrahend: parseInt(subtrahend.value)
    }))
  }
  leftServer.addEventListener('blur', put)
  rightServer.addEventListener('blur', put)
  subtrahend.addEventListener('input', put)
  subtrahend.addEventListener('change', put)
})()
