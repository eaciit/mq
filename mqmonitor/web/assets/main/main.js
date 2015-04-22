(function () {
	'use strict';

	var Main = function () { 
		var self = this;
		var $body = $('body');
		var $sectionNodes = $body.find('.section-nodes');
		var ajaxPullInterval = {};

		// register ajax pull, to make grid shows realtime data
		var registerAjaxPullFor = function (what, success, error) {
			var doRequest = function () {
				$.ajax({
					url: '/data/' + what,
					type: 'get',
					dataType: 'json'
				})
				.success(success)
				.error(error);
			};

			ajaxPullInterval['section-' + what] = setInterval(
				doRequest, 
				self.ajaxPullDelay['section-' + what] * 1000
			);

			doRequest();
		};

		// ajax pull delay in second
		// this duration can be changed on the fly
		this.ajaxPullDelay = {
			'section-nodes': 50
		};

		// initiate all components
		this.init = function () {

			// prepare section nodes grid
			$sectionNodes.find('.grid').kendoGrid({
				chartArea: {
					background: "transparent"
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
						{ field: 'DataCount', title: 'Total', width: 90,
							attributes: { style: 'text-align: right;' } },
						{ field: 'DataSize', title: 'Size (in MB)', width: 90, 
							attributes: { style: 'text-align: right;' } }
					] },
					{ title: 'Time', columns: [
						{ field: 'StartTime', title: 'Start Time', width: 100, 
							attributes: { style: 'text-align: center;' } },
						{ field: 'Duration', title: 'Duration', width: 100, 
							attributes: { style: 'text-align: right;' } }
					] }
				]
			});
		};

		this.registerAjaxPull = function () {
			registerAjaxPullFor('nodes', function (res) {
				var $grid = $sectionNodes.find('.grid').data('kendoGrid');
				$grid.setDataSource(new kendo.data.DataSource({
					data: res.data,
					pageSize: $grid.dataSource.pageSize()
				}));
			}, function (a, b, c) {
				console.log(a, b, c);
				alert('Error occured when fetching data for nodes');
			});
		};
	};

	$(function () {
		// start the magic
		var main = new Main();
		main.init();
		main.registerAjaxPull();
	});
}());





var rez = {
    "data": [
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }, 
    { "Config": { "Name": "127.0.0.1", "Port": 7890, "Role": "Master" }, "DataCount": 0, "DataSize": 0 }
    ],
    "message": "",
    "success": true
};