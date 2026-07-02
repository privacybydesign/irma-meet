<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddBbbPasswordsToIrmaMeetSessionsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * Stores independent random BigBlueButton moderator/attendee passwords per
     * session, replacing the predictable sha256(role . bbb_session_id) scheme.
     * See GHSA-gpgv-24vm-q4vr.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('irma_meet_sessions', function (Blueprint $table) {
            $table->string('bbb_moderator_password')->nullable()->after('bbb_session_id');
            $table->string('bbb_attendee_password')->nullable()->after('bbb_moderator_password');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('irma_meet_sessions', function (Blueprint $table) {
            $table->dropColumn(['bbb_moderator_password', 'bbb_attendee_password']);
        });
    }
}
