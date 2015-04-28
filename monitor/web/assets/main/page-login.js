(function () {
	'use strict';

	var Login = function () { 
		var self = this;
		var $body = $('body');
		var $sectionLogin = $body.find('.section-login');
		var $window = $(window);

		this.init = function () {

		};

		this.registerEventListener = function () {
			$sectionLogin.find('.btn-reset').on('click', function () {
				$(this).closest('.form').find('input').each(function (i, e) {
					$(e).val('');
				});
			});

			$sectionLogin.find('.btn-login').on('click', function () {
				$.ajax({
					url: '/login',
					type: 'post',
					data: $sectionLogin.find('form').serialize(),
					dataType: 'json'
				})
				.success(function (res) {
					if (!res.success) {
						toastr.error(res.message);
						return;
					}

					$sectionUser.find('.btn-search').trigger('click');
					toastr.success('login success');

					setTimeout(function () {
						document.location.href = '/';
					}, 1 * 1000);
				})
				.error(function (a, b, c) {
					toastr.error('error when trying to login');
				});
			});
		};
	};

	$(function () {
		var login = new Login();
		login.init();
		login.registerEventListener();

		setTimeout(function () {
			toastr.info("Welcome to MQ Monitor.<br />Please login.");
		}, 1000 * 1);
	});
}());