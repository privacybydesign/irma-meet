<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class MainController extends Controller
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
    public function index()
    {
        $version = \Tremby\LaravelGitVersion\GitVersionHelper::getVersion();
        return view('layout/mainlayout')->with([
            'message' => '<img class="img-fluid" src="'.url('/') .'/img/IRMA-meet.png" alt="IRMA-meet logo v.' . $version . '">',
            'title' => __('welcome'),
            'buttons' => '<button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/default\', \'./irma_session/create/free\');">'. '<img class="img-fluid" src="'.url('/') .'/img/team-icon.svg">' . __('Start een video gesprek') . '</button>' . '<br>Iedereen met een <a href="https://privacybydesign.foundation/uitgifte/email/">e-mail</a> en naam (<a href="https://services.nijmegen.nl/irma/gemeente/start">BRP</a> of <a href="https://privacybydesign.foundation/uitgifte/social/linkedin/">LinkedIn</a>) in IRMA kan gratis een gesprek starten<br><br>' . '<button type="button" class="btn btn-secondary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/medical\', \'./irma_session/create/medical_consult\');">' . '<img class="img-fluid" src="'.url('/') .'/img/medical_icon.svg">' . __('Start een medisch video consult') . '</button>' . '<br>Alleen een zorgverlener met een <a href="https://irma-agb.nuts.nl/">AGB-code</a> in IRMA kan een medisch gesprek starten<br><br>' . '<button type="button" class="btn btn-secondary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/teacher\', \'./irma_session/create/exam\');">' . '<img class="img-fluid" src="'.url('/') .'/img/exam-icon.svg">' . __('Start een mondeling video tentamen') . '</button>' . '<br>Alleen een docent met een <a href="https://privacybydesign.foundation/uitgifte/surfnet/surfnet/">SURF</a> kaartje in IRMA kan een tentamen starten',
            'more' => __('more')
        ]);
    }
}
