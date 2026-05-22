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
            'message' => '<img id="card-image" class="img-fluid" src="'.url('/') .'/img/yivi-meet.png" alt="Yivi Meet logo v.' . $version . '"><section class="irma-web-form" style="display: none"/>',
            'title' => __('Welcome'),
            'buttons' => view('layout.partials.buttons')->render(),
            'more' => __('more')
        ]);
    }
}
