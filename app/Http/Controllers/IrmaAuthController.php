<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Session;

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
    public function start()
    {
        $irmasession = $this->irma_start_session([
            '@context' => 'https://irma.app/ld/request/disclosure/v2',
            'disclose' => [
                [
                    ['pbdf.pbdf.email.email'],
                ],
                [
                    ['pbdf.gemeente.personalData.fullname'],
                    ['pbdf.pbdf.linkedin.fullname']

                ]
            ],
        ]);
        
        $_SESSION['irmasession'] = $irmasession->token;

        return json_encode($irmasession->sessionPtr);
    }


    /**
     * Show the authenticate screen.
     *
     * @return \Illuminate\Contracts\Support\Renderable
     */
    public function authenticate($url)
    {
        $token = Session::get('irma_session_token', '');
        if ($token === '') {
            $mainContent = view('layout.partials.irma-authenticate')->with([
                'url' => urldecode(urldecode($url))
            ])->render();

            return view('layout/mainlayout')->with(
                [
                'message' => $mainContent,
                'title' => 'Authenticate with IRMA',
                'buttons' => ''
            ]
            );
        } else {
            echo $url;
            return redirect(urldecode(urldecode($url)));
        }
    }
}
