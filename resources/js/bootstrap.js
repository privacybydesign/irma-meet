// Import Bootstrap's JS bundle (includes Popper). Exposed on window so the
// pre-existing data-attribute auto-init markers (data-bs-toggle, data-bs-target,
// data-bs-dismiss) continue to work.
import * as bootstrap from 'bootstrap';

window.bootstrap = bootstrap;
