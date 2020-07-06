<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

//We use icons from https://github.com/linssen/country-flag-icons

class LocalizationController extends Controller
{
    //
    public function index($locale)
    {
        \App::setLocale($locale);
        //storing the locale in session to get it back in the middleware
        session()->put('locale', $locale);
        return redirect()->home();
    }
}
