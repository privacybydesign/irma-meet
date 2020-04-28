<?php

namespace App\Mail;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Mail\Mailable;
use Illuminate\Queue\SerializesModels;

class Invitation extends Mailable
{
    use Queueable, SerializesModels;

    /**
     * The mailinfo instance.
     *
     * @var Order
     */
    protected $mailinfo;

    /**
     * Create a new message instance.
     *
     * @return void
     */
    public function __construct($mailinfo)
    {
        $this->mailinfo = $mailinfo;
    }

    /**
     * Build the message.
     *
     * @return $this
     */
    public function build()
    {
        return $this->from('no-reply@irma-meet.com')->view('emails.invitation')->with([
            'hoster_name' => $this->mailinfo['hoster_name'],
            'invitation_note' => $this->mailinfo['invitation_note'],
            'invitation_link' => $this->mailinfo['invitation_link']
        ]);
    }
}
