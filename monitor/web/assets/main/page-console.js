(function () {
	'use strict';

	var Console = function () { 
		var self = this;
		var $body = $('body');
		var $sectionConsole = $body.find('.section-console');
		var $window = $(window);

		this.init = function () {

		};

		this.registerEventListener = function () {
			$sectionConsole.find('.btn-set-get .btn').on('click', function () {
				$(this).siblings().removeClass('active');
				$(this).addClass('active');

				if ($(this).hasClass('btn-set')) {
					$sectionConsole.find('.inline-set').show();
				} else if ($(this).hasClass('btn-get')) {
					$sectionConsole.find('.inline-set').hide();
				}
			});

			$sectionConsole.find('.btn-run').on('click', function () {
				var $content = $sectionConsole.find('.content');
				var $loader = $sectionConsole.find('.loader');
				var $btnActive = $sectionConsole.find('.btn-set-get .btn.active');
				var param = {
					key: $sectionConsole.find('.input-key').val(),
					value: $sectionConsole.find('.input-value').val(),
				};

				if ($btnActive.hasClass('btn-set')) {
					param.mode = 'set';
				} else if ($btnActive.hasClass('btn-get')) {
					param.mode = 'get';
				}

				$content.hide();
				$loader.show();

				$.ajax({
					url: '/console',
					data: param,
					dataType: 'json',
					type: 'get'
				})
				.success(function (res) {
					if (!res.success) {
						$content.show();
						$loader.hide();

						toastr.error(res.message);
						return;
					}

					setTimeout(function () {
						$content.show();
						$loader.hide();

						console.log(res);
					}, 1 * 1000);
				})
				.error(function () {
					$content.show();
					$loader.hide();
					toastr.error('error when trying to login');
				});
			});
		};
	};

	$(function () {
		var console = new Console();
		console.init();
		console.registerEventListener();

		$('.btn-set-get .btn-get').trigger('click');
	});
}());