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
            <img class="img-fluid" src="{{ url('/') }}/img/IRMA-meet.png" alt="IRMA-meet logo">
        </main>
        <main class="col-4 bg-warmwhite">
        	<div class="content">
        		<p class="cl-grey">AUTHENTICATED VIDEO CHATTING</p>
        		<h1>Meet online, without surprises!</h1>
        		<p>Only people who prove who they are can participate in IRMA-meet video. Chat with certainty, without being fooled, in your business meeting, (medical) consult, or online oral exam.</p>
                <div class="">
                    <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('./irma_auth/start', './irma_session/create');">Start a meeting now!</a>
                </div>
                <div id="errorbox" class="alert alert-danger alert-dismissible fade">
                    <span id="errorboxtext"></span>
                    <button type="button" class="close" data-dismiss="alert">&times;</button>
                </div>
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
