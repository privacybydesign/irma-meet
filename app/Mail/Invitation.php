<?php

namespace App\Mail;

use Illuminate\Bus\Queueable;
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
        return $this->from($this->mailinfo['from'])->
            view($this->mailinfo['content'])->
            subject("IRMA-meet link for " . $this->mailinfo['meeting_name'])->
            with([
                'hoster_name' => $this->mailinfo['hoster_name'],
                'invitation_note' => $this->mailinfo['invitation_note'],
                'invitation_link' => $this->mailinfo['invitation_link'],
            ]);
    }
}
