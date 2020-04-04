// const opts = {
//   mapName: "Worldwide",
//   topojson: topoJSON,
//   projection: projection,
//   objectName: "countries",
//   propertyName: "name",
//   data: dataJSON,
//   dataName: "global",
//   width: 800,
//   height: 600
// };
function makeMap(opts) {
  if (!("topojson" in opts) && !("geojson" in opts)) {
    throw "missing topojson or geojson";
  }

  var selectedData = opts.selectedData || "active";
  var sortData = opts.sortData || "active";
  var sortDirection = opts.sortDirection || "desc";
  var searchFilters = {};

  const paths = d3.geoPath().projection(opts.projection);

  const map = d3
    .select("body")
    .append("div")
    .attr("class", "container")
    .attr("id", opts.id);

  var spinner = new Spinner().spin(map.node());

  function loaded(error, mapdata, data) {
    if (error) throw error;

    spinner.stop();

    var values = [];

    const columns = ["region", "confirmed", "recovered", "deaths", "active"];
    var rows = [];
    for (const key in data[opts.dataName]) {
      if (data[opts.dataName].hasOwnProperty(key)) {
        const d = data[opts.dataName][key];
        rows.push({
          region: key,
          confirmed: d.confirmed,
          recovered: d.recovered,
          deaths: d.deaths,
          active: d.active,
        });
        values.push(d[selectedData]);
      }
    }

    if (sortDirection === "asc") {
      rows.sort((a, b) => (a[sortData] > b[sortData] ? 1 : -1));
    } else {
      rows.sort((a, b) => (a[sortData] < b[sortData] ? 1 : -1));
    }

    const totalCases = values.reduce(function getSum(total, num) {
      return total + num;
    });

    let selectedDataName =
      selectedData.charAt(0).toUpperCase() + selectedData.slice(1);

    if (selectedData !== "deaths") {
      selectedDataName = `${selectedDataName} Cases`;
    }

    const totalCasesHTML = `<h2> ${
      opts.mapName
    } ${selectedDataName}: ${totalCases.toLocaleString()} </h2>`;

    var tooltip = map
      .append("div")
      .attr("class", "toolTip")
      .html(totalCasesHTML);

    var table = map.append("table");
    var thead = table.append("thead");

    thead
      .append("tr")
      .selectAll("th")
      .data(columns)
      .enter()
      .append("td")
      .attr("class", "no-padding")
      .append("input")
      .attr("value", function (column) {
        if (column in searchFilters) {
          return searchFilters[column];
        } else {
          return "";
        }
      })
      .on("keyup", function (column) {
        searchFilters[column] = d3.select(d3.event.target).property("value");
        loadRows(table, rows);
      });

    thead
      .append("tr")
      .selectAll("th")
      .data(columns)
      .enter()
      .append("th")
      .attr("align", "left")
      .html(function (column) {
        let selected = "";
        let arrow = "";

        if (column === sortData) {
          if (sortDirection === "asc") {
            arrow = "&#x25B2;";
          } else {
            arrow = "&#x25BC;";
          }
        }

        if (column === selectedData) {
          selected = "&#x2605;";
        }

        return `${
          column.charAt(0).toUpperCase() + column.slice(1)
        } ${selected} ${arrow}`;
      })
      .on("click", function (d) {
        if (sortDirection === "asc" && sortData === d) {
          sortDirection = "desc";
        } else if (sortData === d) {
          sortDirection = "asc";
        }

        if (d !== "region") {
          selectedData = d;
        }

        sortData = d;

        map.selectAll("*").remove();

        loaded(error, mapdata, data);
      });

    loadRows(table, rows);

    function loadRows(table, rows) {
      table.selectAll("tbody").remove();
      var tbody = table.append("tbody");
      var rows = tbody
        .selectAll()
        .data(rows)
        .data(
          rows.filter(function (row) {
            for (const column in searchFilters) {
              const filter = searchFilters[column];
              const columnValue = row[column];

              if (columnValue.substring) {
                if (!columnValue.toLowerCase().includes(filter.toLowerCase())) {
                  return false;
                }
              } else {
                // do other thing
                if (parseInt(filter) > columnValue) {
                  return false;
                }
              }
            }
            return true;
          })
        )
        .enter()
        .append("tr");

      var cells = rows
        .selectAll("td")
        .data(function (row) {
          return columns.map(function (column) {
            return { value: row[column] };
          });
        })
        .enter()
        .append("td")
        .text(function (d) {
          return d.value.toLocaleString();
        });
    }

    let geometries;

    if ("topojson" in opts) {
      geometries = mapdata.objects[opts.objectName].geometries;
    } else {
      geometries = mapdata.features;
    }

    geometries.forEach(function (geo, i) {
      let locationName = geo.properties[opts.propertyName];
      if (locationName === "Ningxia Hui") {
        locationName = "Ningxia";
      } else if (locationName === "Xizang") {
        locationName = "Tibet";
      } else if (locationName === "Nei Mongol") {
        locationName = "Inner Mongolia";
      } else if (locationName === "United States of America") {
        locationName = "United States";
      } else if (locationName === "Thành phố Hồ Chí Minh") {
        locationName = "Hồ Chí Minh";
      } else if (locationName === "Quebec") {
        locationName = "Québec";
      } else if (locationName === "Bosnia and Herz.") {
        locationName = "Bosnia and Herzegovina";
      } else if (locationName === "Macedonia") {
        locationName = "North Macedonia";
      }

      if (locationName in data[opts.dataName]) {
        const d = data[opts.dataName][locationName];
        geometries[i].properties.name = locationName;
        geometries[i].properties.confirmed = d.confirmed;
        geometries[i].properties.recovered = d.recovered;
        geometries[i].properties.deaths = d.deaths;
        geometries[i].properties.active = d.active;
      } else {
        geometries[i].properties.name = locationName;
        geometries[i].properties.confirmed = 0;
        geometries[i].properties.recovered = 0;
        geometries[i].properties.deaths = 0;
        geometries[i].properties.active = 0;
      }
    });

    const minVal = d3.min(values);
    const maxVal = d3.max(values);

    const lowColor = "#fee";
    const highColor = "#f00";

    const color = d3
      .scaleSqrt()
      .domain([0, 1, maxVal])
      .range(["green", lowColor, highColor]);

    let pathData;
    if ("topojson" in opts) {
      pathData = topojson.feature(mapdata, mapdata.objects[opts.objectName]);
    } else {
      pathData = mapdata;
    }

    const svg = map
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
      .attr("fill", function (d, i) {
        return color(d.properties[selectedData]);
      })
      .attr("d", paths)
      .on("mouseover", function (d) {
        var currentState = this;
        d3.select(this).style("stroke-width", 1.5);
        tooltip.html(
          `<h2> ${d.properties.name} ${selectedDataName}: ${d.properties[
            selectedData
          ].toLocaleString()} </h2>`
        );
      })
      .on("mouseout", function (d) {
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

  var jsonFiles = [d3.json(mapjson), d3.json(opts.data)];

  return new Promise((resolve, reject) => {
    Promise.all(jsonFiles)
      .then(function (values) {
        loaded(null, values[0], values[1]);
      })
      .then(function () {
        return resolve();
      })
      .catch(function (err) {
        return reject(err);
      });
  });
}

function worldMap() {
  const w = 800;
  const h = 583;

  const projection = d3
    .geoMercator()
    .scale(w / Math.PI / 2)
    .translate([w / 2, w / 2]);

  return makeMap({
    mapName: "Worldwide",
    id: "world",
    width: w,
    height: h,
    topojson: "https://cdn.jsdelivr.net/npm/world-atlas@2/countries-50m.json",
    data: "https://montanaflynn.github.io/covid-19/data/current.json",
    projection: projection,
    objectName: "countries",
    dataName: "global",
    propertyName: "name",
  });
}

function usaMap() {
  const w = 800;
  const h = 600;

  const projection = d3
    .geoAlbersUsa()
    .translate([w / 2, h / 2])
    .scale([1000]);

  return makeMap({
    mapName: "United States",
    id: "usa",
    width: w,
    height: h,
    topojson: "https://cdn.jsdelivr.net/npm/us-atlas@3/states-10m.json",
    data: "https://montanaflynn.github.io/covid-19/data/current.json",
    projection: projection,
    objectName: "states",
    dataName: "usa",
    propertyName: "name",
  });
}

function canadaMap() {
  const w = 800;
  const h = 600;

  const projection = d3
    .geoAzimuthalEqualArea()
    .rotate([100, -45])
    .center([5, 20])
    .scale(w)
    .translate([w / 2, h / 2]);

  return makeMap({
    mapName: "Canada",
    id: "canada",
    width: w,
    height: h,
    topojson:
      "https://gistcdn.githack.com/montanaflynn/32f882ec77b0dd15bced6a28fad80028/raw/13f1fb4d257ca2f11dd441003ce578a46ec5097f/canada-provinces.topo.json",
    data: "https://montanaflynn.github.io/covid-19/data/current.json",
    projection: projection,
    objectName: "provinces",
    dataName: "canada",
    propertyName: "name",
  });
}

function germanyMap() {
  const w = 800;
  const h = 600;

  var projection = d3
    .geoMercator()
    .center([10.5, 51.35])
    .scale(2000)
    .translate([w / 2, h / 2]);

  return makeMap({
    mapName: "Germany",
    id: "germany",
    width: w,
    height: h,
    geojson:
      "https://raw.githubusercontent.com/isellsoap/deutschlandGeoJSON/master/2_bundeslaender/4_niedrig.geojson",
    data: "https://montanaflynn.github.io/covid-19/data/current.json",
    projection: projection,
    objectName: "DEU_adm2",
    dataName: "germany",
    propertyName: "NAME_1",
  });
}

function chinaMap() {
  const w = 800;
  const h = 600;

  const projection = d3
    .geoMercator()
    .center([110, 25])
    .scale([700])
    .translate([450, 500]);

  return makeMap({
    mapName: "China",
    id: "china",
    width: w,
    height: h,
    topojson:
      "https://raw.githubusercontent.com/deldersveld/topojson/master/countries/china/china-provinces.json",
    data: "https://montanaflynn.github.io/covid-19/data/current.json",
    projection: projection,
    objectName: "CHN_adm1",
    dataName: "china",
    propertyName: "NAME_1",
  });
}

function vietnamMap() {
  const w = 800;
  const h = 600;

  var projection = d3
    .geoMercator()
    .center([108.5, 14.35])
    .scale(2200)
    .translate([w / 2, h / 2 + 70]);

  return makeMap({
    mapName: "Vietnam",
    id: "vietnam",
    width: w,
    height: h,
    topojson:
      "https://raw.githubusercontent.com/kcjpop/vietnam-topojson/master/adm2/adm2.json",
    data: "https://montanaflynn.github.io/covid-19/data/current.json",
    projection: projection,
    objectName: "adm2",
    dataName: "vietnam",
    propertyName: "name_vi",
  });
}

document.addEventListener("DOMContentLoaded", function () {
  var promises = [
    worldMap(),
    usaMap(),
    canadaMap(),
    germanyMap(),
    chinaMap(),
    vietnamMap(),
  ];

  Promise.all(promises).catch(function (err) {
    throw err;
  });
});
