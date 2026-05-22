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
            'message' => view('layout.partials.main', ['version' => $version])->render(),
            'title' => __('Welcome'),
            'buttons' => view('layout.partials.buttons')->render(),
            'more' => __('more')
        ]);
    }
}
