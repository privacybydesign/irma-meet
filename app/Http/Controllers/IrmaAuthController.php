<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Session;
use Config;

class IrmaAuthController extends Controller
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
     * Show the main screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    #public function start($meetingType)
    #{
    #    print_r(Config::get('meeting-types.' . $meetingType . '.irma_disclosure'));
    #    $irmasession = $this->irma_start_session(Config::get('meeting-types.' . sprintf("%s", $meetingType) . '.irma_disclosure'));
    #
    #    $_SESSION['irmasession'] = $irmasession->token;

    #    return json_encode($irmasession->sessionPtr);
    #}


    /**
     * Show the authenticate screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function authenticate($disclosureType, $url)
    {
        $token = Session::get('irma_session_token', '');
        if ($token === '') {
            $mainContent = view('layout.partials.irma-authenticate')->with([
                'url' => urldecode(urldecode($url)),
                'disclosureType' => $disclosureType
            ])->render();
            return view('layout/mainlayout')->with(
                [
                'message' => $mainContent,
                'title' => __('Authenticate with IRMA'),
                'buttons' => ''
            ]
            );
        } else {
            echo $url;
            return redirect(urldecode(urldecode($url)));
        }
    }
}
