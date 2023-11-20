# MHF Dev Proxy

Proxy that keeps your game open even after the server restarts. Should only be used for development.

## Config

Run with `--help` to see configurable options.

## Usage

1. Set `DevModeOptions.ProxyPort` to `8090` in your Erupe config
2. Run proxy
3. Run game

You should now be able to restart the server while the game is still running.

The proxy can also handle multiple clients, and if the game restarts the connection to the server is restarted as well.
