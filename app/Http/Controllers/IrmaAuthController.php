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
        $target = urldecode(urldecode($url));

        // Only ever continue to a location inside this application. The return
        // URL is always built server-side as URL::to('/') . '/' . path, so any
        // value that does not point back at our own origin is rejected.
        $base = url('/');
        if ($target !== $base && ! str_starts_with($target, $base . '/')) {
            abort(400);
        }

        // The disclosure type must be one of the configured types.
        if (Config::get('disclosure-types.' . $disclosureType) === null) {
            abort(400);
        }

        $token = session()->get('irma_session_token', '');
        if ($token === '') {
            $mainContent = view('layout.partials.irma-authenticate')->with([
                'url' => $target,
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
            return redirect($target);
        }
    }
}
