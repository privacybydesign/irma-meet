<!-- /resources/views/layout/irma-session-form.blade.php -->

		<form method="post" class="" action="{{ route('irma_session.store') }}">
			{{ csrf_field() }}
			<input id="meeting_type" name="meeting_type"  value="exam" type="hidden">
		<div class="form-group">
			<label class="control-label"  for="meeting_name">{{ __('Meeting name') }}</label>
				<input id="meeting_name" name="meeting_name"  value="{{ old('meeting_name') }}" type="text" class="form-control @error('meeting_name') is-invalid @enderror">
		</div>

		<div class="form-group">
			<label class="control-label"  for="hoster_name">{{ __('Hoster name') }}</label>
				<input id="hoster_name" name="hoster_name"  value="{{ $validated_name }}" type="text" class="form-control @error('hoster_name') is-invalid @enderror" readonly>
		</div>

		<div class="form-group">
			<label class="control-label"  for="hoster_email_address">{{ __('Hoster email address') }}</label>
				<input id="hoster_email_address"  value="{{ $validated_email }}" name="hoster_email_address" type="text" class="form-control @error('hoster_email_address') is-invalid @enderror" readonly>
		</div>

		<div>
			<label for="invitation_note">{{ __('Invitation note') }}</label>
			<textarea id="invitation_note" name="invitation_note" class="@error('invitation_note') is-invalid @enderror">{{ old('invitation_note') }}</textarea>
		</div>

		<div>
			<label for="participant_email_address1">{{ __('Student email address') }}</label>
			<input id="participant_email_address1"  value="{{ old('participant_email_address1') }}" name="participant_email_address1" type="text" class="@error('participant_email_address1') is-invalid @enderror">
		</div>
		    <div class="form-group">        
			 	<button type="submit" class="btn btn-primary">{{ __('Send invite') }}</button>
		    </div>


		</form>


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