                        <div id="card-image">
                        <p>{{ __('Please authenticate with IRMA before entering the room.') }}</p>

                        <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('../../irma_auth/start/{{ $disclosureTypeHost }}', '{{ url('/') }}/irma_session/join_host/{{ $irmaSessionId }}');">{{ $hosterMessage }}
                        </button>                        
                        <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('../../irma_auth/start/{{ $disclosureTypeParticipant }}', '{{ url('/') }}/irma_session/join_participant/{{ $irmaSessionId }}');">{{ $clientMessage }}
                        </button>
                    </div>
                        <section class="irma-web-form" style="display: none" />
                          