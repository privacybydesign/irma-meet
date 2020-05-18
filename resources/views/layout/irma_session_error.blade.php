<!DOCTYPE html>

<html lang="en">

<head>

    @include('layout.partials.head')

</head>

<body>


    <div class="container-fluid h-100 d-flex flex-column">

        <div class="row">
            <main class="col-md-4 bg-blue">
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
            <main class="col-md-4 bg-blue">
                <div class="title">
                </div>
            </main>
            <main class="col-md-4 nopadding-md bg-warmwhite text-center">
                <div class="card">
                    <div class="card-header">
                        Error
                    </div>
                    <div class="card-body">
                        {{ $mainContent }}
                    </div>
                </div>

            </main>
            <main class="col-md-4 bg-warmwhite">
                <div class="content">
                    <p class="cl-grey slogan">AUTHENTICATED VIDEO CALLS</p>
                    <h1>Meet online, without surprises!</h1>
                    <p>Only people who prove who they are can participate in IRMA-meet video. Chat with certainty, without being fooled, in your business meeting, (medical) consult, or online oral exam.</p>
                    <div class="">
                    </div>
                    <div id="errorbox" class="alert alert-danger alert-dismissible fade">
                        <span id="errorboxtext"></span>
                        <button type="button" class="close" data-dismiss="alert">&times;</button>
                    </div>
                </div>
            </main>
        </div>

        <div class="row-full bg-grey">
            <div class="row img-warp tight">
                <main class="col-md-7 bg-grey">
                    <div class="info">
                        <h2>Know who you are talking to</h2>
                        <p>Participants in IRMA-meet video calls need to prove who they are by disclosing relevant information (like their name and e-mail address) about themselves before they can enter a video chat.</p>
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
                        <img class="img-fluid irma_logo_info" src="{{ url('/') }}/img/irma_logo.svg" />
                    </div>
                </main>
                <main class="col-md-7">
                    <div class="info">
                        <h2>Prove who you are with IRMA</h2>
                        <p>You can prove who you are via <a href="https://www.irma.app">IRMA</a>. It is a
                            freely available open source identity app; it is a wallet on your phone for collecting personal
                            data. This data is stored only on your phone, and nowhere else. You control where you disclose
                            it.</p>
                    </div>
                    <div>
                        <img class="img-fluid irma_logo_small" src="{{ url('/') }}/img/irma_logo.svg" />
                    </div>
                </main>
            </div>
        </div>

        <div class="row img-warp tight">
            <main class="col-md-7">
                <div class="info">
                    <div>
                        <h2>Get ready</h2>
                    </div>
                    <p>All IRMA-meet users need to have the IRMA app on their phone, filled with the relevant personal
                        data. Which data is required may depend on the video call.</p>
                    <p>If you don't have IRMA yet, install it from the IOS appstore or Google Play store and load the
                        required personal data, like your name or e-mail address.</p>
                    <p>At this early stage, there is only one type of free meeting available, which requires your e-mail
                        address and your name in IRMA. Your name can be obtained either from your LinkedIn profile or
                        for users in the Netherlands from the Dutch citizen registry.</p>
                </div>
            </main>
            <main class="col-md-5">
                <div class="">
                    <img class="img-fluid getready" src="{{ url('/') }}/img/getready.svg" />
                </div>
            </main>
        </div>

        <div class="row-full bg-light-grey">
            <div class="info">
                <h2 style="text-align: center;">How it works</h2>
            </div>
            <div class="row img-warp tight">
                <main class="col-lg-4">
                    <div class="step">
                        <h3>Step 1</h3>
                        <h5>Host proves who (s)he is</h5>
                        <p class="explanation">The host initiates a new meeting on this website and makes
                            herself/himself known,
                            via the IRMA app.</p>
                        <img class="step-img" src="{{ url('/') }}/img/irma_logo.svg" />
                    </div>
                </main>
                <main class="col-lg-4">
                    <div class="step">
                        <h3>Step 2</h3>
                        <h5>Host chooses meeting name</h5>
                        <p class="explanation">The host chooses a name for the meeting and obtains a link for the
                            meeting.</p>
                        <img class="step-img" src="{{ url('/') }}/img/step2.svg" />

                    </div>
                </main>
                <main class="col-lg-4">
                    <div class="step">
                        <h3>Step 3</h3>
                        <h5>Host distributes meeting link</h5>
                        <p class="explanation">The host distributes the meeting link to the intended participants, for
                            instance via email or direct messaging.</p>
                        <img class="step-img" src="{{ url('/') }}/img/step3.svg" />

                    </div>
                </main>
            </div>
            <div class="row img-warp tight">
                <main class="col-lg-4">
                    <div class="step">
                        <h3>Step 4</h3>
                        <h5>Participants prove who they are</h5>
                        <p class="explanation">When a participant clicks on the link (s)he is asked to prove who (s)he
                            via IRMA.</p>
                        <img class="step-img" src="{{ url('/') }}/img/irma_logo.svg" />

                    </div>
                </main>
                <main class="col-lg-4">
                    <div class="step">
                        <h3>Step 5</h3>
                        <h5>Participants enter meeting</h5>
                        <p class="explanation">Each participant joins the meeting automatically. The verified identity
                            of everyone in the meeting is visible.</p>
                        <img class="step-img" src="{{ url('/') }}/img/step5.svg" />

                    </div>
                </main>
                <main class="col-lg-4">
                    <div class="step">
                        <h3>Step 6</h3>
                        <h5>Host leads the meeting</h5>
                        <p class="explanation">If needed, the host can remove unintended participants.</p>
                        <img class="step-img" src="{{ url('/') }}/img/step6.svg" />

                    </div>
                </main>
            </div>
        </div>

        <div class="row-full bg-grey">
            <div class="row tight">
                <main class="col-md-12">
                    <div class="info" style="padding-bottom: 1rem;">
                        <h2>What you get for free</h2>
                        <p>A limited version of IRMA-meet is freely available. Soon a non-free version will become
                            available too, with more options and more functionality. For instance, it allows a host to
                            select which information people must disclose about themselves before participation. In the
                            free version:</p>
                    </div>
                    <div class="two-column-text">
                        <ul>
                            <li>the host reveals his/her name and email address</li>
                            <li>the meeting link remains valid for at most 24 hours</li>
                            <li>only one active link is allowed per host</li>
                            <li>participants also need to disclose their name and email address</li>
                            <li>at most 6 participants can join a video meeting, including the host</li>
                            <li>no BigBlueButton control options are available - in particular, no recordings can be made</li>
                        </ul>
                    </div>
                </main>
            </div>
        </div>

        <div class="row tight">
            <main class="col-md-12">
                <div class="info" style="padding-bottom: 1rem;">
                    <h2>Privacy policy</h2>
                    <p>The IRMA-meet service is offered jointly by the company <a
                            href="https://www.procolix.com">ProcoliX</a> and the <a
                            href="https://privacybydesign.foundation">Privacy by Design Foundation</a>. Both <a
                            href="https://www.procolix.com">ProcoliX</a> and
                        the foundation are jointly the data controller, but only <a
                            href="https://www.procolix.com">ProcoliX</a> is
                        the data processor: it hosts the servers (for BigBlueButton video and
                        for IRMA authentication, and for the webpages) that process the
                        relevant (personal) data. The <a href="https://privacybydesign.foundation">Privacy by Design
                            Foundation</a> does not
                        process any personal data as part of IRMA-meet.</p>
                    <p>IRMA-meet video calls are protected in two ways.</p>
                    <ul>
                        <li> The connections that carry the video and audio (and other) signals between
                            each user and <a href="https://www.procolix.com">ProcoliX</a> are encrypted, using the
                            BigBlueButton
                            software.</li>
                        <li>Hosts and participants in IRMA-meet are authenticated via IRMA. The
                            types of attributes and their trust-worthiness may differ per
                            meeting. It is the responsibility of hosts and participants to judge
                            whether these attributes offer a sufficiently high level of certainty.</li>
                    </ul>
                    <p>IRMA-meet does not offer end-to-end encryption between hosts and
                        participants.</p>
                    <p><a href="https://www.procolix.com">ProcoliX</a> processess only minimal personal data required
                        for the
                        IRMA-meet functionality. In particular, it does not record or store
                        content from meetings, like video, audio, chat, presentations, shared
                        screens etc. The functionality (for hosts) to record meetings within
                        IRMA-meet is not available at this early stage. Also, personal data
                        from authentication with IRMA, for getting access to IRMA-meet calls,
                        is not stored by <a href="https://www.procolix.com">ProcoliX</a> beyond the meeting. Other,
                        technical
                        personal data such as IP-addresses will be stored and deleted in
                        accordance <a href="https://www.procolix.com">ProcoliX</a>'s general company policy. When no
                        technical
                        abnormalities occur, they are deleted after 2 days.</p>
                    <p><a href="https://www.procolix.com">ProcoliX</a> may ask hosts of meetings if their name and/or
                        email address
                        can be retained in order to stay informed about future IRMA-meet
                        developments. Such processing of personal data happens on the (legal)
                        basis of the consent.</p>
                    <p>Statistical information about IRMA-meet video chats may be aggregated
                        (duration, numbers of participants, numbers of meetings per
                        day/week/month/year, origin etc.) and may be used by <a
                            href="https://www.procolix.com">ProcoliX</a> and by
                        the <a href="https://privacybydesign.foundation">Privacy by Design Foundation</a> for
                        documentation and
                        publications. If needed, IRMA-meet data can be used in a minimal way
                        for technical purposes, such as maintenance and performance tuning,
                        but they will be destroyed after a maximum of 2 days.</p>
                    <p>Technical changes in the IRMA-meet system, or possible new services,
                        may lead to adaptation of this privacy policy. <a href="https://www.procolix.com">ProcoliX</a>
                        and the
                        <a href="https://privacybydesign.foundation">Privacy by Design Foundation</a> reserve the right
                        to make such changes
                        and will publish the adapted privacy policy via this website as soon
                        as possible.</p>
                    <p>For questions, remarks, or complaints about this data processing by
                        <a href="https://www.procolix.com">ProcoliX</a> for IRMA-meet functionality, please contact the
                        <a href="https://www.procolix.com">ProcoliX</a> at
                        info@irma-meet.nl. For complaints about the IRMA-meet data processing
                        one can also contact the Data Protection Authority of the Netherlands.</p>
                </div>
            </main>
        </div>

        <div class="row-full bg-blue">
            <div class="row img-warp tight">
                <main class="col-xl-4">
                    <div class="footer">
                        <h5>About</h5>
                        <p>The IRMA-meet video chat service is based on the open source video platform <a
                                href="https://bigbluebutton.org/">BigBlueButton</a>. It
                            runs at the hosting company <a href="https://www.procolix.com">ProcoliX</a> in the
                            Netherlands. It has been developed jointly with
                            the <a href="https://privacybydesign.foundation">Privacy by Design Foundation</a> with
                            support from <a href="https://www.sidnfonds.nl">SIDNfonds</a>.</p>
                    </div>
                </main>
                <main class="col-xl-4">
                    <div class="footer">
                        <h5>Open Source</h5>
                        <p>The free version of IRMA-meet is open source. You can find the source code and documentation
                            on GitHub.</p>
                        <p>Feel free to contribute or set up your own IRMA-meet server. </p>
                        <p> <i class="fa fa-github"></i> Visit the <a
                                href="https://github.com/privacybydesign/irma-meet">GitHub repository</a></p>
                    </div>
                </main>
                <main class="col-xl-4">
                    <div class="footer">
                        <h5>Contact</h5>
                        <p>Visit the <a href="https://www.procolix.com">ProcoliX website</a> for more information about
                            the hosting company.</p>

                        <p>Visit the <a href="https://privacybydesign.foundation">Privacy by Design Foundation</a> for
                            background information and the <a href="https://www.irma.app">IRMA website</a> to learn more
                            about the IRMA app.</p>

                        <p> <i class="fa fa-twitter"></i>
                            Follow us on twitter:</p>

                        <p>IMRA-meet: <a href="https://twitter.com/IRMA_meet/">@IRMA_meet</a></p>
                        <p>Privacy by Design Foundation: <a href="https://twitter.com/IRMA_privacy/">@IRMA_privacy</a>
                        </p>
                    </div>
                </main>
            </div>
        </div>

    </div>

    @include('layout.partials.footer-scripts')

</body>

</html>