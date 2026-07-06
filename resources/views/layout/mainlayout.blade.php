<!DOCTYPE html>

<html lang="en">

<head>

    @include('layout.partials.head')

</head>

<body>

    <div class="text-center py-2" style="background:#f0ad4e; color:#333; font-size:0.9rem;">
        <strong>{{ __('demo_notice_label') }}</strong>
        {{ __('demo_notice_text') }}
    </div>

    <div class="container-fluid h-100 d-flex flex-column">

        <div class="row">
            <main class="col-md-4 bg-blue">
                <div class="title">
                    Yivi Meet
                </div>
            </main>
            <main class="col-md-4 bg-warmwhite">
            </main>
            <main class="col-md-4 bg-warmwhite">
                <ul class="navbar-nav ms-auto">
                    <!-- Authentication Links -->
                    @php $locale = session()->get('locale'); @endphp
                    <li class="nav-item dropdown text-end">
                        <a id="navbarDropdown" class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-haspopup="true" aria-expanded="false" v-pre>
                            @switch($locale)
                            @case('en')
                            <img src="{{asset('img/flags/gbr.svg')}}" width="30px" height="20x"> English
                            @break
                            @default
                            <img src="{{asset('img/flags/nld.svg')}}" width="30px" height="20x"> Nederlands
                            @endswitch
                            <span class="caret"></span>
                        </a>
                        <div class="dropdown-menu dropdown-menu-end" aria-labelledby="navbarDropdown">
                            <a class="dropdown-item" href="{{ url('/') }}/lang/en"><img src="{{asset('img/flags/gbr.svg')}}" width="30px" height="20x"> English</a>
                            <a class="dropdown-item" href="{{ url('/') }}/lang/nl"><img src="{{asset('img/flags/nld.svg')}}" width="30px" height="20x"> Nederlands</a>
                        </div>
                    </li>
                </ul>

            </main>
        </div>
        <div class="row flex-fill">
            <main class="col-md-4 bg-blue">
                <div class="title">
                </div>
            </main>
            <main class="col-md-4 nopadding-md bg-warmwhite text-center">
                <div class="card">
                    <div class="card-header">
                        {!! $title !!}
                    </div>
                    <div class="card-body">
                        {!! $message !!}
                    </div>
                </div>

            </main>
            <main class="col-md-4 bg-warmwhite">
                <div class="content">
                    <p class="cl-grey slogan">{{ __('AUTHENTICATED VIDEO CALLS') }}</p>
                    <h1>{{ __('Meet online, without surprises') }}</h1>
                    <p>{{ __('Only people who prove who they are can participate in IRMA-meet video. Chat with certainty, without being fooled, in your business meeting, (medical) consult, or online oral exam.') }}</p>
                    <br>
                    <div class="">
                        {!! $buttons !!}
                    </div>
                    <div>
                        <div>
                            <br>
                            {!! __('Not sure which meeting type is right for you? <br> Find out more <a href="#threemeetings">below</a>!') !!}
                        </div>
                        <div id="errorbox" class="alert alert-danger alert-dismissible fade">
                            <span id="errorboxtext"></span>
                            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="{{ __('Close') }}"></button>
                        </div>
                    </div>
            </main>
        </div>

        <div class="row-full bg-grey">
            <div class="row img-warp tight">
                <main class="col-md-7 bg-grey">
                    <div class="info">
                        <h2>{{ __('Know who you are talking to') }}</h2>
                        <p>{{ __('Participants in IRMA-meet video calls need to prove who they are by disclosing relevant information (like their name and e-mail address) about themselves before they can enter a video chat.') }}</p>
                    </div>
                </main>
                <main class="col-md-5 bg-grey">
                    <div class="img-warp">
                        <img class="img-fluid know" src="{{ url('/') }}/img/know.svg" />
                    </div>
                </main>
            </div>
        </div>

        <div class="row-full bg-light-grey">
            <div class="row img-warp tight">
                <main class="col-md-5">
                    <div class="">
                        <img class="img-fluid irma_logo_info" src="{{ url('/') }}/img/yivi-logo.svg" />
                    </div>
                </main>
                <main class="col-md-7">
                    <div class="info">
                        <h2>{{ __('Prove who you are with IRMA') }}</h2>
                        <p>{!! __('You can prove who you are via <a href="https://www.irma.app">IRMA</a>. It is a freely available open source identity app; it is a wallet on your phone for collecting personal data. This data is stored only on your phone, and nowhere else. You control where you disclose it.') !!}</p>
                    </div>
                    <div>
                        <img class="img-fluid irma_logo_small" src="{{ url('/') }}/img/yivi-logo.svg" />
                    </div>
                </main>
            </div>
        </div>


        <div class="row img-warp tight">
            <main class="col-md-7">
                <div class="info">
                    <div>
                        <h2>{{ __('Get ready') }}</h2>
                    </div>
                    <p>{{ __('All IRMA-meet users need to have the IRMA app on their phone, filled with the relevant personal data. Which data is required may depend on the video call.') }}</p>
                    <p>{!! __('If you don\'t have IRMA yet, install it from the <a href="https://apps.apple.com/us/app/irma-authentication/id1294092994">iOS App Store</a> or <a href="https://play.google.com/store/apps/details?id=org.irmacard.cardemu&hl=en">Google Play Store</a> and load the required personal data, like your name or e-mail address.') !!}</p>
                    <p>{{ __('At this early stage, there is only one type of free meeting available, which requires your e-mail address and your name in IRMA. Your name can be obtained either from your LinkedIn profile or for users in the Netherlands from the Dutch citizen registry.') }}</p>
                </div>
            </main>
            <main class="col-md-5">
                <div class="">
                    <img class="img-fluid getready" src="{{ url('/') }}/img/getready.svg" />
                </div>
            </main>
        </div>


        <div class="row-full bg-grey" id="threemeetings">
            <div class="info tight" style="padding-bottom: 1rem;">
                <h2>{{ __('Choose among three free options') }}</h2>
            </div>
            <div class="row img-warp tight g-0">
                <main class="col-sm-3">
                    <div class="type">
                        <p class="typeH">{{ __('Video meeting') }}</p><img class="circles" src="{{ url('/') }}/img/team.svg" />
                    </div>
                    <div>

                    </div>
                </main>
                <main class="col-sm-9">
                    <div class="info">
                        <div class="speech-bubble">
                            <p>
                                <b>{{ __('For whom?') }}</b> {{ __('Video meetings with up to six participants can be hosted by everyone, globally.') }}
                            </p>
                            <span id="free" style="display: none;">
                                <p>
                                    <b>{{ __('How does it work?') }}</b> {{ __('The host proves who they are with the IRMA app, sets up the meeting and sends the meeting link to all participants. To enter the meeting, both host and participants must make themselves known by disclosing their email address and name, via the IRMA app.') }}
                                </p>
                                <p>
                                    <b>{{ __('What do you need?') }}</b> {!! __('For video meetings, both the host and the participants need their <a href="https://privacybydesign.foundation/uitgifte/email/" class="light">email</a> and name in their IRMA app. Users who live in the Netherlands can use their official name, from the <a href="https://services.nijmegen.nl/irma/gemeente/start" class="light">national civil registry</a>. International users use their name from <a href="https://privacybydesign.foundation/uitgifte/social/linkedin/" class="light">LinkedIn</a>. The latter is not really reliable because it is self-chosen and can be changed at any moment. The email address does offer more certainty, since it is verified.') !!}
                                </p>
                            </span>
                            <p class="toggle" onclick="toggleMoreInfo('free')" id="freeBtn">{{ __('Show more') }}</p>
                        </div>
                        </p>
                    </div>
                </main>
            </div>
            <div class="row img-warp tight g-0">
                <main class="col-sm-3">
                    <div class="type">
                        <p class="typeH">{{ __('Oral video exam') }}</p><img class="circles" src="{{ url('/') }}/img/exam.svg" />
                    </div>
                    <div>

                    </div>
                </main>
                <main class="col-sm-9">
                    <div class="info">
                        <div class="speech-bubble">
                            <p>
                                <b>{{ __('For whom?') }}</b> {{ __('Academics in the Netherlands can use IRMA-meet for one-on-one oral exams via video, with verified identity of the student.') }}
                            </p>
                            <span id="exam" style="display: none;">
                                <p>
                                    <b>{{ __('How does it work?') }}</b> {{ __('The teacher proves that they are a teacher with the IRMA app, sets up the meeting and sends the meeting link to the student. The student has to prove who (s)he is with the IRMA app before entering the video exam. The teacher then knows for sure that the right student is taking the exam.') }}
                                </p>
                                <p>
                                    <b>{{ __('What do you need?') }}</b> {!! __('For video exams, both the teacher and the student need the academic personal data from <a href="https://privacybydesign.foundation/uitgifte/surfnet/surfnet/" class="light">SURFconext</a> in their IRMA app.') !!}
                                </p>
                            </span>
                            <p class="toggle" onclick="toggleMoreInfo('exam')" id="examBtn">{{ __('Show more') }}</p>
                        </div>
                        </p>
                    </div>
                </main>
            </div>
            {{-- Medical consult section hidden: requires deprecated Yivi credentials (AGB-code/BSN) --}}
        </div>

        <script>
            function toggleMoreInfo(meetingType) {
                var moreText = document.getElementById(meetingType);
                var toggleText = document.getElementById(meetingType + "Btn");

                if (toggleText.innerHTML == "{{ __('Show less') }}") {
                    toggleText.innerHTML = "{{ __('Show more') }}";
                    moreText.style.display = "none";
                } else {
                    toggleText.innerHTML = "{{ __('Show less') }}";
                    moreText.style.display = "inline";
                }
            }
        </script>


    </div>
    </div>

    <div class="row-full bg-light-grey">
        <div class="info">
            <h2 style="text-align: center;">{{ __('How it works') }}</h2>
        </div>
        <div class="row img-warp tight">
            <main class="col-lg-4">
                <div class="step">
                    <h3>{{ __('Step 1') }}</h3>
                    <h5>{{ __('Host proves who (s)he is') }}</h5>
                    <p class="explanation">{!! __('The host initiates a new meeting on <a href="#">this website</a> and makes herself/himself known, via the IRMA app.') !!}</p>
                    <img class="step-img" src="{{ url('/') }}/img/yivi-logo.svg" />
                </div>
            </main>
            <main class="col-lg-4">
                <div class="step">
                    <h3>{{ __('Step 2') }}</h3>
                    <h5>{{ __('Host chooses meeting name') }}</h5>
                    <p class="explanation">{{ __('The host chooses a name for the meeting and obtains a link for the meeting.') }}</p>
                    <img class="step-img" src="{{ url('/') }}/img/step2.svg" />

                </div>
            </main>
            <main class="col-lg-4">
                <div class="step">
                    <h3>{{ __('Step 3') }}</h3>
                    <h5>{{ __('Host distributes meeting link') }}</h5>
                    <p class="explanation">{{ __('The host distributes the meeting link to the intended participants, for instance via email or direct messaging.') }}</p>
                    <img class="step-img" src="{{ url('/') }}/img/step3.svg" />

                </div>
            </main>
        </div>
        <div class="row img-warp tight">
            <main class="col-lg-4">
                <div class="step">
                    <h3>{{ __('Step 4') }}</h3>
                    <h5>{{ __('Participants prove who they are') }}</h5>
                    <p class="explanation">{{ __('When a participant clicks on the link (s)he is asked to prove who (s)he via IRMA.') }}</p>
                    <img class="step-img" src="{{ url('/') }}/img/yivi-logo.svg" />

                </div>
            </main>
            <main class="col-lg-4">
                <div class="step">
                    <h3>{{ __('Step 5') }}</h3>
                    <h5>{{ __('Participants enter meeting') }}</h5>
                    <p class="explanation">{{ __('Each participant joins the meeting automatically. The verified identity of everyone in the meeting is visible.') }}</p>
                    <img class="step-img" src="{{ url('/') }}/img/step5.svg" />

                </div>
            </main>
            <main class="col-lg-4">
                <div class="step">
                    <h3>{{ __('Step 6') }}</h3>
                    <h5>{{ __('Host leads the meeting') }}</h5>
                    <p class="explanation">{{ __('If needed, the host can remove unintended participants.') }}</p>
                    <img class="step-img" src="{{ url('/') }}/img/step6.svg" />

                </div>
            </main>
        </div>
    </div>

    <div class="row-full bg-grey">
        <div class="row tight">
            <main class="col-md-12">
                <div class="info" style="padding-bottom: 1rem;">
                    <h2>{{ __('What you get for free') }}</h2>
                    <p>{{ __('A limited version of IRMA-meet is freely available. Soon a non-free version will become available too, with more options and more functionality. For instance, it allows a host to select which information people must disclose about themselves before participation. In the free version:') }}</p>
                </div>
                <div class="two-column-text">
                    <ul>
                        <li>{{ __('the host reveals his/her name and email address') }}</li>
                        <li>{{ __('the meeting link remains valid for at most 24 hours') }}</li>
                        <li>{{ __('only one active link is allowed per host') }}</li>
                        <li>{{ __('participants also need to disclose their name and email address') }}</li>
                        <li>{{ __('at most 6 participants can join a video meeting, including the host') }}</li>
                        <li>{{ __('no BigBlueButton control options are available - in particular, no recordings can be made') }}</li>
                    </ul>
                </div>
            </main>
        </div>
    </div>

    <div class="row tight">
        <main class="col-md-12">
            <div class="info" style="padding-bottom: 1rem;">
                <h2>{{ __('Privacy policy') }}</h2>
                <p>{!! __('privacy_data_controller') !!}</p>
                <p>{{ __('IRMA-meet video calls are protected in two ways.') }}</p>
                <ul>
                    <li>{!! __('privacy_connections_encrypted') !!}</li>
                    <li>{{ __('Hosts and participants in IRMA-meet are authenticated via IRMA. The types of attributes and their trust-worthiness may differ per meeting. It is the responsibility of hosts and participants to judge whether these attributes offer a sufficiently high level of certainty.') }}</li>
                </ul>
                <p>{{ __('IRMA-meet does not offer end-to-end encryption between hosts and participants.') }}</p>
                <p>{!! __('privacy_minimal_data') !!}</p>
                <p>{{ __('privacy_bbb_retention') }}</p>
                <p>{!! __('privacy_policy_changes') !!}</p>
                <p>{!! __('privacy_contact') !!}</p>
            </div>
        </main>
    </div>

    <div class="row-full bg-blue">
        <div class="row img-warp tight">
            <main class="col-xl-4">
                <div class="footer">
                    <h5>{{ __('About') }}</h5>
                    <p>{!! __('about_service_description') !!}</p>
                </div>
            </main>
            <main class="col-xl-4">
                <div class="footer">
                    <h5>{{ __('Open Source') }}</h5>
                    <p>{{ __('The free version of IRMA-meet is open source. You can find the source code and documentation on GitHub.') }}</p>
                    <p>{{ __('Feel free to contribute or set up your own IRMA-meet server.') }} </p>
                    <p> <i class="fa fa-github"></i> {!! __('Visit the <a href="https://github.com/privacybydesign/irma-meet">GitHub repository</a>') !!}</p>
                </div>
            </main>
            <main class="col-xl-4">
                <div class="footer">
                    <h5>{{ __('Contact') }}</h5>
                    <p>{!! __('Visit the <a href="https://www.irma.app">irma.app</a> website to learn more about the IRMA app and go to <a href="https://privacybydesign.foundation">Privacy by Design Foundation</a> for background information.') !!}</p>
                </div>
            </main>
        </div>
    </div>

    </div>

    @include('layout.partials.footer-scripts')

</body>

</html>