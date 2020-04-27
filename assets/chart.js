// const opts = {
//   name: "Percentage of Population",
//   data: csvData,
//   selectedData: "confirmed",
//   percentageOfPopulation: false,
//   yScale: 'linear',
//   width: 800,
//   height: 600
// };
function makeChart(opts) {  
    var selectedData = opts.selectedData || "confirmed";
    var percentageOfPopulation = opts.percentageOfPopulation || false;
    var yScale = opts.yScale || 'linear';
    var width = opts.width || 1200;
    var height = opts.height || 600;

    // Set the dimensions of the canvas / graph
    var margin = {top: 60, right: 300, bottom: 60, left: 60},
    width = width - margin.left - margin.right,
    height = height - margin.top - margin.bottom;


    function getPercentage(d, selectedData) {
        return (d[selectedData]/d.population)*100
    }

    function getSelected(d, selectedData) {
        return d[selectedData]
    }

    function chooseData(d, selectedData) {

        if (selectedData === "active") {
            console.log(d.label_en, d.confirmed, d.recovered, d.deaths)
            var activeCases = d.confirmed - d.recovered - d.deaths
            console.log(activeCases)
            if (percentageOfPopulation){
                return (activeCases/d.population)*100
            } else {
                return activeCases
            }
        }

        if (percentageOfPopulation){
            return getPercentage(d, selectedData)
        } else {
            return getSelected(d, selectedData)
        }
    }

    // Parse the date / time
    var parseDate = d3.timeParse("%Y%m%d");

    // Set the ranges
    var x = d3.scaleTime().range([0, width]);  
    var y = d3.scaleSqrt().range([height, 0]);
    var y = d3.scaleLinear().range([height, 0]);

    // Define the line
    var line = d3.line()	
        .x(function(d) { return x(d.date); })
        .y(function(d) { return y(chooseData(d, selectedData)); });

    var container = d3.select("body")
        .append("div")
        .attr("class", "container")
        .attr("id", opts.id);

    container.append("h2").text(opts.name)

    // Adds the svg canvas
    var svg = container
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    function render(data){

        data = data.filter(function(row) {
            return row['label_parent_en'] === 'null';
        });

        data.forEach(function(d) {
            d.date = parseDate(d.date);
            d.val = +chooseData(d, selectedData);
        });

        var max = function(data){
            return d3.max(data, function(d) { return d.val; })
        }

        // Scale the range of the data
        x.domain(d3.extent(data, function(d) { return d.date; }));
        y.domain([0, max(data)]);

        // set the colour scale
        var color = d3.scaleOrdinal(d3.schemeCategory10);

        // Nest the entries by country
        var dataNest = d3.nest()
            .key(function(d) {return d.label_en;})
            .entries(data);

        dataNest.sort(function (a, b){
            if (max(a.values) > max(b.values)) {return -1;} 
            else if (max(a.values) < max(b.values)) { return 1;} 
            else return 0;
        })

        dataNest = dataNest.slice(0, 20);

        legendSpace = height/dataNest.length;

        // Loop through each symbol / key
        dataNest.forEach(function(d,i) { 
            svg.append("path")
                .attr("class", "line")
                .style("stroke", function() {
                    return d.color = color(d.key); })
                .attr("d", line(d.values));

            // Add the Legend
            svg.append("text")
                .attr("y", (legendSpace/2)+i*legendSpace)  
                .attr("x", width + 20)
                .attr("class", "legend")
                .style("fill", function() {
                    return d.color = color(d.key); })
                .text(d.key); 

        });

    // Add the X Axis
    svg.append("g")
        .attr("class", "axis")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Add the Y Axis
    var axisLeft = d3.axisLeft(y)
    if (percentageOfPopulation) {
        axisLeft = d3.axisLeft(y).tickFormat(d => d + "%")
    } 

    svg.append("g")
        .attr("class", "axis")
        .call(axisLeft);
    }

    // Get the data
    d3.csv(opts.data).then(function(data) {
        render(data)
    });
}
  
function getHistoricalData() {
    if (location.hostname === "localhost" || location.hostname === "127.0.0.1") {
        return "data/historical.csv";
    }
    return "https://montanaflynn.github.io/covid-19/data/historical.csv";
}

document.addEventListener("DOMContentLoaded", function () {
    historicalData = getHistoricalData()
    makeChart({
        name: "Confirmed Cases",
        id: "confirmed",
        selectedData: "confirmed",
        data: historicalData
    });

    makeChart({
        name: "Confirmed Cases By Population",
        id: "confirmed_by_population",
        selectedData: "confirmed",
        percentageOfPopulation: true,
        data: historicalData
    });

    makeChart({
        name: "Deaths Cases",
        id: "deaths",
        selectedData: "deaths",
        data: historicalData
    });

    makeChart({
        name: "Deaths Cases By Population",
        id: "deaths_by_population",
        selectedData: "deaths",
        percentageOfPopulation: true,
        data: historicalData
    });

    makeChart({
        name: "Recovered Cases",
        id: "recovered",
        selectedData: "recovered",
        data: historicalData
    });

    makeChart({
        name: "Recovered Cases By Population",
        id: "recovered_by_population",
        selectedData: "active",
        percentageOfPopulation: true,
        data: historicalData
    });

    makeChart({
        name: "Active Cases",
        id: "active",
        selectedData: "active",
        data: historicalData
    });

    makeChart({
        name: "Active Cases By Population",
        id: "active_by_population",
        selectedData: "active",
        percentageOfPopulation: true,
        data: historicalData
    });
  });
  