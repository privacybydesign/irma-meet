<?php

namespace App\Http\Controllers;

use App\Mail\Invitation;
use BigBlueButton\BigBlueButton;
use BigBlueButton\Parameters\CreateMeetingParameters;
use BigBlueButton\Parameters\JoinMeetingParameters;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Mail;
use Session;

class IrmaSessionController extends Controller
{
    /**
     * Create a new controller instance.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    /**
     * Show the store session screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function authenticate($url)
    {
        $token = Session::get('irma_session_token', '');
        if ($token === '') {
            return view('layout/irma_session_authenticate')->with('url', urldecode(urldecode($url)));
        } else {
            echo $url;
            return redirect(urldecode(urldecode($url)));
        }
    }

    /**
     * Show the store session screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function create()
    {
        $validated_email = Session::get('validated_email', '');
        return view('layout/irma_session_create')->with('validated_email', $validated_email);
    }

    /**
     * Store a session.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function store(Request $request)
    {
        $validated_email = Session::get('validated_email', '');
        $validatedData = $request->validate([
            'meeting_name' => 'required|max:255',
            'hoster_name' => 'required',
            'hoster_email_address' => 'required|in:' . $validated_email,
            'start_time' => '', //not yet used
            'invitation_note' => 'nullable|max:255',
            'participant_email_address1' => 'nullable|email',
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
            if (in_array('participant_email_address' . $i, $validatedData)) {
                $irmaParticipant = [
                    'irma_session_id' => $irma_session['irma_session_id'],
                    'email_address' => $validatedData['participant_email_address' . $i],
                    'authentication' => 1,
                ];
                $irma_participants = \App\IrmaMeetParticipants::create($irmaParticipant);
                array_push($participantsEmails, $irma_participants->email_address);
            }
        }

        $participantLink = url('/') . "/irma_session/join/participant/" . $uniqueId;
        $hosterLink = url('/') . "/irma_session/join/hoster/" . $uniqueId;

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
            //Send mail with links to hoster
            Mail::to($validatedData['hoster_email_address'])
                ->send(new Invitation([
                    'hoster_name' => $validatedData['hoster_name'],
                    'invitation_note' => in_array('invitation_note', $validatedData) ? $validatedData['invitation_note'] : '',
                    'hoster_invitation_link' => $hosterLink,
                    'participant_invitation_link' => $participantLink,
                ]));
            if (false && count($participantsEmails)) {
                //Send mail with links to participants
                //TODO: create different invitation for participants
                Mail::from($validatedData['hoster_email_address'])
                ->bcc($participantsEmails)
                ->send(new Invitation([
                    'hoster_name' => $validatedData['hoster_name'],
                    'invitation_note' => $validatedData['invitation_note'],
                    'participant_invitation_link' => $participantLink,
                ]));
            }
            $mainContent = 'Meeting is successfully validated and data has been saved';
            return view(
                'layout/irma_session_success',
                [
                    'mainContent' => $mainContent,
                    'participantInvitationLink' => $participantLink,
                    'hosterInvitationLink' => $hosterLink,
                ]
            );
        }
    }

    /**
     * Redirect authenticated participant to session
     *
     * @return redirect to bbb
     */
    public function joinParticipant($irmaSessionId)
    {
        $irmaSession = \App\IrmaMeetSessions::where('irma_session_id', $irmaSessionId)->first();
        $bbbSessionId = $irmaSession->bbb_session_id;
        $participantEmail = Session::get('validated_email', '');
        if (($participantEmail !== '') && ($irmaSessionId !== '')) {
            $bbb = new BigBlueButton();
            //redirect to bbb
            $joinParams = new JoinMeetingParameters($bbbSessionId, '[' . $participantEmail . ']', md5('participant' . $bbbSessionId));
            $joinParams->setRedirect(true);
            // Join the meeting by redirecting the user to the generated URL
            $url = $bbb->getJoinMeetingURL($joinParams);
            return \Redirect::to($url);
        } else {
            $mainContent = 'Your email address could not be verified by IRMA.';
            return view('layout/irma_session_error')->with('mainContent', $mainContent);
        }
    }

    /**
     * Redirect authenticated hoster to session
     *
     * @return redirect to bbb
     */
    public function joinHoster($irmaSessionId)
    {
        $irmaSession = \App\IrmaMeetSessions::where('irma_session_id', $irmaSessionId)->first();
        $hosterName = $irmaSession->hoster_name;
        $hosterEmailAddress = $irmaSession->hoster_email_address;
        $bbbSessionId = $irmaSession->bbb_session_id;
        $bbb = new BigBlueButton();
        //redirect to bbb
        $joinParams = new JoinMeetingParameters($bbbSessionId, $hosterName . ' [' . $hosterEmailAddress . ']', md5('hoster' . $bbbSessionId));
        $joinParams->setRedirect(true);
        // Join the meeting by redirecting the user to the generated URL
        $url = $bbb->getJoinMeetingURL($joinParams);
        return \Redirect::to($url);
    }
}
