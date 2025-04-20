# TorGate

TorGate is a proxy service that enables access to [Tor Network](https://torproject.org/) hidden services via regular web browsers. 

## Installation

TorGate compiles to a single executable. If you have Go installed, you can install it using `go get github.com/gideongrinberg/tor-gate`, or you can download and build the repo yourself. You can also download the latest release from Github. You will also need to have the [Tor CLI](https://community.torproject.org/onion-services/setup/install/) installed. 

## Usage

To run TorGate, start a Tor SOCKS5 proxy by running the `tor` command. Then, invoke the TorGate binary. The service will listen for incoming requests on port 8080.

### Configuration

You can configure TorGate by creating a `torgate.json` file in the directory in which you are running the program. You can configure the following options:

- `domain`: The domain on which the service is being hosted. This is used when rewriting onion links. (default: `localhost:8080`)
- `port`: The port for the service to listen on (default: `8080`)
- `whitelistOnly`: If this is true, users will only be able to access onion sites that are whitelisted. (default: `false`)
- `whitelist`: A list of whitelisted onion URLs, without the `.onion` at the end (default: `[]`)
- `blacklist`: A list of onion URLs that users are not permitted to access, without the `.onion` at the end (default: `[]`)
- `translations`: Map subdomains to specific onion sites (example: `times: tortimeswqlzti2aqbjoieisne4ubyuoeiiugel2layyudcfrwln76qd` will serve the Tor Times from `times.example.com`)
- `enableTranslations`: Enables the translation feature described above. (default: `false`)
- `showDisclaimer`: If this is enabled, TorGate will display a disclaimer page informing users that it does not provide anonymity and explaining how to download a Tor browser. A single anonymous cookie is used to show the page only on the first visit.

### HTTPS Support

TorGate does not include built-in HTTPS. To serve it securely over HTTPS, use a reverse proxy like [Caddy](https://caddyserver.com/docs/), which handles automatic HTTPS and can host TorGate alongside other applications.

## Credits

This project was inspired by [tor2web](https://github.com/tor2web), which was originally developed by Aaron Swartz and Virgil Griffith and subsequently maintained by Giovanni Pellerano of the GlobaLeaks project. 

## License

Copyright &copy; 2025 Gideon Grinberg

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.