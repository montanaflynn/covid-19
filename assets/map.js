// const opts = {
//     "mapName": "China",
//     "width": w,
//     "height": h,
//     "topojson": topoJSON,
//     "geojson": null,
//     "data": dataJSON,
//     "projection": projection,
//     "objectName":"CHN_adm1",
//     "dataName":"china",
//      "propertyName":"NAME_1"
//   }
function makeMap(opts) {
  if (!("topojson" in opts) && !("geojson" in opts)) {
    throw "missing topojson or geojson";
  }

  const paths = d3.geoPath().projection(opts.projection);

  d3.select("body")
    .append("div")
    .attr("class", "container");

  var targets = document.getElementsByClassName("container");
  var spinner = new Spinner().spin(targets[0]);

  function loaded(error, mapdata, data) {
    if (error) throw error;

    spinner.stop();

    let geometries;

    if ("topojson" in opts) {
      geometries = mapdata.objects[opts.objectName].geometries;
    } else {
      geometries = mapdata.features;
    }

    geometries.forEach(function(geo, i) {
      let locationName = geo.properties[opts.propertyName];
      if (locationName === "Ningxia Hui") {
        locationName = "Ningxia";
      } else if (locationName === "Xizang") {
        locationName = "Tibet";
      } else if (locationName === "United States of America") {
        locationName = "United States";
      }

      if (locationName in data[opts.dataName]) {
        geometries[i].properties.name = locationName;
        geometries[i].properties.confirmed =
          data[opts.dataName][locationName].confirmed;
      } else {
        geometries[i].properties.confirmed = 0;
      }
    });

    let values;
    if ("topojson" in opts) {
      values = d3
        .entries(mapdata.objects[opts.objectName].geometries)
        .map(function(d) {
          return d.value.properties.confirmed;
        });
    } else {
      values = d3.entries(mapdata.features).map(function(d) {
        return d.value.properties.confirmed;
      });
    }

    const totalCases = values.reduce(function getSum(total, num) {
      return total + num;
    });

    const totalCasesHTML =
      "<h2>" + opts.mapName + " Confirmed Cases: " + totalCases + "</h2>";

    const minVal = d3.min(values);
    const maxVal = d3.max(values);

    const lowColor = "#fee";
    const highColor = "#f00";

    const color = d3
      .scaleSqrt()
      .domain([0, 1, maxVal])
      .range(["#fff", lowColor, highColor]);

    var tooltip = d3
      .select(".container")
      .append("div")
      .attr("class", "toolTip")
      .html(totalCasesHTML);

    let pathData;
    if ("topojson" in opts) {
      pathData = topojson.feature(mapdata, mapdata.objects[opts.objectName]);
    } else {
      pathData = mapdata;
    }

    const svg = d3
      .select(".container")
      .append("svg")
      .attr("width", opts.width)
      .attr("height", opts.height)
      .append("g")
      .selectAll("path")
      .data(pathData.features)
      .enter()
      .append("path")
      .attr("stroke", "#000")
      .attr("stroke-width", 0.5)
      .attr("fill", function(d, i) {
        return color(d.properties.confirmed);
      })
      .attr("d", paths)
      .on("mouseover", function(d) {
        var currentState = this;
        d3.select(this).style("stroke-width", 1.5);
        tooltip.html(
          "<h2>" +
            d.properties.name +
            " Confirmed Cases: " +
            d.properties.confirmed +
            "</h2>"
        );
      })
      .on("mouseout", function(d) {
        d3.select(this).style("stroke-width", 0.5);
        tooltip.html(totalCasesHTML);
      });
  }

  let mapjson;
  if ("topojson" in opts) {
    mapjson = opts.topojson;
  } else {
    mapjson = opts.geojson;
  }

  console.log(mapjson);

  d3.queue()
    .defer(d3.json, mapjson)
    .defer(d3.json, opts.data)
    .await(loaded);
}
