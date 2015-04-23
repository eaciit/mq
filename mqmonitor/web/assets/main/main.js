(function () {
	'use strict';

	var Main = function () { 
		var self = this;
		var $body = $('body');
		var $sectionNodes = $body.find('.section-nodes');
		var ajaxPullInterval = {};
		var windowResizeTimeout = setTimeout(function () {}, 0);
		var $window = $(window);
		var isServerAlive = true;
		var ajaxPullDelay = 7;
		var notifyDelay = 5;

		// register ajax pull, to make grid shows realtime data
		var registerAjaxPullFor = function (what, data, success, error) {
			var doRequest = function () {
				$.ajax({
					url: '/data/' + what,
					data: $.extend(true, { isServerAlive: isServerAlive }, data),
					type: 'get',
					dataType: 'json'
				})
				.success(success)
				.error(error);
			};

			ajaxPullInterval['section-' + what.replace(/\//g, '-')] = setInterval(
				doRequest, 
				ajaxPullDelay * 1000
			);

			doRequest();
		};

		// initiate all components
		this.init = function () {
			// toastr init
			toastr.options.closeButton = true

			// notify about delay every some minutes
			setInterval(function () {
				toastr['error']('data refreshed every ' + ajaxPullDelay + ' seconds');
			}, 1000 * 60 * notifyDelay);

			// prepare section nodes grid
			$sectionNodes.find('.grid').kendoGrid({
				chartArea: {
					background: 'transparent'
				},
				dataSource: { 
					data: [], 
					pageSize: 5 
				},
				pageable: {
					pageSizes: [5, 10, 15]
				},
				sortable: true, 
				scrollable: false,
				columns: [
					{ title: 'Configuration', columns: [
						{ field: 'ConfigName', title: 'Host', width: 110,
							template: '#: ConfigName #:#: ConfigPort #' },
						{ field: 'ConfigRole', title: 'Role', width: 90 }
					] },
					{ title: 'Data', columns: [
						{ field: 'DataCount', title: 'Total', width: 90, format: '{0:N0}',
							attributes: { style: 'text-align: right;' } },
						{ field: 'DataSize', title: 'Size (in KB)', width: 90, format: '{0:N2}',
							attributes: { style: 'text-align: right;' } }
					] },
					{ title: 'Time', columns: [
						{ field: 'StartTime', title: 'Start Time', width: 100, 
							attributes: { style: 'text-align: center;' } },
						{ field: 'Duration', title: 'Duration', width: 100, 
							attributes: { style: 'text-align: right;' } }
					] }
				],
				dataBound: function () {
					$(this.element).find('td:contains(Master)').css({
						fontWeight: 'bold'
					});
				}
			});

			// prepare chart
			$sectionNodes.find('.chart').kendoChart({
				chartArea: {
					background: 'transparent'
				},
				transitions: false,
				dataSource: {
					data: []
				},
				seriesDefaults: {
					type: 'line',
					style: 'smooth',
					markers: {
						visible: true
					},
				},
				series: [
					{ field: 'TotalHost', name: 'Total Host', axis: 'TotalHost',
						color: '#99b433',
						markers: {
							background: '#99b433'
						}
					},
					{ field: 'TotalDataCount', name: 'Total Data Count', axis: 'TotalDataCount',
						color: '#ee1111',
						markers: {
							background: '#ee1111'
						}
					},
					{ field: 'TotalDataSize', name: 'Total Data Size', axis: 'TotalDataSize',
						color: '#ffc40d',
						markers: {
							background: '#ffc40d'
						}
					}
				],
				categoryAxis: {
					field: 'Time',
					axisCrossingValues: [0, 0, 100],
					majorGridLines: {
						color: '#F9F9F9'
					}
				},
				valueAxes: [
					{ name: 'TotalHost', title: { text: 'Total Host' }, 
						min: 0, max: 0,
						majorGridLines: {
							color: '#F9F9F9'
						} 
					},
					{ name: 'TotalDataCount', title: { text: 'Total Data Count' }, 
						min: 0, max: 0 },
					{ name: 'TotalDataSize', title: { text: 'Total Data Size' }, 
						min: 0, max: 0 }
				],
				tooltip: {
					visible: true,
					template: '#= series.name # at #: category # => #= value #'
				},
				legend: {
					position: 'bottom'
				}
			});
		};

		// register ajax pull, 
		// make data semi real time
		// interval changeable
		this.registerAjaxPull = function () {

			// prepare ajax pull for nodes,
			// return data which used in both node grid & chart
			registerAjaxPullFor('nodes', {
				seriesLimit: 6,
				seriesDelay: ajaxPullDelay
			}, function (res) {
				if (!res.success) {
					if (res.message === 'connection is shut down')
						isServerAlive = false;

					toastr["error"](res.message);
					return;
				}

				if (isServerAlive == false) {
					isServerAlive = true;
					toastr["success"]("connected to server");
				}

				var $grid = $sectionNodes.find('.grid').data('kendoGrid');
				var $chart = $sectionNodes.find('.chart').data('kendoChart');

				$grid.setDataSource(new kendo.data.DataSource({
					data: Lazy(res.data.grid).sortBy(function (d) { return -d.StartTime; }).toArray(),
					pageSize: $grid.dataSource.pageSize()
				}));

				// get max value of each series,
				// then use it as valueAxis.max of each series
				Lazy($chart.options.valueAxis).each(function(v) {
					var max = Lazy(res.data.chart).max(function (d) {
						return parseInt(d[v.name], 10)
					})[v.name];

					v.max = max + Math.ceil(max / 5);
				});

				// sort data using time ascending
				$chart.setDataSource(new kendo.data.DataSource({
					data: Lazy(res.data.chart).sortBy(function (d) { return d.TimeInt; }).toArray()
				}));

				$chart.redraw();
			}, function (a, b, c) {
				toastr["error"]("Error occured when fetching data for nodes")
			});
		};

		// register event listener
		this.registerEventListener = function () {

			// when browser resized, do some changes
			$window.on('resize', function () {
				clearTimeout(windowResizeTimeout);

				windowResizeTimeout = setTimeout(function () {

					// redraw chart
					$sectionNodes.find('.chart').data('kendoChart').redraw();
				}, 500);
			});
		};
	};

	// start the magic
	$(function () {
		var main = new Main();
		main.init();
		main.registerAjaxPull();
		main.registerEventListener();
	});
}());
