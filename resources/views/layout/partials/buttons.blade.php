<button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('{{ url('/') }}/irma_auth/start/default', '{{ url('/') }}/irma_session/create/free');">
	<img class="img-fluid" src="{{ url('/') }}/img/team-icon.svg"> {{ __('Create video meeting') }}
</button>
<br>
<p style="font-size: 14px; color:#1B4D8C;"><br>{{ __('For use in the Netherlands:') }}</p>
<button type="button" class="btn btn-secondary btn-lg btn-blue" onclick="startInvitation('{{ url('/') }}/irma_auth/start/teacher', '{{ url('/') }}/irma_session/create/exam');">
	<img class="img-fluid" src="{{ url('/') }}/img/exam-icon.svg">{{ __('Start oral video exam') }}
</button>
<br><br>
<button type="button" class="btn btn-secondary btn-lg btn-blue" onclick="startInvitation('{{ url('/') }}/irma_auth/start/medical', '{{ url('/') }}/irma_session/create/medical_consult');">
	<img class="img-fluid" src="{{ url('/') }}/img/medical_icon.svg">{{ __('Start medical video consult') }}
</button>