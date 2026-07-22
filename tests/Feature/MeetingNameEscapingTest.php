<?php

namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Mail;
use Tests\TestCase;

/**
 * Regression test: the user-supplied meeting name must be HTML-escaped (via
 * Laravel's e() helper) before it is rendered into the session-creation success
 * output, so that any markup in the name is shown as text rather than
 * interpreted by the browser. See the private security advisory linked from
 * issue #32 for the detailed rationale.
 */
class MeetingNameEscapingTest extends TestCase
{
    use RefreshDatabase;

    public function testMeetingNameIsHtmlEscapedInSuccessOutput()
    {
        Mail::fake();

        $payload = '<script>alert(1)</script>';

        // _getEmailAddress() reads the disclosed email from the session; the
        // "free" meeting type discloses pbdf.pbdf.email.email.
        $response = $this->withSession(['pbdf.pbdf.email.email' => 'host@example.com'])
            ->post('/irma_session/store', [
                'meeting_name' => $payload,
                'hoster_name' => 'Host',
                'hoster_email_address' => 'host@example.com',
                'meeting_type' => 'free',
                'agreed' => 'on',
            ]);

        $response->assertStatus(200);
        // The raw script tag must never reach the response unescaped ...
        $response->assertDontSee($payload, false);
        // ... it must be HTML-entity encoded instead.
        $response->assertSee('&lt;script&gt;alert(1)&lt;/script&gt;', false);
    }
}
