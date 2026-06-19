import _ from 'lodash';
window._ = _;

import * as Popper from '@popperjs/core';
window.Popper = Popper;

import $ from 'jquery';
window.$ = window.jQuery = $;

import axios from 'axios';
window.axios = axios;
window.axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
