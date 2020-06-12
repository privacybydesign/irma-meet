<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class IrmaMeetSessions extends Model
{
    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'irma_session_id', 'meeting_name', 'hoster_name', 'hoster_email_address', 'start_time', 'invitation_note', 'bbb_session_id', 'meeting_type'
    ];

    /**
     * The attributes that should be hidden for arrays.
     *
     * @var array
     */
    protected $hidden = [
        
    ];
}
