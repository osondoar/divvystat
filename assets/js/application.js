(function(window){

  google.load("visualization", "1", {packages:["annotatedtimeline"]});
  google.setOnLoadCallback(laodData);


  function laodData(){
    $.get( "api/loads", function( statuses ) {
      var statusesArray  = new Array()
        for (var key in statuses) {

        var time = new Date(key);
        statusesArray.push([time, statuses[key]]);
      }
      console.log(statusesArray);
      drawChart(statusesArray, 30);
    });
  }

  function drawChart(status_data, total_docks){
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn('datetime', 'Time');
    dataTable.addColumn('number', 'Divvy Load');
    dataTable.addRows(status_data);

    var dataView = new google.visualization.DataView(dataTable);
    var options = {
      title: 'Divvy Load',
      max: total_docks
    };
    var chart = new google.visualization.AnnotatedTimeLine(document.getElementById('chart_div'));
    chart.draw(dataView, options);
  }

})(window);
