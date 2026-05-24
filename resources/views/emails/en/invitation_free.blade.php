<div>
  Dear Mr / Ms,
</div>

<br>

<div>
  {{ $hoster_name }} has made a reservation for a Yivi Meet
  video call with you as a participant. This meeting has the link:
</div>

<br>

<div>
  <a href="{{ $invitation_link }}">{{ $invitation_link }}</a><br />
</div>

@if (!empty($invitation_note))
<div><br />
*******************
</div>
<div>
  {{ $invitation_note }}
</div>
<div>
*******************
</div>
@endif

<div>
  You can join the meeting by following the above link. You will then have to login at yivi-meet.nl, with your name and email address from your Yivi app to join the meeting. Your name and email address will then be visible to everyone else in the meeting.
</div>

<br>

<div>
  In case you do not have the Yivi app on your phone yet, you need to get the Yivi app first. You need to fill it with at least your email address and name (from <a href="https://privacybydesign.foundation/issuance/social/linkedin/">LinkedIn</a> or from the <a href="https://services.nijmegen.nl/irma/gemeente/start"> Dutch civil registry BRP</a>). You can find out more at the website <a href="https://yivi.app">https://yivi.app</a>.</li>
</div>
<br>

<div>
  Enjoy your authenticated meeting!
</div>

<div>
  Yivi Meet team.
</div>
