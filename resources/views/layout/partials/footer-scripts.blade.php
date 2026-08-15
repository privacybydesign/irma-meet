<!-- Bootstrap's JavaScript is bundled via Vite (loaded from the <head>).
     Placed at the end of the document so the pages load faster. -->

<script src="{{ url('/') }}/js/irma.js?v2&i={{rand()}}" defer></script>
<script type="text/javascript">
  // Human-readable feedback for failed Yivi disclosure sessions, keyed by the
  // status the Yivi frontend rejects its promise with. Translated server-side so
  // custom.js never has to show the raw status string (e.g. "CANCELLED").
  window.yiviSessionMessages = @js([
    'cancelled' => __('The meeting could not be started because the Yivi session was cancelled. This usually means your Yivi app does not yet contain the credentials required for this meeting, or you declined the disclosure request. Please make sure the requested data is loaded in your Yivi app and try again.'),
    'timeout'   => __('The meeting could not be started because the Yivi session timed out. Please try again and approve the disclosure request in your Yivi app.'),
    'error'     => __('The meeting could not be started because the Yivi server could not be reached. Please try again in a moment.'),
    'fallback'  => __('The meeting could not be started because the Yivi session did not complete. Please make sure your Yivi app contains the required credentials and try again.'),
    'docsLabel' => __('Read which credentials you need and how to load them'),
    'docsUrl'   => 'https://privacybydesign.foundation/uitgifte/',
  ]);
</script>
<script src="{{ url('/') }}/js/custom.js?i={{rand()}}"?v2></script>
