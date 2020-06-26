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
            'message' => '<img class="img-fluid" src="'.url('/') .'/img/IRMA-meet.png" alt="IRMA-meet logo ' . $version . '">',
            'title' => __('welcome'),
            'buttons' => '<button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/default\', \'./irma_session/create/free\');">' . __('Start a free meeting now') . '</button> <button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/medical\', \'./irma_session/create/medical_consult\');">' . __('Start a medical consult now') . '</button><button type="button" class="btn btn-primary btn-lg btn-blue" onclick="startInvitation(\'./irma_auth/start/teacher\', \'./irma_session/create/exam\');">' . __('Start an exam now') . '</button>'
        ]);
    }
}
