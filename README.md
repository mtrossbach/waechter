# ![Wächter](https://raw.githubusercontent.com/mtrossbach/waechter/main/logo.png)

![License](https://img.shields.io/github/license/mtrossbach/waechter) ![GitHub last commit](https://img.shields.io/github/last-commit/mtrossbach/waechter) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mtrossbach/waechter)

UNDER CONSTRUCTION!

Wächter is a home alarm system. It supports standard smart home devices such as motion sensors, contact sensors, smoke sensors, keypads and sirens provided via Zigbee2Mqqt and Home Assistant. In case of an alarm, SMS can be sent via a connected mobile modem. Support for the WhatsApp Business Cloud API is also available.

The project is still in an early phase and not all planned features have been implemented yet.

Currently there are integrations with Zigbee2Mqtt and Home Assistant.

**Supported device types:**

|                           | **Zigbee2Mqtt** | **Home Assistant** |
|:-------------------------:|:---------------:|:------------------:|
| **Motion sensor**         |:white_check_mark:|:white_check_mark:|
| **Contact/window sensor** |:white_check_mark:|:white_check_mark:|
| **Smoke sensor**          |:white_check_mark:|:white_check_mark:|
| **Siren**                 |:white_check_mark:| :x:                  |
| **Keypad**                |:white_check_mark:| :x:                  |

To increase security, the alarm system can respond to tampering, poor radio link quality (for wireless devices) and running low batteries.

**Supported device state attributes:**

|                           | **Zigbee2Mqtt** | **Home Assistant** |
|:-------------------------:|:---------------:|:------------------:|
| **`tamper` flag**         |:white_check_mark:|:white_check_mark: <br />(not working with zigbee devices and ZHA)|
| **Link quality** |:white_check_mark:|:x:|
| **Battery**          |:white_check_mark:|:white_check_mark:|

In the event of an alarm, a notification can be sent.

**Currently supported notification channels:**

- :white_check_mark: WhatsApp Business Cloud API

## Security Limitations
Because of the way the zigbee protocol works, battery-powered devices only report when they want to transmit information to the network. Otherwise these devices are in sleep mode. It is therefore not possible to determine whether someone is intentionally interrupting the radio contact or blocking it with an interference signal.  

## TODO (not implemented yet)
- Siren support via Home Assistant
- Home Assistant link quality
- SMS sending
- More configuration options
