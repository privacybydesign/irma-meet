<!DOCTYPE html>

<html lang="en">

 <head>

   @include('layout.partials.head')

 </head>

 <body>


<div class="container-fluid h-100 d-flex flex-column">

    <div class="row">
        <main class="col-4 bg-blue" >
        	<div class="title">
        		IRMA-meet.nl
	        </div>
        </main>
        <main class="col-4 bg-warmwhite">
        </main>
        <main class="col-4 bg-warmwhite">
        </main>
    </div>
    <div class="row flex-fill">
        <main class="col-4 bg-blue" >
        	<div class="title">
	        </div>
        </main>
        <main class="col-4 bg-warmwhite text-center">
        	<div class="">
    	    	<button type="button" class="btn btn-primary btn-lg btn-blue">Start a meeting now!</button>
	        </div>
        </main>
        <main class="col-4 bg-warmwhite">
        	<div class="content">
        		<p class="cl-grey">THE SECURE VIDEO CHAT</p>
        		<h1>Meet online, safe and sound!</h1>
        		<p>Invite your guests and IRMA-meet will do the rest. No surprise visitors guaranteed</p>
	        </div>
        </main>
    </div>
    <div class="row footer">
        <main class="col-4 bg-blue" >
        	<div class="title">
	        </div>
        </main>
        <main class="col-4 bg-warmwhite">
        </main>
        <main class="col-4 bg-warmwhite">
        	<p class="cl-grey">FAQ</p>
        </main>
    </div>
</div>

@include('layout.partials.footer-scripts')

 </body>

</html>