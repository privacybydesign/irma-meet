<!DOCTYPE html>

<html lang="en">

 <head>

   @include('layout.partials.head')

 </head>

 <body>


<div class="container-fluid h-100 d-flex flex-column">

    <div class="row">
        <main class="col-md-4 bg-blue" >
        	<div class="title">
        		IRMA-meet.nl
	        </div>
        </main>
        <main class="col-md-4 bg-warmwhite">
        </main>
        <main class="col-md-4 bg-warmwhite">
        </main>
    </div>
    <div class="row flex-fill">
        <main class="col-md-4 bg-blue" >
            <div class="title">
            </div>
        </main>
        <main class="col-md-4 nopadding-md bg-warmwhite text-center">
            <div class="card">
                <div class="card-header">
                    Welcome
                </div>              
                <div class="card-body">
                   <img class="img-fluid" src="{{ url('/') }}/img/IRMA-meet.png" alt="IRMA-meet logo">
                </div>
            </div>
            
        </main>
        <main class="col-md-4 bg-warmwhite">
        	<div class="content">
        		<p class="cl-grey">AUTHENTICATED VIDEO CALLS</p>
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

    <div class="row img-warp">
        <main class="col-md-7 bg-grey" >
            <div class="info">
                <h2>Know who you are talking to</h2>
                <p>Participants in IRMA-meet video calls need to prove who they are by disclosing relevant information (like their name and e-mail address) about themselves before they can enter a video chat.</p>
            </div>
        </main>
        <main class="col-md-5 bg-grey" >
            <div class="img-warp">
                <img class="img-fluid know" src="{{ url('/') }}/img/know.svg" />
            </div>
        </main>
    </div>
    <div class="row img-warp">
        <div class="row img-warp">
            <main class="col-md-5 bg-light-grey" >
                <div class="col-sm-5"><img class="img-fluid prove" src="{{ url('/') }}/img/prove.svg" /></div>
                <div class="col-sm-7"><img class="img-fluid irma_logo_info" src="{{ url('/') }}/img/irma_logo.svg" /></div>
            </main>
        </div>
        <main class="col-md-7 bg-light-grey" >
            <div class="info">
                <h2>Prove who you are with IRMA</h2>
                <p>You can prove who you are via <a href="https://privacybydesign.foundation">IRMA</a>. It is a freely available open source identity app; it is a wallet on your phone for collecting personal data. This data is stored only on your phone, and nowhere else. You control where you disclose it.</p>
            </div>
        </main>
    </div>
    <div class="row">
        <main class="col-md-7" >
            <div class="info">
                <div>
                    <h2>Get ready</h2>
                </div>
                <p>All IRMA-meet users need to have the IRMA app on their phone, filled with the relevant personal data. Which data is required may depend on the video call.</p>
                <p>If you don't have IRMA yet, install it from the IOS appstore or Google Play store and load the required personal data, like your name or e-mail address.</p>
                <p>At this early stage, there is only one type of free meeting available, which requires your e-mail address and your name in IRMA. Your name can be obtained either from your LinkedIn profile or for users in the Netherlands from the Dutch citizen registry.</p>
            </div>
        </main>
        <main class="col-md-5" >
            <div class="">
                <img class="img-fluid " src="{{ url('/') }}/img/getready.svg" />
            </div>
        </main>
    </div>
    

</div>

@include('layout.partials.footer-scripts')

 </body>

</html>
