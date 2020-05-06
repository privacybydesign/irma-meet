<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateIrmaMeetSessionsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('irma_meet_sessions', function (Blueprint $table) {
            $table->id();
            $table->char('irma_session_id', 48)->unique();
            $table->string('meeting_name');
            $table->string('hoster_name');
            $table->string('hoster_email_address');
            $table->datetime('start_time');
            $table->string('invitation_note')->nullable()->default('');
            $table->string('bbb_session_id');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('irma_meet_sessions');
    }
}
