VirtualBox Host Daemon
======================

`vboxhostd` executes commands passed via VirtualBox guest properties on the host side.

| Commands | Description |
|----------|-------------|
| open     | Run `open` (macOS), `xdg_open` (Linux) or `start` (Windows XXX not implemented) with the given value, which is supposed to be an URL (starting by `http:`, `https:` or `mailto:`) |
|          |             |

Usage:

- On host, start `vbhostd`
- On host, run `VBoxManage guestproperty set go-virtualbox vbhostd/open http://www.apple.com`, to open the default browser at http://www.apple.com.
- In Linux/macOS guest, run `sudo VBoxControl guestproperty set vbhostd/open http://www.hp.com`, to open the default browser at http://www.hp.com.
