(function () {
	'use strict';

	var Main = function () {
		var $body = $('body');

		this.init = function () {
			// toastr init
			toastr.options.closeButton = true
			toastr.options.positionClass = 'toast-bottom-right';

			if ($body.find('[data-page=login]').size() > 0)
				$body.find('.menu-nav').remove()
			
			$body.find('.menu-nav .page-' + $body.find('[data-page]').attr('data-page')).addClass('active');
		}

		this.registerEventListener = function () {
			$body.find('.page-header .pull-left').on('click', function () {
				location.href = '/';
			});
		}
	}

	$(function () {
		var main = new Main();
		main.init();
		main.registerEventListener();
	})
}());