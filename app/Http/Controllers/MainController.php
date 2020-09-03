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
            'message' => '<img id="card-image" class="img-fluid" src="'.url('/') .'/img/IRMA-meet.png" alt="IRMA-meet logo v.' . $version . '"><section class="irma-web-form" style="display: none"/>',
            'title' => __('Welcome'),
            'buttons' => '<button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/default\', \'./irma_session/create/free\');">'. '<img class="img-fluid" src="'.url('/') .'/img/team-icon.svg">' . __('Create video meeting') . '</button>' . '<br>' .  __('<p style="font-size: 14px; color:#1B4D8C;"><br>For use in the Netherlands:</p>') . '<button type="button" class="btn btn-secondary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/teacher\', \'./irma_session/create/exam\');">' . '<img class="img-fluid" src="'.url('/') .'/img/exam-icon.svg">' . __('Start oral video exam') . '</button>' .  '<br><br>' . '<button type="button" class="btn btn-secondary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/medical\', \'./irma_session/create/medical_consult\');">' . '<img class="img-fluid" src="'.url('/') .'/img/medical_icon.svg">' . __('Start medical video consult') . '</button>',
            'more' => __('more')
        ]);
    }
}
