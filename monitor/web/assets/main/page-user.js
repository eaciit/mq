(function () {
	'use strict';

	var User = function () { 
		var self = this;
		var $body = $('body');
		var $sectionUser = $body.find('.section-user');
		var $sectionInsert = $body.find('.section-insert');
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
				]
			});

			$sectionUser.show();
			$sectionInsert.hide();
			$sectionInsert.find('.loader').hide();
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

			$sectionUser.find('.btn-add').on('click', function () {
				$sectionInsert.find('.btn-reset').trigger('click');
				$sectionUser.hide();
				$sectionInsert.show();
			});

			$sectionInsert.find('.btn-back').on('click', function () {
				$sectionInsert.hide();
				$sectionUser.show();
			});

			$sectionInsert.find('.btn-reset').on('click', function () {
				$(this).closest('.nav-menu').next().find('input').each(function (i, e) {
					$(e).val('');
				})
			});

			$sectionInsert.find('.btn-save').on('click', function () {
				var $loader = $(this).closest('.nav-menu').next().find('.loader');
				var $form = $(this).closest('.nav-menu').next().find('form');
				var username = $form.find('[name=username]').val();
				var password = $form.find('[name=password]').val();
				var passwordConfirmation = $form.find('[name=password-confirmation]').val();

				var isValid = (function (inputs) {
					for (var input in inputs) {
						if (inputs.hasOwnProperty(input)) {
							if (inputs[input].length === 0) {
								toastr.error(input + ' cannot be empty');
								return false;
							} else if (inputs[input].length < 7) {
								toastr.error(input + ' need minimum 6 character');
								return false;
							}
						}
					}

					return true;
				}({
					username: username,
					password: password,
					'password confirmation': passwordConfirmation
				}));

				if (!isValid)
					return;

				if (password !== passwordConfirmation) {
					$form.find('[name=password]').val('');
					$form.find('[name=password-confirmation]').val('');
					toastr.error('password do not match');
					return;
				}

				$loader.show();
				$form.hide();

				$.ajax({
					url: 'data/users',
					data: { username: username, password: password },
					type: 'post',
					dataType: 'json'
				})
				.success(function (res) {
					setTimeout(function () {
						$loader.hide();
						$form.show();
						
						if (!res.success) {
							toastr.error(res.message);
							$form.find('[name=password]').val('');
							$form.find('[name=password-confirmation]').val('');
							return;
						}

						toastr.success('user ' + username + ' saved!');
						$sectionInsert.find('.btn-back').trigger('click');
						$(this).closest('.nav-search').find('.btn-search').trigger('click');
					}, 500);
				})
				.error(function (a, b, c) {
					setTimeout(function () {
						$loader.hide();
						$form.show();
						toastr.error('error occured when saving user');
					}, 500);
				});
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
