<?php

namespace App\Http\Controllers;

use App\Mail\Invitation;
use BigBlueButton\BigBlueButton;
use BigBlueButton\Parameters\CreateMeetingParameters;
use BigBlueButton\Parameters\JoinMeetingParameters;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Mail;
use Session;
use Config;

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
     * Show the create session screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function create($meetingType)
    {
        $disclosureType = Config::get('meeting-types.' . $meetingType . '.irma_disclosure');
        $disclosureTypeHost = Config::get('meeting-types.' . $meetingType . '.irma_disclosure_host', $disclosureType);
        $validatedEmail = strtolower(Session::get(Config::get('disclosure-types.' . $disclosureTypeHost . '.email'), ''));
        $validatedName = $this->_get_name($disclosureTypeHost);
        $form = view('layout.partials.irma-session-form-' . $meetingType)->with([
            'validated_email' => $validatedEmail,
            'validated_name' => $validatedName
        ])->render();
        $buttons = view('layout.partials.buttons')->render();
        return view('layout/mainlayout')->with(
            [
                'message' => $form,
                'title' => __('Create video meeting'),
                'buttons' => $buttons
                ]
        );
    }

    /**
     * Store a session.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function store(Request $request)
    {
        $meetingType = $request->get('meeting_type');
        $disclosureType = Config::get('meeting-types.' . $meetingType . '.irma_disclosure');
        $disclosureTypeHost = Config::get('meeting-types.' . $meetingType . '.irma_disclosure_host', $disclosureType);
        $validatedEmail = strtolower(Session::get(Config::get('disclosure-types.' . $disclosureTypeHost . '.email'), ''));
        //TODO: find a way to have infinite participan_email_addresses in validation
        $validatedData = $request->validate([
            'meeting_name' => 'required|max:255',
            'hoster_name' => 'required',
            'hoster_email_address' => 'required|in:' . sprintf("%s", $validatedEmail),
            'start_time' => '', //not yet used
            'invitation_note' => 'required_with:participant_email_address1|max:255',
            'participant_email_address1' => 'nullable|email',
            'participant_email_address2' => 'nullable|email',
            'participant_email_address3' => 'nullable|email',
            'participant_email_address4' => 'nullable|email',
            'participant_email_address5' => 'nullable|email',
            'participant_email_address6' => 'nullable|email',
            'meeting_type' => 'required'
        ]);
        $uniqueId = bin2hex(openssl_random_pseudo_bytes(4));
        //we use another session id for bbb so the bbb session id is not exposed in the url
        $bbbSessionId = bin2hex(openssl_random_pseudo_bytes(12));
        $validatedData = array_merge($validatedData, ['irma_session_id' => $uniqueId, 'start_time' => now(), 'bbb_session_id' => $bbbSessionId]);
        $irma_session = \App\IrmaMeetSessions::create($validatedData);

        // for all participants store data
        $participantsEmails = [];
        for ($i = 1; $i < 7; $i++) {
            // if (in_array('participant_email_address' . 1, $validatedData)) {
            if (!empty($validatedData['participant_email_address'. $i])) {
                $irmaParticipant = [
                    'irma_session_id' => $irma_session['irma_session_id'],
                    'email_address' => $validatedData['participant_email_address' . $i],
                    'authentication' => 1,
                ];
                $irma_participants = \App\IrmaMeetParticipants::create($irmaParticipant);
                array_push($participantsEmails, strtolower($irma_participants->email_address));
            }
        }

        $invitationLink = url('/') . "/irma_session/join/" . $uniqueId;

        //send emails to hoster and participants
        $this->_send_mail($validatedData, $invitationLink);
        $mainContent = '<p>' . __('Meeting is successfully validated and data has been saved.') . '</p>';
        $mainContent .= '<p>' . __('Use the link below to share with your participants:') . '</p>';
        $mainContent .= '<a href="' . $invitationLink . '">' . $invitationLink . '</a>';
        $buttons = view('layout.partials.buttons')->render();
        return view('layout/mainlayout')->with(
            [
                    'message' => $mainContent,
                    'title' =>  __('Success'),
                    'buttons' => $buttons
                    ]
        );
    }

    /**
     * Redirect authenticated hoster to session
     *
     * @return redirect to bbb
     */
    public function join($irmaSessionId)
    {
        $irmaSession = \App\IrmaMeetSessions::where('irma_session_id', $irmaSessionId)->first();
        $meetingType = $irmaSession->meeting_type;
        $hosterEmailAddress = $irmaSession->hoster_email_address;
        $disclosureType = Config::get('meeting-types.' . $meetingType . '.irma_disclosure');
        $disclosureTypeHost = Config::get('meeting-types.' . $meetingType . '.irma_disclosure_host', $disclosureType);
        
        //Create session in BBB
        $bbbSessionId = $irmaSession->bbb_session_id;
        $meeting_name = $irmaSession->meeting_name;
        $bbb = new BigBlueButton();
        $createParams = new CreateMeetingParameters($bbbSessionId, $meeting_name);
        //we use an md5 hash from the bbb session id. The bbb session id is not exposed, so md5 should be enough protection
        $createParams->setAttendeePassword(hash('sha256', 'participant' . $bbbSessionId));
        $createParams->setModeratorPassword(hash('sha256', 'hoster' . $bbbSessionId));
        $isMeetingRunning = $bbb->isMeetingRunning($createParams);
        if (! $isMeetingRunning->isRunning()) {
            $response = $bbb->createMeeting($createParams);
            if ($response->getReturnCode() == 'FAILED') {
                return __('Can\'t create room! please contact our administrator.');
            }
        }
    
        $email = strtolower(Session::get(Config::get('disclosure-types.' . $disclosureTypeHost . '.email'), Config::get('disclosure-types.' . $disclosureType . '.email')));
        if (($email !== '') && ($email === $hosterEmailAddress)) {
            //hoster is already logged in
            //TODO validate attributes again?
            return $this->_join($irmaSessionId);
        }
        $mainContent = view('layout.partials.irma-session-join')->with([
            'irmaSessionId' => $irmaSessionId,
            'disclosureTypeHost' => $disclosureTypeHost,
            'disclosureTypeParticipant' => $disclosureType
        ])->render();
        $buttons = view('layout.partials.buttons')->render();
        return view('layout/mainlayout')->with(
            [
                'message' => $mainContent,
                'title' => __('Choose your role'),
                'buttons' => $buttons
                ]
        );
    }

    public function join_host($irmaSessionId)
    {
        return $this->_join($irmaSessionId);
    }

    public function join_participant($irmaSessionId)
    {
        return $this->_join($irmaSessionId);
    }


    private function _validate($authentications, $visibleName)
    {
        foreach ($authentications as $authentication) {
            $validAttr = true;
            $json = '';
            foreach ($authentication as $attribute) {
                $value = session()->get($attribute, '');
                if ($value == '') {
                    $validAttr = false;
                } else {
                    $json .= sprintf('"%s"&#58; "%s",', $attribute, $value);
                }
            }
            if ($validAttr) {
                break;
            };
        }
        return '{ "name"&#58; "' . $visibleName . '",' . substr($json, 0, -1) . '}';
    }

    private function _get_name($disclosureType)
    {
        $names = Config::get('disclosure-types.' . $disclosureType . '.name');
        $visibleName = '';
        foreach ($names as $name) {
            $validAttr = true;
            foreach ($name as $attribute) {
                $value = session()->get($attribute, '');
                if ($value == '') {
                    $validAttr = false;
                } else {
                    $visibleName .= sprintf("%s ", $value);
                }
            }
            if ($validAttr) {
                break;
            };
        }
        return substr($visibleName, 0, -1);
    }

    private function _send_mail($validatedData, $invitationLink)
    {
        // Send emails in the language in which the site is viewed
        $locale = session()->get('locale', 'nl');

        // If no participant-email is provided, the host needs to invite participants themselves
        $noteForHost = __('Please send the above meeting link to people that you wish to invite to the video meeting.');
        
        // If participant-email is provided, send individual mails to all participants
        for ($i = 1; $i < 7; $i++) {
            if (!empty($validatedData['participant_email_address'. $i])) {
                Mail::to($validatedData['participant_email_address'. $i])->send(new Invitation([
                        'from' => env('MAIL_FROM_ADDRESS'),
                        'reply_to' => $validatedData['hoster_email_address'],
                        'content' => 'emails.' . $locale . '.invitation_' . $validatedData['meeting_type'],
                        'meeting_name' => $validatedData['meeting_name'],
                        'hoster_name' => $validatedData['hoster_name'],
                        'invitation_note' => !empty($validatedData['invitation_note']) ? __('This mail mostly contains general information, but specifically for this meeting the host has added the following personal message: ') . $validatedData['invitation_note'] : '',
                        'invitation_link' => $invitationLink,
                     ]));
                // The host does not have to send the link to the participant themselves
                $noteForHost = __('The meeting link has also been sent to the participant.');
                // if host also provided a note
                if (!empty($validatedData['invitation_note'])) {
                    $noteForHost = __('The meeting link has also been sent to the participant, together with the following personal message: ') . $validatedData['invitation_note'];
                }
            }
        }
        
        // TODO consider multiple participants (change words to plural)

        //Send mail with links to hoster
        Mail::to($validatedData['hoster_email_address'])
                ->send(new Invitation([
                    'reply_to' => env('MAIL_FROM_ADDRESS'),
                    'from' => env('MAIL_FROM_ADDRESS'),
                    'content' => 'emails.' . $locale . '.confirmation_' . $validatedData['meeting_type'],
                    'meeting_name' => $validatedData['meeting_name'],
                    'hoster_name' => $validatedData['hoster_name'],
                    'invitation_note' => $noteForHost,
                    'invitation_link' => $invitationLink,
                ]));
    }

    private function _join($irmaSessionId)
    {
        $irmaSession = \App\IrmaMeetSessions::where('irma_session_id', $irmaSessionId)->first();
        //$hosterName = $irmaSession->hoster_name;
        $hosterEmailAddress = $irmaSession->hoster_email_address;
        $bbbSessionId = $irmaSession->bbb_session_id;
        $meetingType = $irmaSession->meeting_type;
        $disclosureType = Config::get('meeting-types.' . $meetingType . '.irma_disclosure');

        $disclosureTypeHost = Config::get('meeting-types.' . $meetingType . '.irma_disclosure_host', $disclosureType);

        $email = strtolower(Session::get(Config::get('disclosure-types.' . $disclosureTypeHost . '.email'), Config::get('disclosure-types.' . $disclosureType . '.email')));
        if ($email === $hosterEmailAddress) {
            $disclosureType = Config::get('meeting-types.' . $meetingType . '.irma_disclosure_host', $disclosureType);
        }

        $bbb = new BigBlueButton();
        $authentications = Config::get('disclosure-types.' . $disclosureType . '.valid_authentication');
        $validatedName = $this->_get_name($disclosureTypeHost);
        $nameFormat = getenv('IRMA_MEET_NAME_FORMAT');
        if ($nameFormat === 'json') {
            $bbbName = $this->_validate($authentications, $validatedName);
        } else {
            $bbbName = $validatedName;
        }

        if (($email !== '') && ($bbbName !== '') && ($irmaSessionId !== '')) {
            if ($email === $hosterEmailAddress) {
                $password = hash('sha256', 'hoster' . $bbbSessionId);
            } else {
                $password = hash('sha256', 'participant' . $bbbSessionId);
            }
            //redirect to bbb
            $joinParams = new JoinMeetingParameters($bbbSessionId, $bbbName, $password);
            $joinParams->setRedirect(true);
            $joinParams->setJoinViaHtml5(true);
            // Join the meeting by redirecting the user to the generated URL
            $url = $bbb->getJoinMeetingURL($joinParams);
            return \Redirect::to($url);
        } else {
            $mainContent = __('Your attributes could not be verified by IRMA.');
            $buttons = view('layout.partials.buttons')->render();
            return view('layout/mainlayout')->with(
                [
                    'message' => $mainContent,
                    'title' => 'Error',
                    'buttons' => $buttons
                    ]
            );
        }
    }
}
