                        {{ __('Please authenticate with IRMA before entering the room.') }}'
                        <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation('../../irma_auth/start', '{{ $url }}');">{{ __('Authenticate with IRMA') }}</button>
                        </button>