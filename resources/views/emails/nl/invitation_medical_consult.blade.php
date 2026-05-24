<div>
Beste mevrouw, meneer,
</div>

<br>

<div>
Beste student,
</div>

<br>

<div>
  Zorgverlener {{ $hoster_name }} heeft een reservering gemaakt voor een Yivi Meet medisch videogesprek
  met jou als patiënt. Dit gesprek heeft de link:
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
  Je kunt deelnemen aan het videogesprek via bovenstaande link. Je zult dan moeten inloggen op yivi-meet.nl, met je persoonsgegevens uit je Yivi app. Je naam, geboortedatum, BSN en e-mail zijn zichtbaar voor iedereen in het gesprek.
</div>

<br>

<div>
  Staat de Yivi-app nog niet op je mobiel? Installeer die app dan eerst en laad je persoonsgegevens via de <a href="https://services.nijmegen.nl/irma/gemeente/start">Basisregistratie Personen BRP</a>. Voor meer informatie, zie de website: <a href="https://yivi.app">https://yivi.app</a>.</li>
</div>

<br>

<div>
  Een goed geprek gewenst.
</div>

<div>
  Yivi Meet team.
</div>
