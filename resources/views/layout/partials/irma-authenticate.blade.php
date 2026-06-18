                        {{ __('Please authenticate with IRMA before entering the room.') }}
                        <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation({{ \Illuminate\Support\Js::from('../../irma_auth/start/' . $disclosureType) }}, {{ \Illuminate\Support\Js::from($url) }});">{{ __('Authenticate with IRMA') }}
                        </button>
