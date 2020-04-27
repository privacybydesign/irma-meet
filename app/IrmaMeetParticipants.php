<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class IrmaMeetParticipants extends Model
{
    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'irma_session_id', 'email_address', 'authentication',
    ];
}
