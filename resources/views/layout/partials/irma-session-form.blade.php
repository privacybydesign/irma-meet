<!-- /resources/views/layout/irma-session-form.blade.php -->

<h1>Create meeting</h1>

<div>
<form method="post" action="{{ route('irma_session.store') }}">
	{{ csrf_field() }}

<div>
	<label for="meeting_name">Meeting name</label>
	<input id="meeting_name" name="meeting_name"  value="{{ old('meeting_name') }}" type="text" class="@error('meeting_name') is-invalid @enderror">
</div>

<div>
	<label for="hoster_name">Hoster name</label>
	<input id="hoster_name" name="hoster_name"  value="{{ old('hoster_name') }}" type="text" class="@error('hoster_name') is-invalid @enderror">
</div>

<div>
	<label for="hoster_email_address">Hoster email address</label>
	<input id="hoster_email_address"  value="{{ $validated_email }}" name="hoster_email_address" type="text" class="@error('hoster_email_address') is-invalid @enderror" readonly>
</div>

<div>
	<label for="invitation_note">Invitation note</label>
	<textarea id="invitation_note" name="invitation_note" class="@error('invitation_note') is-invalid @enderror">{{ old('invitation_note') }}</textarea>
</div>

@for ($i = 1; $i < 7; $i++)
<div>
	<label for="participant_email_address{{ $i }}">Participant {{ $i }} email address</label>
	<input id="participant_email_address{{ $i }}"  value="{{ old('participant_email_address'.$i) }}" name="participant_email_address{{ $i }}" type="text" class="@error('participant_email_address{{ $i }}') is-invalid @enderror">
</div>
@endfor

 	<button type="submit" class="btn btn-primary">Send invite</button>
</form>
</div>


@if ($errors->any())
    <div class="alert alert-danger">
        <ul>
            @foreach ($errors->all() as $error)
                <li>{{ $error }}</li>
            @endforeach
        </ul>
    </div>
@endif

<!-- Create Post Form -->