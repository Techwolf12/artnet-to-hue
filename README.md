# Artnet-To-Hue
artnet-to-hue is a bridge between Art-Net and Philips Hue. 
It allows you to control Philips Hue lights in an entertainment zone using Art-Net, which is commonly used in lighting control systems.
It only has support for color lights and entertainment zones, since an entertainment zone is a maximum of 10 lights, you can use it to control a maximum of 10 lights at once.

## Why
In my house full of Philips Hue lights, I wanted to be able to control some with a proper light setup during a party.

## Installation

macOS users can install `artnet-to-hue` using Homebrew Tap:

```bash
brew tap techwolf12/tap
brew install techwolf12/tap/artnet-to-hue
```

For Docker users, you can use the Docker image:

```bash
docker run --network host ghcr.io/techwolf12/artnet-to-hue:latest server -i <ip-address> 
```

For other systems, see the [releases page](https://github.com/Techwolf12/artnet-to-hue/releases/).

## Usage
Be sure to create an entertainment zone in your Philips Hue app before using this tool.
First you can discover your Hue Bridge by running:

```bash
artnet-to-hue discover
```
This will output the IP address of your Hue Bridge along with the command to pair.
Next, you can pair your Hue Bridge by running the command provided in the previous step:

```bash
artnet-to-hue pair -i <ip-address>
```
Be sure to save the username and client key generated after pairing, as you will need them to control your lights.

After pairing, you can run bridgeInfo to see the entertainment zones available:
```bash
artnet-to-hue bridgeInfo -i <ip-address> -u <username>
```
Finally, you can start the server to listen for Art-Net packets and control your Hue lights:

```bash
artnet-to-hue server -i <ip-address> -u <username> -c <client-key> -e <entertainment-zone> -l <amount-of-lights>
```

Be sure to use help to see all available options.

## Options

## `artnet-to-hue server` Flags

| Flag | Shorthand | Type     | Default | Description                                                  |
|------|-----------|----------|---------|--------------------------------------------------------------|
| `--hue-bridge-ip` | `-i` | IP Address | *none*  | IP address of the Hue bridge                                 |
| `--username`      | `-u` | String     | *none*  | Username for the Hue bridge                                  |
| `--client-key`    | `-c` | String     | *none*  | Client key for the Hue bridge (used for DTLS authentication) |
| `--entertainment-zone` | `-e` | String | *none*  | Entertainment zone ID for the Hue bridge                     |
| `--lights`        | `-l` | Integer    | `10`    | Number of lights in the entertainment zone                   |
| `--artnet-universe` | `-n` | UInt16   | `0`     | Art-Net universe to listen on                                |
| `--artnet-dmx-start` | `-a` | Integer | `1`     | Art-Net DMX start channel                                    |
| `--debug`         | `-d` | Boolean    | `false` | Debug logging )                                              |

---

## `artnet-to-hue pair` Flags

| Flag | Shorthand | Type      | Default | Description |
|------|-----------|-----------|---------|-------------|
| `--hue-bridge-ip` | `-i` | IP Address | *none*  | IP address of the Hue bridge |

---

## `artnet-to-hue bridgeInfo` Flags

| Flag | Shorthand | Type      | Default | Description |
|------|-----------|-----------|---------|-------------|
| `--hue-bridge-ip` | `-i` | IP Address | *none*  | IP address of the Hue bridge |
| `--username`      | `-u` | String     | *none*  | Username for the Hue bridge |

## Contributing
If you want to contribute to this project, feel free to open an issue or a pull request.
You can also help by reporting bugs or suggesting features.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.