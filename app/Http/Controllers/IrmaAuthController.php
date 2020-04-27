<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

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
            ],
        ]);

        $_SESSION['irmasession'] = $irmasession->token;

        return json_encode($irmasession->sessionPtr);
    }


}
