(function () {
	'use strict';

	var User = function () { 
		var self = this;
		var $body = $('body');
		var $sectionUser = $body.find('.section-user');
		var $window = $(window);

		this.init = function () {
			$sectionUser.find('.grid').kendoGrid({
				dataSource: { 
					data: [], 
					pageSize: 10
				},
				pageable: {
					pageSizes: [5, 10, 15, 20]
				},
				sortable: true, 
				scrollable: false,
				columns: [
					{ field: 'UserName', title: 'User Name' },
					{ title: 'Options', width: 130, 
						template: '<button class="btn btn-xs btn-primary row-edit"><i class="fa fa-edit"></i>&nbsp;edit</button>&nbsp;<button class="btn btn-xs btn-danger row-delete"><i class="fa fa-remove"></i>&nbsp;delete</button>',
						attributes: { style: 'text-align: center' }
				 	}
				],
			});
		}

		// register event listener
		this.registerEventListener = function () {
			$body.find('.btn-search').on('click', function () {
				var what = $(this).attr('data-ajax');

				$.ajax({
					url: '/data/users',
					data: {
						search: $sectionUser.find('.nav-search .input-search').val()
					},
					type: 'get',
					dataType: 'json'
				})
				.success(function (res) {
					if (!res.success) {
						toastr.error(res.message);
						return;
					}

					var $userGrid = $sectionUser.find('.grid').data('kendoGrid');

					$userGrid.setDataSource(new kendo.data.DataSource({
						data: Lazy(res.data.grid).sortBy(function (d) { 
							return d.UserName; 
						}).toArray(),
						pageSize: $userGrid.dataSource.pageSize()
					}));
				})
				.error(function (a, b, c) {
					toastr.error('error occured when fetching user data');
				});
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
		var user = new User();
		user.init();
		user.registerEventListener();

		$('.btn-search').trigger('click');

		setTimeout(function () {
			toastr.info("Welcome to User Management");
		}, 1000 * 1);
	});
}());
