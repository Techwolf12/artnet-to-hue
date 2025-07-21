# Artnet-To-Hue
artnet-to-hue is a bridge between Art-Net and Philips Hue. It allows you to control Philips Hue lights in an entertainment zone using Art-Net, which is commonly used in lighting control systems.

## Why
In my house full of Philips Hue lights, I wanted to be able to control them with a proper light setup during a party.

## Installation

macOS users can install `artnet-to-hue` using Homebrew Tap:

```bash
brew tap techwolf12/tap
brew install techwolf12/tap/artnet-to-hue
```

For Docker users, you can use the Docker image:

```bash
docker run -v "$(pwd)":/mnt ghcr.io/techwolf12/artnet-to-hue:latest server -
```

For other systems, see the [releases page](https://github.com/Techwolf12/artnet-to-hue/releases/).

## Usage
First you can discover your Hue Bridge by running:

```bash
artnet-to-hue discover
```
This will output the IP address of your Hue Bridge along with the command to pair.
Next, you can pair your Hue Bridge by running the command provided in the previous step:

```bash
artnet-to-hue pair -i <ip-address>
```
After pairing, it shows you the command to run the server:

```bash
artnet-to-hue server -i <ip-address>
```

Be sure to use help to see all available options.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.