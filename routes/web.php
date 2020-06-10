<?php

use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
 */

Route::get('/', 'MainController@index')->name('home');

Route::get('/irma_auth/start/{meetingType}', 'IrmaAuthController@start')->name('irma_auth.start')->middleware('irma_auth');

Route::get('/irma_session/authenticate/{url}', 'IrmaAuthController@authenticate')->name('irma_session.authenticate');

Route::get('/irma_session/create/{meetingType}', 'IrmaSessionController@create')->name('irma_session.create')->middleware('irma_auth');

Route::post('/irma_session/store', 'IrmaSessionController@store')->name('irma_session.store');

Route::get('/irma_session/join/{irmaSessionId}', 'IrmaSessionController@join')->name('irma_session.join')->middleware('irma_auth');

Route::get('lang/{locale}', 'LocalizationController@index');

Auth::routes();

Route::get('/admin', 'AdminController@index')->name('admin');
