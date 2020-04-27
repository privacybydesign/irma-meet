<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
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
    public function create()
    {
        $validated_email = Session::get('validated_email', '');
        return view('layout/irma-session')->with('validated_email', $validated_email);
    }

    /**
     * Store a session.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function store(Request $request)
    {
        //@TODO: verify and store all
        $validatedData = $request->validate([
            'meeting_name' => 'required|max:255',
            'hoster_name' => 'required',
            'hoster_email_address' => 'required',
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
        $validatedData = array_merge($validatedData, [ 'irma_session_id' => $uniqueId, 'start_time' => now() ]);
        $irma_session = \App\IrmaMeetSessions::create($validatedData);

       // TODO verify emial address with session
       //for all participants store data
        for ($i = 1; $i < 7; $i++) {
            if ( $validatedData['participant_email_address'.$i] ) {
                $irmaParticipant=[
                    'irma_session_id' => $irma_session['irma_session_id'], 
                    'email_address' => $validatedData['participant_email_address'.$i], 
                    'authentication' => 1, 
                ];
                $irma_participants = \App\IrmaMeetParticipants::create($irmaParticipant);
            }
        } 
        
        return redirect()->route('irma_session.success');
    }

    /**
     * Show the store session screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function success()
    {
        $mainContent = 'Meeting is successfully validated and data has been saved';
        return view('layout/irma-session-success')->with('mainContent', $mainContent);
    }

}
