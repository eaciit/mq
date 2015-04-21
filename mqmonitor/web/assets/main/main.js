(function () {
	'use strict';

	var Upload = function () { 
		var self = this;

		this.init = function () {
			$("[name=source-path]").val("/Users/novalagung/Desktop/windows/ecfz/TmpStatus/JSON/");
			$("[name=destination-path]").val("/tmp/disturbed/input/TmpStatus");
		};

		this.getParam =  function () {
			return {
				sourcepath: $("[name=source-path]").val(),
				targetpath: $("[name=destination-path]").val()
			};
		};

		this.registerEvent = function () {
			$(".btn-upload").on("click", function () {
				var $btnUpload = $(this);
				var $loader = $(".loader");
				var $form = $(".path-form");

				$btnUpload.hide();
				$loader.show();
				$form.hide();

				$.ajax({
					url: "/upload",
					type: "POST",
					dataType: "JSON",
					data: self.getParam(),
					timeout: 1000 * 60 * 10
				}).done(function (res) {
					if (res) {
						alert("Upload success!");
					}

					$btnUpload.show();
					$loader.hide();
					$form.show();
				}).error(function (a, b, c) {
					console.log(a);
					console.log(b);
					console.log(c);
					alert("Error occured when uploading the data !");

					$btnUpload.show();
					$loader.hide();
					$form.show();
				});
			});
		};
	};

	$(function () {
		var upload = new Upload();

		upload.init();
		upload.registerEvent();
	});
}());