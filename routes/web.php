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

Route::get('/irma_auth/start', 'IrmaAuthController@start')->name('irma_auth.start')->middleware('irma_auth');

Route::get('/irma_session/create', 'IrmaSessionController@create')->name('irma_session.create')->middleware('irma_auth');

Route::post('/irma_session/store', 'IrmaSessionController@store')->name('irma_session.store');

Route::get('/irma_session/success', 'IrmaSessionController@success')->name('irma_session.success');

//TODO: Next 2 statements needs irmaauth
Route::get('/irma_session/join/hoster/{irmaSessionId}', 'IrmaSessionController@join_hoster')->name('irma_session.join_hoster')->middleware('irma_auth');

Route::get('/irma_session/join/participant/{irmaSessionId}', 'IrmaSessionController@join_participant')->name('irma_session.join_participant')->middleware('irma_auth');

Auth::routes();

Route::get('/admin', 'AdminController@index')->name('admin');
