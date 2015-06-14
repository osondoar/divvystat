(function(window){

  var d3ChartElement;
  var lineChart
  initialize();

  function initialize(){
    d3ChartElement = d3.select('#chart svg');

    initializeDatePicker();
    loadData(function(statuses){
      drawD3Chart(statuses);
    });
  }

  function initializeDatePicker(){
    timeZoneData = "America/Chicago|CST CDT EST CWT CPT|60 50 50 50 50|01010101010101010101010101010101010102010101010103401010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010|-261s0 1nX0 11B0 1nX0 1wp0 TX0 WN0 1qL0 1cN0 WL0 1qN0 11z0 1o10 11z0 1o10 11z0 1o10 11z0 1o10 11z0 1qN0 11z0 1o10 11z0 1o10 11z0 1o10 11z0 1o10 11z0 1qN0 WL0 1qN0 11z0 1o10 11z0 11B0 1Hz0 14p0 11z0 1o10 11z0 1qN0 WL0 1qN0 11z0 1o10 11z0 RB0 8x30 iw0 1o10 11z0 1o10 11z0 1o10 11z0 1o10 11z0 1qN0 WL0 1qN0 11z0 1o10 11z0 1o10 11z0 1o10 11z0 1o10 1fz0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1fz0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1fz0 1a10 1fz0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1fz0 1cN0 1cL0 1cN0 1cL0 s10 1Vz0 LB0 1BX0 1cN0 1fz0 1a10 1fz0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 1cN0 1fz0 1a10 1fz0 1cN0 1cL0 1cN0 1cL0 1cN0 1cL0 14p0 1lb0 14p0 1nX0 11B0 1nX0 11B0 1nX0 14p0 1lb0 14p0 1lb0 14p0 1nX0 11B0 1nX0 11B0 1nX0 14p0 1lb0 14p0 1lb0 14p0 1lb0 14p0 1nX0 11B0 1nX0 11B0 1nX0 14p0 1lb0 14p0 1lb0 14p0 1nX0 11B0 1nX0 11B0 1nX0 Rd0 1zb0 Op0 1zb0 Op0 1zb0 Rd0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Rd0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Rd0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Rd0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Rd0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0 Op0 1zb0"
    moment.tz.add(timeZoneData);
    jQuery('#datetimepicker-from').datetimepicker();
    jQuery('#datetimepicker-to').datetimepicker();

    $('#datetimepicker-from, #datetimepicker-to').on('change', function() {
      loadData(function(data){
        d3ChartElement.datum(data);
        lineChart.update();
      });
    });
  }

  function loadData(callback){
    var from = new Date( $('#datetimepicker-from').val());
    var to = new Date( $('#datetimepicker-to').val());
    var fromIso8601 = from != 'Invalid Date' ? moment.tz(from, "America/Chicago").format() : "";
    var toIso8601 = to != 'Invalid Date' ? moment.tz(to, "America/Chicago").format() : ""
    $.get( "api/loads?from=" + fromIso8601 + '&to='+toIso8601, function( statuses ) {
      var data = buildDataForChart(statuses)
      callback(data);
    });
  }

  function buildDataForChart(data) {
    var loads = [];

    for (var status of data) {
      // var date = parseDate(status['execution_time']); // Use this instead if there are issues parsing dates across browsers
      var date = new Date(status['execution_time']);
      loads.push({ x: date, y: status['load'] });
    }
    return [
      {
        values: loads,
        key: 'Divvy Activity',
        color: '#3db7e4'
      }
    ];
  }

  function drawD3Chart(data){
    nv.addGraph(function() {
      lineChart = nv.models.lineChart()
                .margin({left: 60})  //Adjust lineChart margins to give the x-axis some breathing room.
                .useInteractiveGuideline(false)  //It doesn't seem to work with dates. We want nice looking tooltips and a guideline!
                .showLegend(true)       //Show the legend, allowing users to turn on/off line series.
                .showYAxis(true)        //Show the y-axis
                .showXAxis(true)        //Show the x-axis
                .tooltipContent(function (key, x, y, e) {
                                return  '<span class="time">' + x + '</span>' +
                                        '<div class="activity">' +
                                          '<span class="bullet"></span>' +
                                          '<span>Divvy Activity: ' + e.point.y + '</span>' +
                                        '</div>';
                            })

                            lineChart.xAxis
        // .axisLabel('Time (ms)')
        .tickFormat(function(d) { return d3.time.format('%H:%M, %b %d')(new Date(d)); })


        lineChart.yAxis     //Chart y-axis settings
        .axisLabel('Activity');


        // d3.time.scale().ticks(d3.time.minute, 10);
        lineChart.xScale(d3.time.scale());
        // lineChart.xAxis.ticks(d3.time.minutes,10);

        d3ChartElement            //Select the <svg> element you want to render the lineChart in.
          .datum(data)         //Populate the <svg> element with lineChart data...
          .call(lineChart);        //Finally, render the lineChart!

      //Update the lineChart when window resizes.
      nv.utils.windowResize(function() { lineChart.update() });
      return lineChart;
    });
  }

  // function getParameterByName(name) {
  //     name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
  //     var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
  //         results = regex.exec(location.search);
  //     return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
  // }


})(window);
