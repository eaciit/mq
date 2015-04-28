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

				$body.find('.input-key').focus();
			});

			$body.find('.input-key').on('keyup', function (e) {
				if (e.keyCode !== 13)
					return;
				
				var $btnActive = $sectionConsole.find('.btn-set-get .btn.active');
				var isModeSet = $btnActive.hasClass('btn-set');

				if (isModeSet) {
					$body.find('.input-value').focus();
					return;
				}

				$sectionConsole.find('.btn-run').trigger('click');
			});

			$body.find('.input-value').on('keyup', function (e) {
				if (e.keyCode !== 13)
					return;

				$sectionConsole.find('.btn-run').trigger('click');
			});

			$sectionConsole.find('.btn-run').on('click', function () {
				var $content = $sectionConsole.find('.content');
				var $loader = $sectionConsole.find('.loader');
				var $btnActive = $sectionConsole.find('.btn-set-get .btn.active');
				var isModeSet = $btnActive.hasClass('btn-set');
				var param = {
					key: $sectionConsole.find('.input-key').val(),
					value: $sectionConsole.find('.input-value').val(),
				};

				if (isModeSet) {
					param.mode = 'set';
				} else {
					param.mode = 'get';
				}

				$content.html('');
				$content.hide();
				$loader.show();

				$.ajax({
					url: '/console',
					data: param,
					dataType: 'json',
					type: 'post'
				})
				.success(function (res) {
					if (!res.success) {
						$content.show();
						$loader.hide();

						toastr.error(res.message);
						return;
					}

					$content.show();
					$loader.hide();

					if (isModeSet) {
						toastr.success('set value success');
					} else {
						$content.html(res.data);
					}
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