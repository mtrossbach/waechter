# ![Wächter](https://raw.githubusercontent.com/mtrossbach/waechter/main/logo.png)

![License](https://img.shields.io/github/license/mtrossbach/waechter) ![GitHub last commit](https://img.shields.io/github/last-commit/mtrossbach/waechter) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mtrossbach/waechter)

UNDER CONSTRUCTION!

Wächter is a home alarm system. It supports standard smart home devices such as motion sensors, contact sensors, smoke sensors, keypads and sirens provided via Zigbee2Mqqt and Home Assistant. In case of an alarm, SMS can be sent via a connected mobile modem. Support for the WhatsApp Business Cloud API is also planned.

The project is still in an early phase and not all planned features have been implemented yet.

Supported device types via Zigbee2Mqqt:
- Motion sensors
- Contact sensors
- Smoke sensors
- Sirens
- Keypads

The Zigbee2Mqqt integration also supports reading the `tampered` flag which indicates that the device is being opened or dismounted without authorization and leads to an alarm.

Supported device types via Home Assistant:
- Motion sensors
- Contact sensors
- Smoke sensors

For HomeAssistant the `tampered` functionality is current not supported!

## Security Limitations
Because of the way the zigbee protocol works, battery-powered devices only report when they want to transmit information to the network. Otherwise these devices are in sleep mode. It is therefore not possible to determine whether someone is intentionally interrupting the radio contact or blocking it with an interference signal.  

As soon as I have a wired zigbee device like e.g. a stationary siren I will check if it is possible to verify at least with this kind of devices if there is a clean radio connection. If this is possible I will integrate such a test into the system.


## License
[MIT](https://choosealicense.com/licenses/mit/)
