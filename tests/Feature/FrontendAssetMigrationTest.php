<?php

namespace Tests\Feature;

use Tests\TestCase;

/**
 * Pins the Node toolchain migration from issue #16: Laravel Mix -> Vite,
 * Bootstrap 4 -> 5, and the removal of jQuery. These are static assertions on
 * the committed sources (blade views, JS entrypoint, package manifest and the
 * public/ tree) so they run without a database or a compiled Vite manifest and
 * would fail the moment a Bootstrap 4 idiom or a jQuery/Mix reference sneaks
 * back in.
 */
class FrontendAssetMigrationTest extends TestCase
{
    private function base(string $relative): string
    {
        return file_get_contents(base_path($relative));
    }

    /**
     * @return string[] every blade template under resources/views
     */
    private function bladeFiles(): array
    {
        $dir = base_path('resources/views');
        $iterator = new \RecursiveIteratorIterator(
            new \RecursiveDirectoryIterator($dir, \FilesystemIterator::SKIP_DOTS)
        );

        $files = [];
        foreach ($iterator as $file) {
            if (str_ends_with($file->getFilename(), '.blade.php')) {
                $files[] = $file->getPathname();
            }
        }

        return $files;
    }

    public function testLayoutsLoadAssetsViaVite()
    {
        foreach (['resources/views/app.blade.php', 'resources/views/layout/app.blade.php'] as $layout) {
            $this->assertStringContainsString('@vite([', $this->base($layout), "$layout should load assets via @vite");
        }

        // The public landing page pulls its head/footer from partials.
        $this->assertStringContainsString('@vite([', $this->base('resources/views/layout/partials/head.blade.php'));
    }

    public function testNoBootstrap4DataAttributesRemainInViews()
    {
        foreach ($this->bladeFiles() as $file) {
            $contents = file_get_contents($file);
            foreach (['data-toggle=', 'data-target=', 'data-dismiss=', 'data-ride=', 'data-parent='] as $bs4Attr) {
                $this->assertStringNotContainsString(
                    $bs4Attr,
                    $contents,
                    sprintf('%s still uses the Bootstrap 4 attribute "%s"; use the "data-bs-*" form.', $file, $bs4Attr)
                );
            }
        }
    }

    public function testNoBootstrap4UtilityRenamesRemainInViews()
    {
        // Utility classes renamed in Bootstrap 5 that appear in this app.
        $renamed = ['ml-auto', 'mr-auto', 'no-gutters', 'dropdown-menu-right', 'class="close"'];

        foreach ($this->bladeFiles() as $file) {
            $contents = file_get_contents($file);
            foreach ($renamed as $needle) {
                $this->assertStringNotContainsString(
                    $needle,
                    $contents,
                    sprintf('%s still uses the Bootstrap 4 class "%s".', $file, $needle)
                );
            }
        }
    }

    public function testViewsNoLongerReferenceRemovedVendorBundles()
    {
        foreach ($this->bladeFiles() as $file) {
            $contents = file_get_contents($file);
            foreach (['jquery-3.4.1', 'popper.min.js', 'bootstrap.min.js', 'bootstrap.min.css'] as $needle) {
                $this->assertStringNotContainsString(
                    $needle,
                    $contents,
                    sprintf('%s still references the removed committed bundle "%s".', $file, $needle)
                );
            }
        }
    }

    public function testRemovedCommittedBundlesAreGone()
    {
        $removed = [
            'public/js/jquery-3.4.1.slim.min.js',
            'public/js/popper.min.js',
            'public/js/bootstrap.min.js',
            'public/css/bootstrap.min.css',
            'public/css/app.css',
            'public/js/app.js',
        ];

        foreach ($removed as $path) {
            $this->assertFileDoesNotExist(base_path($path), "$path should have been removed by the Vite/Bootstrap 5 migration");
        }
    }

    public function testJqueryIsNoLongerADependencyOrImported()
    {
        $this->assertStringNotContainsString('jquery', $this->base('package.json'), 'jQuery must be dropped from package.json');
        $this->assertStringNotContainsString('jquery', $this->base('resources/js/bootstrap.js'), 'bootstrap.js must not import jQuery');
    }

    public function testLaravelMixIsFullyReplacedByVite()
    {
        $this->assertFileDoesNotExist(base_path('webpack.mix.js'), 'webpack.mix.js should be gone after the Vite migration');
        $this->assertFileExists(base_path('vite.config.js'));

        $package = $this->base('package.json');
        $this->assertStringNotContainsString('laravel-mix', $package);
        $this->assertStringContainsString('vite', $package);
    }
}
