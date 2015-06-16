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
				if ($(this).hasClass('btn-set') && !$(this).hasClass('active')) {
					$sectionConsole.find('.input-value').val('');
				}

				$("button").siblings().removeClass('active');
				$(this).addClass('active');

				if ($(this).hasClass('btn-set')) {
					$sectionConsole.removeClass('mode-get');
				} else if ($(this).hasClass('btn-get')) {
					$sectionConsole.addClass('mode-get');
				}

				$body.find('.input-key').focus();
				$(window).trigger('resize');
			});

			$sectionConsole.find('.btn-keys-info .btn').on('click', function () {
				if ($(this).hasClass('btn-keys') && !$(this).hasClass('active')) {
					$sectionConsole.find('.input-value').val('');
				}

				$("button").siblings().removeClass('active');
				$(this).addClass('active');

				$sectionConsole.addClass('mode-get');

				$body.find('.input-key').focus();
				$(window).trigger('resize');
			});

			$sectionConsole.find('.btn-read-write .btn').on('click', function () {
				if ($(this).hasClass('btn-read') && !$(this).hasClass('active')) {
					$sectionConsole.find('.input-value').val('');
				}

				$("button").siblings().removeClass('active');
				$(this).addClass('active');

				$sectionConsole.addClass('mode-get');

				$body.find('.input-key').focus();
				$(window).trigger('resize');
			});

			$sectionConsole.find('.input-key').on('keyup', function (e) {
				if (e.keyCode !== 13)
					return;

				var $btnActive = $sectionConsole.find('.nav-button .btn.active');
				var isModeSet = $btnActive.hasClass('btn-set');

				if (isModeSet) {
					$body.find('.input-value').focus();
					return;
				}

				$sectionConsole.find('.btn-run').trigger('click');
			});

			$sectionConsole.find('.input-value').on('keyup', function (e) {
				if (e.keyCode !== 13)
					return;

				$sectionConsole.find('.btn-run').trigger('click');
			});

			$sectionConsole.find('.btn-run').on('click', function () {

				var $content = $sectionConsole.find('.content');
				var $loader = $sectionConsole.find('.loader');
				var $btnActive = $sectionConsole.find('.nav-button .btn.active');
				var mode;

				if($btnActive.hasClass('btn-set')){
					mode = 'set';
				}else if ($btnActive.hasClass('btn-keys')) {
					mode = 'keys';
				}else if ($btnActive.hasClass('btn-infoo')) {
					mode = 'info';
				}else if ($btnActive.hasClass('btn-write')) {
					mode = 'write';
				}else if ($btnActive.hasClass('btn-read')) {
					mode = 'read';
				}else{
					mode = 'get';
				}

				var param = {
					mode: mode,
					key: $.trim($sectionConsole.find('.input-key').val()),
					value: $.trim($sectionConsole.find('.input-value').val()),
					owner: $.trim($sectionConsole.find('.input-owner').val()),
					table: $.trim($sectionConsole.find('.input-table').val()),
					duration: $.trim($sectionConsole.find('.input-duration').val()),
					permission: $.trim($sectionConsole.find('.input-permission').val()),
				};

				if (param.key.length === 0) {
					toastr.error('key cannot be empty');
					return;
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

					if (mode == "set") {
						toastr.success('set value success');
					}else {
						$content.html(res.data);
					}
				})
				.error(function () {
					$content.show();
					$loader.hide();
					toastr.error('error when trying to login');
				});
			});

			$sectionConsole.find('.btn-detail').on('click', function () {
				var isActive = $(this).hasClass('active');

				if (isActive) {
					$(this).removeClass('active');
					$('.nav-button:eq(0)').addClass('bordered');
					$('.nav-button:eq(1)').hide();
				} else {
					$(this).addClass('active');
					$('.nav-button:eq(0)').removeClass('bordered');
					$('.nav-button:eq(1)').show();
				}
			});

			$window.on('resize', function () {
				var widthOfNav = $('.nav-button:eq(0)').width();
				var widthOfSetGet = $('.nav-button:eq(0) .inline:eq(0)').width();
				var widthOfInputKey = 200;
				var widthOfInputValue = $('.nav-button:eq(0) .inline:eq(2)').width();
				var widthOfButtonRun = $('.nav-button:eq(0) .inline:eq(3)').width();
				var widthWithoutSetGetRun = (widthOfNav - widthOfSetGet - widthOfButtonRun);
				var widthDetailEach = widthWithoutSetGetRun / 4;

				if ($sectionConsole.hasClass('mode-get')) {
					$('.nav-button:eq(0) .inline:eq(1)').width(widthWithoutSetGetRun);
				} else {
					$('.nav-button:eq(0) .inline:eq(1)').width(widthOfInputKey);
					$('.nav-button:eq(0) .inline:eq(2)').width(widthWithoutSetGetRun - widthOfInputKey);
				}

				$('.nav-button:eq(1) .inline:eq(0)').width($('.nav-button:eq(0) .inline:eq(0)').width());
				$('.nav-button:eq(1) .inline:eq(1)').width(widthDetailEach);
				$('.nav-button:eq(1) .inline:eq(2)').width(widthDetailEach);
				$('.nav-button:eq(1) .inline:eq(3)').width(widthDetailEach);
				$('.nav-button:eq(1) .inline:eq(4)').width(widthDetailEach);
			});
		};
	};

	$(function () {
		var console = new Console();
		console.init();
		console.registerEventListener();

		$('.btn-set-get .btn-get').trigger('click');
		$(window).trigger('resize');
	});
}());
