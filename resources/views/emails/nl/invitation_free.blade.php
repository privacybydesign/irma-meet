<div>
Beste meneer, mevrouw,
</div>

<br>

<div>
  {{ $hoster_name }} heeft een reservering gemaakt voor een IRMA-meet videogesprek
  met jou als deelnemer. Dit gesprek heeft de link:
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
  Je kunt deelnemen aan het videogesprek via bovenstaande link. Je zult dan moeten inloggen op irma-meet.nl, met je naam en e-mailadres uit je IRMA app. Je naam en e-mail zijn zichtbaar voor iedereen in het gesprek.
</div>

<br>

<div>
  Staat de IRMA-app nog niet op je mobiel? Installeer die app dan eerst en laad je e-mailadres en naam (uit de <a href="https://services.nijmegen.nl/irma/gemeente/start">Basisregistratie
  Personen BRP</a> of uit <a href="https://privacybydesign.foundation/uitgifte/social/linkedin/">LinkedIn</a>). Voor meer informatie, zie de website: <a href="https://irma.app">https://irma.app</a>.</li>
</div>

<br>

<div>
  Veel plezier met je videogesprek met zekerheid!
</div>

<div>
  IRMA-meet team.
</div>
