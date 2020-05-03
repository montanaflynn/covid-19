// const opts = {
//   name: "Percentage of Population",
//   data: csvData,
//   selectedData: "confirmed",
//   percentageOfPopulation: false,
//   yScale: 'linear',
//   width: 800,
//   height: 600
// }
function makeChart(opts) {  
    var selectedData = opts.selectedData || "confirmed"
    var percentageOfPopulation = opts.percentageOfPopulation || false
    var width = opts.width || window.innerWidth-300 || 1200
    var height = opts.height || window.innerHeight-300 || 600
    var originalData = opts.data
    var chartName = opts.name || "Confirmed Cases"
    var yScale = getScaleStyle(opts.yScale)

    // Set the dimensions of the canvas / graph
    var margin = {top: 10, right: 300, bottom: 60, left: 60},
    width = width - margin.left - margin.right,
    height = height - margin.top - margin.bottom

    function getChartName() {
        switch(selectedData) {
            case "recovered":
                chartName = "Recovered Cases"
                break;
            case "recovered":
                chartName = "Recovered Cases"
                break;
            case "deaths":
                chartName = "Deaths"
                break;
            case "active":
                chartName = "Active Cases"
                break;
            default:
                chartName = "Confirmed Cases"
            }
        return percentageOfPopulation ? chartName + " as Percentage of Population" : chartName 
    }

    function getScaleStyle(scaleString) {
        console.log(scaleString)
        switch(scaleString) {
            case "sqrt":
                return d3.scaleSqrt()
            case "pow":
                return d3.scalePow()
            case "log":
                return d3.scaleSymlog()
            default:
                return d3.scaleLinear()
            }
    }

    function getPercentage(d, selectedData) {
        return (d[selectedData]/d.population)*100
    }

    function getSelected(d, selectedData) {
        return d[selectedData]
    }

    function chooseData(d, selectedData) {
        if (percentageOfPopulation){
            return getPercentage(d, selectedData)
        } else {
            return getSelected(d, selectedData)
        }
    }

    // Parser for the date / time
    var parseDate = d3.timeParse("%Y%m%d")

    function render(data){

        d3.selectAll(".container").remove()

        var data = JSON.parse(JSON.stringify(data))

        var container = d3.select("body")
        .append("div")
        .attr("class", "container")

        container.append("h2").text(getChartName())

        // Adds the svg canvas
        var svg = container
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")")

        data = data.filter(function(row) {
            return row['label_parent_en'] === 'null'
        })

        data.forEach(function(d) {
            d.date = parseDate(d.date)
            d.val = +chooseData(d, selectedData)
        })

        var max = function(data){
            return d3.max(data, function(d) { return d.val })
        }

        // Set the ranges
        var x = d3.scaleTime().range([0, width])  
        var y = yScale.range([height, 0])

        // Define the line
        var line = d3.line()	
        .x(function(d) { return x(d.date) })
        .y(function(d) { return y(chooseData(d, selectedData)) })
        

        // Scale the range of the data
        x.domain(d3.extent(data, function(d) { return d.date }))
        y.domain([0, max(data)])

        // set the colour scale
        var color = d3.scaleOrdinal(d3.schemeCategory10)

        // Nest the entries by country
        var dataNest = d3.nest()
            .key(function(d) {return d.label_en})
            .entries(data)

        dataNest.sort(function (a, b){
            if (max(a.values) > max(b.values)) {return -1} 
            else if (max(a.values) < max(b.values)) { return 1} 
            else return 0
        })

        dataNest = dataNest.slice(0, 20)

        legendSpace = height/dataNest.length

        // Loop through each symbol / key
        dataNest.forEach(function(d,i) { 
            svg.append("path")
                .attr("class", "line")
                .style("stroke", function() {
                    return d.color = color(d.key) })
                .attr("d", line(d.values))

            // Add the Legend
            svg.append("text")
                .attr("y", (legendSpace/2)+i*legendSpace)  
                .attr("x", width + 20)
                .attr("class", "legend")
                .style("fill", function() {
                    return d.color = color(d.key) })
                .text(d.key) 
        })

        // Add the X Axis
        svg.append("g")
            .attr("class", "axis")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x))

        // Add the Y Axis
        var axisLeft = d3.axisLeft(y)
        if (percentageOfPopulation) {
            axisLeft = d3.axisLeft(y).tickFormat(d => d + "%")
        } 

        svg.append("g")
            .attr("class", "axis")
            .call(axisLeft.tickArguments([5]))
    
    }

    const checkbox = d3.select("#by_population")
    checkbox.on("change", function(){
        percentageOfPopulation = d3.select(this).property("checked")
        render(originalData)
    })

    const metricRadioButtons = d3.selectAll('input[name="metric"]')
    metricRadioButtons.on("change", function(){
        selectedData = this.value
        render(originalData)
    })

    const scaleRadioButtons = d3.selectAll('input[name="scale"]')
    scaleRadioButtons.on("change", function(){
        yScale = getScaleStyle(this.value)
        render(originalData)
    })

    render(originalData)
}
  
function getHistoricalDataURL() {
    if (location.hostname === "localhost" || location.hostname === "127.0.0.1") {
        return "data/historical.csv"
    }
    return "https://montanaflynn.github.io/covid-19/data/historical.csv"
}

document.addEventListener("DOMContentLoaded", function () {

    // get the data
    d3.csv(getHistoricalDataURL()).then(function(data) {

        // add active cases property
        data.map(d => d.active = d.confirmed - d.recovered - d.deaths)
        makeChart({
            selectedData: "confirmed",
            data: data
        })
    })
  })
  