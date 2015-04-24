(function () {
	'use strict';

	var Main = function () { 
		var self = this;
		var $body = $('body');
		var $sectionNodesGrid = $body.find('.section-nodes-grid');
		var $sectionNodesChart = $body.find('.section-nodes-chart');
		var $sectionItemsGrid = $body.find('.section-items-grid');
		var windowResizeTimeout = setTimeout(function () {}, 0);
		var $window = $(window);
		var isServerAlive = true;
		var ajaxPullDelay = 7;
		var notifyDelay = 5;
		var seriesLimit = 6;
		var dataSizeUnit = 1024 * 1024;

		// register ajax pull, to make grid shows realtime data
		var registerAjaxPullFor = function (what, data, success, error, after) {
			success = (typeof success !== String(undefined)) ? success : function () {};
			error 	= (typeof error   !== String(undefined)) ? error   : function () {};
			after 	= (typeof after   !== String(undefined)) ? after   : function () {};

			data = $.extend(true, { 
				isServerAlive: isServerAlive
			}, data);

			var doRequest = function () {
				$.ajax({
					url: '/data/' + what,
					data: data,
					type: 'get',
					dataType: 'json'
				})
				.success(success)
				.error(error);
			};

			after(doRequest);
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

			// prepare section nodes, grid
			$sectionNodesGrid.find('.grid').kendoGrid({
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
						{ field: 'ConfigHost', title: 'Host', width: 100 },
						{ field: 'ConfigRole', title: 'Role', width: 60 }
					] },
					{ title: 'Data', columns: [
						{ field: 'DataCount', title: 'Total', width: 80,
							format: '{0:N0}', attributes: { style: 'text-align: right;' } },
						{ field: 'DataSize', title: 'Size<br />(in byte)', width: 80,
							format: '{0:N2}', attributes: { style: 'text-align: right;' } },
						{ field: 'AllocatedSize', title: 'Allocated<br />(in Mb)', width: 80,
							format: '{0:N2}', attributes: { style: 'text-align: right;' } }
					] },
					{ title: 'Time', columns: [
						{ field: 'StartTime', title: 'Start Time', width: 90,
							attributes: { style: 'text-align: center;' } },
						{ field: 'Duration', title: 'Duration', width: 80, 
							attributes: { style: 'text-align: center;' } }
					] }
				],
				dataBound: function () {
					$(this.element).find('td:contains(Master)').css({
						fontWeight: 'bold'
					});
				}
			});

			// prepare section nodes, chart
			$sectionNodesChart.find('.chart').kendoChart({
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
						markers: { background: '#99b433' }
					},
					{ field: 'TotalDataCount', name: 'Total Data Count', axis: 'TotalDataCount',
						markers: { background: '#ee1111' }
					},
					{ field: 'TotalDataSize', name: 'Total Data Size', axis: 'TotalDataSize',
						markers: { background: '#ffc40d' }
					},
					{ field: 'TotalAllocatedSize', name: 'Total Allocated Size', axis: 'TotalAllocatedSize',
						markers: { background: '#337ab7' }
					}
				],
				seriesColors: ['#99b433', '#ee1111', '#ffc40d', '#337ab7'],
				categoryAxis: {
					field: 'Time',
					axisCrossingValues: [0, 0, 100, 100],
					majorGridLines: {
						color: '#F9F9F9'
					}
				},
				valueAxes: [
					{ name: 'TotalHost', title: { text: 'Total Host' }, min: 0,
						majorGridLines: {
							color: '#F9F9F9'
						} 
					},
					{ name: 'TotalDataCount', title: { text: 'Total Data Count' }, 
						min: 0 },
					{ name: 'TotalDataSize', title: { text: 'Total Data Size' }, 
						min: 0 },
					{ name: 'TotalAllocatedSize', title: { text: 'Total Allocated Size' }, 
						min: 0 }
				],
				tooltip: {
					visible: true,
					template: '#= series.name # at #: category # => #= value #'
				},
				legend: {
					position: 'bottom'
				}
			});

			// prepare section items, grid
			$sectionItemsGrid.find('.grid').kendoGrid({
				dataSource: { 
					data: [], 
					pageSize: 10
				},
				pageable: {
					pageSizes: [5, 10, 15]
				},
				sortable: true, 
				scrollable: false,
				columns: [
					{ field: 'Key', title: 'Key' },
					{ field: 'Value', title: 'Value' },
					{ field: 'Created', title: 'Created', width: 90,
						attributes: { style: 'text-align: center;' } },
					{ field: 'Expiry', title: 'Expiry', width: 80,
						attributes: { style: 'text-align: center;' } },
				]
			});
			
			/* // prepare section items, node selection
			$sectionItemsGrid.find('select.nodes').kendoDropDownList({
				dataSource: { 
					data: [
						{ text: 'Select one ...', value: '' }
					]
				},
				dataTextField: 'text',
				dataValueField: 'value',
				select: function (e) {
					var value = $sectionItemsGrid.find('select.nodes').data('kendoDropDownList').value();
					var valueComp = String(value).split(':');

					if (valueComp.length === 0)
						return;

					registerAjaxPullFor('items', {
						host: valueComp[0],
						port: valueComp[1]
					}, function (res) {
						if (!res.success) {
							return;
						}

						var $grid = $sectionItemsGrid.find('.grid').data('kendoGrid');
						$grid.setDataSource(new kendo.data.DataSource({
							data: res.data.grid,
							pageSize: $grid.dataSource.pageSize()
						}));
					}, function () {
						toastr["error"]("Error occured when fetching items data for selected node")
					});
				}
			});
			$sectionItemsGrid.find('select.nodes').closest('.selector').remove(); */
		};

		// register ajax pull, 
		// make data semi real time
		// interval changeable
		this.registerAjaxPull = function () {

			// prepare ajax pull for nodes,
			// return data which used in both node grid & chart
			registerAjaxPullFor('nodes', {
				seriesLimit: seriesLimit,
				seriesDelay: ajaxPullDelay,
				dataSizeUnit: dataSizeUnit
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

				var $nodeGrid = $sectionNodesGrid.find('.grid').data('kendoGrid');
				var $nodeChart = $sectionNodesChart.find('.chart').data('kendoChart');
				/* var $itemNodeSelect = $sectionItemsGrid.find('select.nodes').data('kendoDropDownList'); */

				$nodeGrid.setDataSource(new kendo.data.DataSource({
					data: Lazy(res.data.grid).map(function (d) {
						d.ConfigHost = (d.ConfigName + ':' + d.ConfigPort);
						return d;
					}).sortBy(function (d) { 
						return -d.StartTime; 
					}).toArray(),
					pageSize: $nodeGrid.dataSource.pageSize()
				}));

				// get max value of each series,
				// then use it as valueAxis.max of each series
				Lazy($nodeChart.options.valueAxis).each(function(v, i) {
					var max = Lazy(res.data.chart).max(function (d) {
						return parseInt(d[v.name], 10)
					})[v.name];

					v.max = max + (5 * String(max).length);
				});

				// sort data using time ascending
				$nodeChart.setDataSource(new kendo.data.DataSource({
					data: Lazy(res.data.chart).sortBy(function (d) { 
						return d.TimeInt; 
					}).toArray()
				}));

				$nodeChart.redraw();

				/** // pupulate nodes data as options in item section
				$itemNodeSelect.setDataSource({
					data: Lazy(res.data.grid).map(function (d) {
						var host = (d.ConfigName + ':' + d.ConfigPort);

						return { 
							value: host, 
							text: (host + ' (' + d.ConfigRole + ')') 
						};
					}).sortBy(function (d) { 
						return d.TimeInt; 
					}).toArray()
				});

				// populate to nodes for the first time
				if (!$('select.nodes').hasClass('first-load')) {
					$('select.nodes').data('kendoDropDownList').options.select();
					$('select.nodes').addClass('first-load');
				}*/
			}, function (a, b, c) {
				toastr["error"]("Error occured when fetching data for nodes")
			}, function (doRequest) {
				setInterval(doRequest, ajaxPullDelay * 1000);
			});

			// prepare ajax pull for items, grid
			registerAjaxPullFor('items', {}, function (res) {
				if (!res.success) {
					return;
				}

				var $itemGrid = $sectionItemsGrid.find('.grid').data('kendoGrid');

				$itemGrid.setDataSource(new kendo.data.DataSource({
					data: Lazy(res.data.grid).map(function (d) {
						return d;
					}).sortBy(function (d) { 
						return -d.LastAccess; 
					}).toArray(),
					pageSize: $itemGrid.dataSource.pageSize()
				}));
			}, null, function (doRequest) {
				setInterval(doRequest, ajaxPullDelay * 1000);
			});
		};

		// register event listener
		this.registerEventListener = function () {

			// when browser resized, do some changes
			$window.on('resize', function () {
				clearTimeout(windowResizeTimeout);

				windowResizeTimeout = setTimeout(function () {

					// redraw chart
					$sectionNodesChart.find('.chart').data('kendoChart').redraw();
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


var rez = {
	success: true,
	data: {
		grid: [
			{ Key: 'name', Value: 'noval', Created: '12-12-2015', Expiry: '6m 2s' },
			{ Key: 'name', Value: 'agung', Created: '08-11-2015', Expiry: '9m 12s' },
			{ Key: 'name', Value: 'prayogo', Created: '10-10-2015', Expiry: '12m 12s' }
		]
	}
}

