<div>
  Dear Mr / Ms,
</div>

<br>

<div>
  {{ $hoster_name }} has made a reservation for an IRMA-meet
  video exam with you as a student. This meeting has the link:
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
  You can join the meeting by following the above link. You will then have to login at irma-meet.nl, with your academic personal data from SURFconext from your IRMA app. Your name and student number will then be visible to everyone else in the meeting.
</div>

<br>

<div>
  In case you do not have the IRMA app on your phone yet, you need to get the IRMA app first. You need to fill it with at least
  your academic personal data from <a href="https://privacybydesign.foundation/issuance/surfnet/surfnet/">SURFconext</a> . You can find out more at the website <a href="https://irma.app">https://irma.app</a>.</li>
</div>

<br>

<div>
  Enjoy your authenticated meeting!
</div>

<div>
  IRMA-meet team.
</div>