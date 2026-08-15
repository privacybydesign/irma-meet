<?php

namespace Tests\Feature;

use App\Http\Controllers\IrmaSessionController;
use ReflectionMethod;
use Session;
use Tests\TestCase;

/**
 * IrmaSessionController::_getEmailAddress() reads the disclosed email attribute
 * from the session for a given disclosure type. It must lowercase the value
 * when present and return an empty string (never null) when the attribute is
 * absent, so that the PHP 8.1+ "passing null to strtolower()" deprecation can
 * never fire. See issue #26.
 */
class EmailAddressDisclosureTest extends TestCase
{
    private function callGetEmailAddress($disclosureType)
    {
        $method = new ReflectionMethod(IrmaSessionController::class, '_getEmailAddress');
        $method->setAccessible(true);

        return $method->invoke(new IrmaSessionController(), $disclosureType);
    }

    public function testReturnsLowercasedEmailWhenAttributePresent()
    {
        // 'pbdf.pbdf.email.email' is one of the configured email fields for the
        // 'default' disclosure type (see config/disclosure-types.php).
        Session::put('pbdf.pbdf.email.email', 'Alice@Example.COM');

        $this->assertSame('alice@example.com', $this->callGetEmailAddress('default'));
    }

    public function testReturnsEmptyStringWithoutDeprecationWhenAttributeAbsent()
    {
        // No email attribute in the session at all. The helper must hand
        // strtolower() an empty string rather than null: strtolower(null) emits
        // a "passing null to parameter #1" deprecation on PHP 8.1+, which
        // PHPUnit's error handler turns into a failed test.
        $email = $this->callGetEmailAddress('default');

        $this->assertSame('', $email);
    }
}
