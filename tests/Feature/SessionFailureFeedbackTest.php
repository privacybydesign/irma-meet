<?php

namespace Tests\Feature;

use Tests\TestCase;

/**
 * The home page must ship descriptive, localized feedback for failed Yivi
 * disclosure sessions so the front-end never has to show the raw status string
 * (e.g. "CANCELLED") to the user. See issue #1.
 */
class SessionFailureFeedbackTest extends TestCase
{
    public function testEnglishFailureMessagesAreInjected()
    {
        $response = $this->withSession(['locale' => 'en'])->get('/');

        $response->assertStatus(200);
        $response->assertSee('window.yiviSessionMessages', false);
        // The cancelled case is the one reported in the issue.
        $response->assertSee('the Yivi session was cancelled', false);
        $response->assertSee('Yivi session timed out', false);
        // The feedback points the user to the documentation.
        $response->assertSee('Read which credentials you need and how to load them', false);
        $response->assertSee('https://privacybydesign.foundation/uitgifte/', false);
    }

    public function testDutchFailureMessagesAreTranslated()
    {
        $response = $this->withSession(['locale' => 'nl'])->get('/');

        $response->assertStatus(200);
        $response->assertSee('de Yivi-sessie is geannuleerd', false);
        $response->assertSee('Lees welke gegevens je nodig hebt', false);
    }
}
