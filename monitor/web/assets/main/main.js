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
		var ajaxPullFor = {};
		var timeoutFor = {};

		// register ajax pull, to make grid shows realtime data
		var registerAjaxPullFor = function (what, data, success, error, after) {
			var doProcess = ajaxPullFor[what] = function () {
				success = (typeof success !== String(undefined)) ? success : function () {};
				error 	= (typeof error   !== String(undefined)) ? error   : function () {};

				console.log($.extend(true, { 
						isServerAlive: isServerAlive
					}, data()));

				$.ajax({
					url: '/data/' + what,
					data: $.extend(true, { 
						isServerAlive: isServerAlive
					}, data()),
					type: 'get',
					dataType: 'json'
				})
				.success(function (res) {
					if (typeof success !== String(undefined))
						success(res);
					if (typeof after !== String(undefined))
						timeoutFor[what] = after(doProcess);
				})
				.error(function (a, b, c) {
					if (typeof error !== String(undefined))
						error(a, b, c);
					if (typeof after !== String(undefined))
						timeoutFor[what] = after(doProcess);
				});
			};

			doProcess();
		};

		// initiate all components
		this.init = function () {
			// toastr init
			toastr.options.closeButton = true
			toastr.options.positionClass = 'toast-bottom-right';

			// notify about delay every some minutes
			setInterval(function () {
				toastr.error('data refreshed every ' + ajaxPullDelay + ' seconds');
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
					{ field: 'Created', title: 'Created', width: 140,
						attributes: { style: 'text-align: center;' } },
					{ field: 'Expiry', title: 'Expiry', width: 80,
						attributes: { style: 'text-align: center;' } },
				]
			});

			$body.find('.menu-nav .page-' + $body.find("[data-page]").attr("data-page")).addClass('active');
		};

		// register ajax pull, 
		// make data semi real time
		// interval changeable
		this.registerAjaxPull = function () {

			// prepare ajax pull for nodes,
			// return data which used in both node grid & chart
			registerAjaxPullFor('nodes', function () {
				return {
					seriesLimit: seriesLimit,
					seriesDelay: ajaxPullDelay,
					dataSizeUnit: dataSizeUnit,
					search: $sectionNodesGrid.find('.nav-search .input-search').val()
				};
			}, function (res) {
				if (!res.success) {
					if (res.message === 'connection is shut down')
						isServerAlive = false;

					toastr.error(res.message);
					return;
				}

				if (isServerAlive == false) {
					isServerAlive = true;
					toastr.success("connected to server");
				}

				var $nodeGrid = $sectionNodesGrid.find('.grid').data('kendoGrid');
				var $nodeChart = $sectionNodesChart.find('.chart').data('kendoChart');

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
			}, function (a, b, c) {
				toastr.error("Error occured when fetching data for nodes")
			}, function (doProcess) {
				return setTimeout(doProcess, ajaxPullDelay * 1000);
			});

			// prepare ajax pull for items, grid
			registerAjaxPullFor('items', function () { 
				return {
					search: $sectionItemsGrid.find('.nav-search .input-search').val()
				};
			}, function (res) {
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
			}, null, function (doProcess) {
				return setTimeout(doProcess, ajaxPullDelay * 1000);
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

			$body.find('.btn-search').on('click', function () {
				var what = $(this).attr('data-ajax');

				clearTimeout(timeoutFor[what]);
				ajaxPullFor[what]();
			});

			$body.find('.input-search').on('keyup', function (e) {
				if (e.keyCode !== 13)
					return;

				$(this).closest('.nav-search').find('.btn-search').trigger('click');
			});
		};
	};

	// start the magic
	$(function () {
		var main = new Main();
		main.init();
		main.registerAjaxPull();
		main.registerEventListener();

		setTimeout(function () {
			toastr.info("Welcome to MQ Monitor");
		}, 1000 * 1);
	});
}());
