<?php

namespace Tests\Feature;

use App\Http\Controllers\IrmaSessionController;
use App\IrmaMeetSessions;
use Illuminate\Foundation\Testing\RefreshDatabase;
use ReflectionMethod;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Tests\TestCase;

/**
 * Regression tests for the hardening in GHSA-gpgv-24vm-q4vr:
 * - the session-creation endpoint must sit behind IRMA auth;
 * - the meeting join endpoints must sit behind IRMA auth;
 * - an unknown session id must 404 instead of dereferencing null;
 * - the user-supplied meeting type must be allowlisted before it reaches a
 *   Blade view path;
 * - BigBlueButton room passwords must be independent random values, not
 *   derived from the (stored) session id.
 */
class SessionSecurityTest extends TestCase
{
    use RefreshDatabase;

    private function reflect($method)
    {
        $m = new ReflectionMethod(IrmaSessionController::class, $method);
        $m->setAccessible(true);
        return $m;
    }

    public function testStoreRouteRequiresIrmaAuth()
    {
        $route = app('router')->getRoutes()->getByName('irma_session.store');
        $this->assertContains('irma_auth', $route->gatherMiddleware());
    }

    public function testJoinRoleRoutesRequireIrmaAuth()
    {
        foreach (['irma_session.join_host', 'irma_session.join_participant'] as $name) {
            $route = app('router')->getRoutes()->getByName($name);
            $this->assertContains('irma_auth', $route->gatherMiddleware(), $name);
        }
    }

    public function testUnknownSessionJoinReturns404()
    {
        // No side effects (no BBB room creation) and no null dereference.
        $this->get('/irma_session/join/deadbeefdeadbeefdeadbeefdeadbeef')
            ->assertStatus(404);
    }

    public function testInvalidMeetingTypeIsRejected()
    {
        $this->expectException(NotFoundHttpException::class);
        $this->reflect('_validMeetingType')
            ->invoke(new IrmaSessionController(), 'evil-not-a-type');
    }

    public function testConfiguredMeetingTypeIsAccepted()
    {
        // 'free' is defined in config/meeting-types.php.
        $this->assertSame(
            'free',
            $this->reflect('_validMeetingType')->invoke(new IrmaSessionController(), 'free')
        );
    }

    public function testBbbPasswordsAreRandomAndIndependent()
    {
        $ensure = $this->reflect('_ensureBbbPasswords');

        $sessionA = $this->makeSession('aaaa1111');
        $sessionB = $this->makeSession('bbbb2222');

        [$modA, $attA] = $ensure->invoke(new IrmaSessionController(), $sessionA);
        [$modB, $attB] = $ensure->invoke(new IrmaSessionController(), $sessionB);

        // Moderator and attendee passwords differ from each other...
        $this->assertNotSame($modA, $attA);
        // ...and are not shared across sessions...
        $this->assertNotSame($modA, $modB);
        $this->assertNotSame($attA, $attB);
        // ...and are not derivable from the (stored) bbb_session_id.
        $this->assertNotSame(hash('sha256', 'hoster' . $sessionA->bbb_session_id), $modA);
        $this->assertNotSame(hash('sha256', 'participant' . $sessionA->bbb_session_id), $attA);
        // ...and are persisted so join() and the room creation agree.
        $this->assertSame($modA, $sessionA->fresh()->bbb_moderator_password);
        $this->assertSame($attA, $sessionA->fresh()->bbb_attendee_password);
    }

    private function makeSession($sessionId)
    {
        return IrmaMeetSessions::create([
            'irma_session_id' => $sessionId,
            'meeting_name' => 'Test meeting',
            'hoster_name' => 'Host',
            'hoster_email_address' => 'host@example.com',
            'start_time' => now(),
            'invitation_note' => '',
            'bbb_session_id' => bin2hex(random_bytes(12)),
            'meeting_type' => 'free',
        ]);
    }
}
