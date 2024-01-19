import './styles/main.css';
import 'htmx.org';

window.htmx = require('htmx.org');

function component() {
    const element = document.createElement('div');

    // Lodash, currently included via a script, is required for this line to work
    element.innerHTML = "Hello, Webpack!";

    return element;
}

document.body.appendChild(component());
