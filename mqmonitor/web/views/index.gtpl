<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />

		<title>MQ Monitor</title>

		<!-- res/jquery -->
		<script src="/res/jquery/jquery-2.1.3.min.js"></script>

		<!-- res/bootstrap -->
		<link rel="stylesheet" href="/res/bootstrap/css/bootstrap.min.css">
		<link rel="stylesheet" href="/res/bootstrap/css/bootstrap-theme.min.css">
		<script src="/res/bootstrap/js/bootstrap.min.js"></script>

		<!-- res/kendoui -->
		<link rel="stylesheet" href="/res/kendoui/styles/kendo.common.min.css" />
		<link rel="stylesheet" href="/res/kendoui/styles/kendo.silver.min.css" />
		<link rel="stylesheet" href="/res/kendoui/styles/kendo.dataviz.min.css" />
		<link rel="stylesheet" href="/res/kendoui/styles/kendo.dataviz.silver.min.css" />
		<script src="/res/kendoui/js/kendo.all.min.js"></script>

		<!-- res/toastr -->
		<link rel="stylesheet" href="/res/toastr/toastr.min.css">
		<script src="/res/toastr/toastr.min.js"></script>

		<!-- css/font-awesome -->
		<link rel="stylesheet" href="/res/font-awesome/css/font-awesome.min.css">

		<!-- res/lazy -->
		<script src="/res/lazy/lazy.min.js"></script>

		<!-- res/main -->
		<link rel="stylesheet" href="/res/main/main.css">
		<script src="/res/main/main.js"></script>
	</head>

	<body>
		<div class="container-fluid main no-padding">
			<div class="page-header">
				<div class="pull-left">
					<img class="logo" src="/res/images/logo.png" />
				</div>
				<div class="pull-left">
					<h1>MQ Monitor</h1>
				</div>
				<div class="clearfix"></div>
			</div>

			<div class="col-md-12">
				<div class="col-md-6 section section-nodes-grid">
					<div class="panel panel-primary">
						<div class="panel-heading">
							<i class="fa fa-cogs"></i> Nodes Information
						</div>
						<div class="panel-body">
							<div class="col-md-12 nav-search">
								<div class="input-group input-sm">
									<div class="input-group-addon input-sm">Search</div>
									<input type="text" class="form-control input-sm input-search" placeholder="Type search keyword here ..." />
									<button class="btn btn-sm btn-success btn-search" data-ajax="nodes">
										<span class="glyphicon glyphicon-search"></span> Search
									</button>
								</div>
							</div>
							<div class="row no-padding no-margin">
								<div class="grid"></div>
							</div>
						</div>
					</div>
				</div>

				<div class="col-md-6 section section-nodes-chart">
					<div class="panel panel-primary">
						<div class="panel-heading">
							<i class="fa fa-bar-chart-o"></i> Nodes Information
						</div>
						<div class="panel-body">
							<div class="chart"></div>
						</div>
					</div>
				</div>

				<div class="col-md-6 section section-items-grid">
					<div class="panel panel-primary">
						<div class="panel-heading">
							<i class="fa fa-files-o"></i> Items Information
						</div>
						<div class="panel-body">
							<div class="col-md-12 nav-search">
								<div class="input-group input-sm">
									<div class="input-group-addon input-sm">Search</div>
									<input type="text" class="form-control input-sm input-search" placeholder="Type search keyword here ..." />
									<button class="btn btn-sm btn-success btn-search" data-ajax="items">
										<span class="glyphicon glyphicon-search"></span> Search
									</button>
								</div>
							</div>
							<div class="row no-padding no-margin">
								<div class="grid"></div>
							</div>
						</div>
					</div>
				</div>

				<div class="clearfix"></div>
			</div>

			<div class="loader">
				<img src="/res/images/495.gif" />
			</div>
		</div>
	</body>
</html>
