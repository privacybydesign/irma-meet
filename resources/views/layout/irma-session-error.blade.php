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
            <a href="{{ $mainContent }}">{{ $mainContent }}</a>
        </main>
        <main class="col-4 bg-warmwhite">
        	<div class="content">
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