<?php

use App\Http\Controllers\AdminController;
use App\Http\Controllers\IrmaAuthController;
use App\Http\Controllers\IrmaSessionController;
use App\Http\Controllers\LocalizationController;
use App\Http\Controllers\MainController;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Route;

Route::get('/', [MainController::class, 'index'])->name('home');

Route::get('/irma_auth/start/{disclosureType}', [IrmaAuthController::class, 'start'])->name('irma_auth.start')->middleware('irma_auth');

Route::get('/irma_session/authenticate/{disclosureType}/{url}', [IrmaAuthController::class, 'authenticate'])->name('irma_session.authenticate');

Route::get('/irma_session/create/{meetingType}', [IrmaSessionController::class, 'create'])->name('irma_session.create')->middleware('irma_auth');

Route::post('/irma_session/store', [IrmaSessionController::class, 'store'])->name('irma_session.store');

Route::get('/irma_session/join/{irmaSessionId}', [IrmaSessionController::class, 'join'])->name('irma_session.join');

Route::get('/irma_session/join_host/{irmaSessionId}', [IrmaSessionController::class, 'join_host'])->name('irma_session.join_host')->middleware('irma_auth');

Route::get('/irma_session/join_participant/{irmaSessionId}', [IrmaSessionController::class, 'join_participant'])->name('irma_session.join_participant')->middleware('irma_auth');

Route::get('lang/{locale}', [LocalizationController::class, 'index']);

Auth::routes();

Route::get('/admin', [AdminController::class, 'index'])->name('admin');
