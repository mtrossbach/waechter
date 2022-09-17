# ![Wächter](https://raw.githubusercontent.com/mtrossbach/waechter/main/logo.png)

![License](https://img.shields.io/github/license/mtrossbach/waechter) ![GitHub last commit](https://img.shields.io/github/last-commit/mtrossbach/waechter) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mtrossbach/waechter)

UNDER CONSTRUCTION!

Wächter is a home alarm system. It supports standard smart home devices such as motion sensors, contact sensors, smoke sensors, keypads and sirens provided via Zigbee2Mqqt and Home Assistant. In case of an alarm, SMS can be sent via a connected mobile modem. Support for the WhatsApp Business Cloud API is also available.

The project is still in an early phase and not all planned features have been implemented yet.

Supported device types via Zigbee2Mqtt:
- Motion sensors
- Contact sensors
- Smoke sensors
- Sirens
- Keypads

The Zigbee2Mqqt integration also supports reading the `tamper` flag which indicates that the device is being opened or dismounted without authorization and leads to an alarm.

Supported device types via Home Assistant:
- Motion sensors
- Contact sensors
- Smoke sensors

For HomeAssistant the `tamper` functionality is current not supported via ZHA integration. If you have integrated zigbee devices via Zigbee2Mqtt it should work.

## Security Limitations
Because of the way the zigbee protocol works, battery-powered devices only report when they want to transmit information to the network. Otherwise these devices are in sleep mode. It is therefore not possible to determine whether someone is intentionally interrupting the radio contact or blocking it with an interference signal.  

## TODO (not implemented yet)
- Siren support via Home Assistant
- SMS sending
- More configuration options
