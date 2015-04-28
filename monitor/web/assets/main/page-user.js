(function () {
	'use strict';

	var User = function () { 
		var self = this;
		var $body = $('body');
		var $sectionUser = $body.find('.section-user');
		var $sectionInsert = $body.find('.section-insert');
		var $window = $(window);
		var userEdit = false;

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
					{ field: 'Role', title: 'Role' },
					{ title: 'Options', width: 130, 
						template: '<button class="btn btn-xs btn-primary btn-row-edit"><i class="fa fa-edit"></i>&nbsp;edit</button>&nbsp;<button class="btn btn-xs btn-danger btn-row-delete"><i class="fa fa-remove"></i>&nbsp;delete</button>',
						attributes: { style: 'text-align: center' }
				 	}
				]
			});

			$sectionInsert.find('[name=role]').kendoDropDownList({
				dataSource: {
					data: ['admin']
				},
				optionLabel: 'Select role'
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

			$sectionUser.find('.k-grid').on('click', '.btn-row-edit', function () {
				var $form = $sectionInsert.find('form');
				var uid = $(this).closest('tr[data-uid]').attr('data-uid');
				var data = $sectionUser.find('.k-grid').data('kendoGrid').dataSource.data();
				var rowData = Lazy(data).find(function (d) { return d.uid === uid; });
				rowData.Role = 'admin';

				$sectionUser.find('.btn-add').trigger('click');

				$form.find('[name=username]').prop('disabled', !false);
				$form.find('[name=role]').data('kendoDropDownList').enable(false);

				userEdit = rowData.UserName;
				$form.find('[name=username]').val(rowData.UserName);
				$form.find('[name=role]').data('kendoDropDownList').value(rowData.Role);
			});

			$sectionUser.find('.k-grid').on('click', '.btn-row-delete', function () {
				var uid = $(this).closest('tr[data-uid]').attr('data-uid');
				var data = $sectionUser.find('.k-grid').data('kendoGrid').dataSource.data();
				var rowData = Lazy(data).find(function (d) { return d.uid === uid; });

				if (!confirm('Are you sure want to delete user ' + rowData.UserName + ' ?'))
					return;

				$.ajax({
					url: '/data/users?' + $.param({ username: rowData.UserName }),
					type: 'delete',
					dataType: 'json'
				})
				.success(function (res) {
					if (!res.success) {
						toastr.error(res.message);
						return;
					}

					$sectionUser.find('.btn-search').trigger('click');
					toastr.success('user ' + rowData.UserName + ' successfully deleted');
				})
				.error(function (a, b, c) {
					toastr.error('error when deleting user ' + rowData.UserName);
				});

				console.log('delete');
				console.log(rowData);
			});

			$sectionUser.find('.btn-add').on('click', function () {
				var $form = $sectionInsert.find('form');
				$form.find('[name=username]').prop('disabled', !true);
				$form.find('[name=role]').data('kendoDropDownList').enable(true);

				userEdit = false;
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
					if (userEdit !== false)
						if ($(e).attr('type') !== 'password')
							return;

					$(e).val('');
				});

				if (userEdit === false)
					$sectionInsert.find('[name=role]').data('kendoDropDownList').value('');
			});

			$sectionInsert.find('.btn-save').on('click', function () {
				var $loader = $(this).closest('.nav-menu').next().find('.loader');
				var $form = $(this).closest('.nav-menu').next().find('form');
				var username = $form.find('[name=username]').val();
				var password = $form.find('[name=password]').val();
				var passwordConfirmation = $form.find('[name=password-confirmation]').val();
				var role = $form.find('[name=role]').data('kendoDropDownList').value();

				var isValid = (function (inputs) {
					for (var input in inputs) {
						if (inputs.hasOwnProperty(input)) {
							if (inputs[input].length === 0) {
								toastr.error(input + ' cannot be empty');
								return false;
							} else if (inputs[input].length < 3) {
								if (input === 'role')
									continue;
								
								toastr.error(input + ' need minimum 3 character');
								return false;
							}
						}
					}

					return true;
				}({
					username: username,
					password: password,
					'password confirmation': passwordConfirmation,
					role: role
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
					data: { 
						username: username, 
						password: password, 
						role: role,
						oldusername: (userEdit !== false ? userEdit : ''),
						edit: (userEdit !== false)
					},
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

						if (userEdit !== false)
							toastr.success('changes saved!');
						else
							toastr.success('user ' + username + ' saved!');

						$sectionInsert.find('.btn-back').trigger('click');
						$sectionUser.find('.btn-search').trigger('click');

						userEdit = false;
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
