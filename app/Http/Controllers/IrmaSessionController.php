<?php

namespace App\Http\Controllers;

use App\Mail\Invitation;
use BigBlueButton\BigBlueButton;
use BigBlueButton\Parameters\CreateMeetingParameters;
use BigBlueButton\Parameters\JoinMeetingParameters;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Mail;
use Session;

class IrmaSessionController extends Controller {
	/**
	 * Create a new controller instance.
	 *
	 * @return void
	 */
	public function __construct() {
		//
	}

	/**
	 * Show the store session screen.
	 *
	 * @return \Illuminate\Contracts\Support\Renderable
	 */
	public function create() {
		$validated_email = Session::get('validated_email', '');
		return view('layout/irma-session')->with('validated_email', $validated_email);
	}

	/**
	 * Store a session.
	 *
	 * @return \Illuminate\Contracts\Support\Renderable
	 */
	public function store(Request $request) {
		$validated_email = Session::get('validated_email', '');
		$validatedData = $request->validate([
			'meeting_name' => 'required|max:255',
			'hoster_name' => 'required',
			'hoster_email_address' => 'required|in:' . $validated_email,
			'start_time' => '', //not yet used
			'invitation_note' => 'required|max:255',
			'participant_email_address1' => 'required|email',
			'participant_email_address2' => 'nullable|email',
			'participant_email_address3' => 'nullable|email',
			'participant_email_address4' => 'nullable|email',
			'participant_email_address5' => 'nullable|email',
			'participant_email_address6' => 'nullable|email',
		]);
		$uniqueId = bin2hex(openssl_random_pseudo_bytes(24));
		//we use another session id for bbb so the bbb session id is not exposed in the url
		$bbbSessionId = bin2hex(openssl_random_pseudo_bytes(12));
		$validatedData = array_merge($validatedData, ['irma_session_id' => $uniqueId, 'start_time' => now(), 'bbb_session_id' => $bbbSessionId]);
		$irma_session = \App\IrmaMeetSessions::create($validatedData);

		//for all participants store data
		$participantsEmails = [];
		for ($i = 1; $i < 7; $i++) {
			if ($validatedData['participant_email_address' . $i]) {
				$irmaParticipant = [
					'irma_session_id' => $irma_session['irma_session_id'],
					'email_address' => $validatedData['participant_email_address' . $i],
					'authentication' => 1,
				];
				$irma_participants = \App\IrmaMeetParticipants::create($irmaParticipant);
				array_push($participantsEmails, $validatedData['participant_email_address' . $i]);
			}
		}

		$participantInvitationLink = url('/') . "/irma_session/join/participant/" . $uniqueId;
		$hosterInvitationLink = url('/') . "/irma_session/join/hoster/" . $uniqueId;

		//Create session in BBB
		$bbb = new BigBlueButton();
		$createParams = new CreateMeetingParameters($bbbSessionId, $validatedData['meeting_name']);
		//we use an md5 hash from the bbb session id. The bbb session id is not exposed, so md5 should be enough protection
		$createParams->setAttendeePassword(md5('participant' . $bbbSessionId));
		//$createParams->setAttendeePassword('participant');
		$createParams->setModeratorPassword(md5('hoster' . $bbbSessionId));
		//$createParams->setModeratorPassword('hoster');
		$response = $bbb->createMeeting($createParams);

		if ($response->getReturnCode() == 'FAILED') {
			return 'Can\'t create room! please contact our administrator.';
		} else {
			$joinParams = new JoinMeetingParameters($bbbSessionId, $validatedData['hoster_name'], md5('hoster' . $bbbSessionId));
			$joinParams->setRedirect(true);
			$url = $bbb->getJoinMeetingURL($joinParams);
			//Send mail with links to hoster
			Mail::to($validatedData['hoster_email_address'])
				->send(new Invitation([
					'hoster_name' => $validatedData['hoster_name'],
					'invitation_note' => $validatedData['invitation_note'],
					'hoster_invitation_link' => $hosterInvitationLink,
					'participant_invitation_link' => $participantInvitationLink,
				]));
			$mainContent = 'Meeting is successfully validated and data has been saved';
			return view('layout/irma-session-success', [
				'mainContent' => $mainContent,
				'participantInvitationLink' => $participantInvitationLink,
				'hosterInvitationLink' => $hosterInvitationLink,
			]
			);
		}

	}

	/**
	 * Show the store success session screen.
	 *
	 * @return \Illuminate\Contracts\Support\Renderable
	 */
	public function success() {
		$mainContent = 'Meeting is successfully validated and data has been saved';
		return view('layout/irma-session-success')->with('mainContent', $mainContent);
	}

	/**
	 * Redirect authenticated participant to session
	 *
	 * @return redirect to bbb
	 */
	public function join_participant($irmaSessionId) {
		//TODO: fix participants link
		//get irmasessionid from link
		//var_dump($irmaSessionId);
		//search email address

		//get validated email address and name

		//get bbb session id

		$bbbSessionId = '';
		//redirect to bbb
		$joinParams = new JoinMeetingParameters($bbbSessionId, 'joeri', 'yy');
		$joinParams->setRedirect(true);
		// Join the meeting by redirecting the user to the generated URL
		$url = $bbb->getJoinMeetingURL($joinParams);
		return \Redirect::to($url);
	}

	/**
	 * Redirect authenticated hoster to session
	 *
	 * @return redirect to bbb
	 */
	public function join_hoster($irmaSessionId) {
		$irmaSession = \App\IrmaMeetSessions::where('irma_session_id', $irmaSessionId)->first();
		$meetingName = $irmaSession->meeting_name;
		$hosterName = $irmaSession->hoster_name;
		$hosterEmailAddress = $irmaSession->hoster_email_address;
		$bbbSessionId = $irmaSession->bbb_session_id;
		$bbb = new BigBlueButton();
		//redirect to bbb
		$joinParams = new JoinMeetingParameters($bbbSessionId, $hosterName . '<' . $hosterEmailAddress . '>', md5('hoster' . $bbbSessionId));
		$joinParams->setRedirect(true);
		// Join the meeting by redirecting the user to the generated URL
		$url = $bbb->getJoinMeetingURL($joinParams);
		return \Redirect::to($url);
	}
}
