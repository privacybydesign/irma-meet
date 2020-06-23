                        <p>{{ __('Please authenticate with IRMA before entering the room.') }}</p>
                        <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('../../irma_auth/start/{{ $disclosureTypeHost }}', '{{ url('/') }}/irma_session/join_host/{{ $irmaSessionId }}');">{{ __('Join as host') }}
                        </button>                        
                        <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('../../irma_auth/start/{{ $disclosureTypeParticipant }}', '{{ url('/') }}/irma_session/join_participant/{{ $irmaSessionId }}');">{{ __('Join as participant') }}
                        </button>