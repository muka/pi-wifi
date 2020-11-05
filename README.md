# pi-wifi

Simple WIFI setup over bluetooth

This servie uses network manager API to connect to a WIFI network.

The implementation offers two interfaces

- a bluetooth service
- a rest API for local(host) usage by other services

## Example

You can use `go run main.go` on a linux machine to start the software.

Then point your browser to http://localhost:9099/ to see an example usage.

## Implementation 

### Bluetooth

The BLE server exposes one service `12342233-0000-1000-8000-00805f9b34fb` with two characteristics

1. `0x3344` that supports
   - read - return the connectivity status as enum with values connected, limited, unknown, disconnected
   - write - accept a UTF8 string in the format WIFI:T:WPA;S:<ssid>;P:<password>;H:false; (if ssid or password contains : or ; they must be backslashed eg \;
2. `0x4455` that support read and list the available APs the wifi device found. The response is in the format `SSID;strength\n` a double \n indicates the end of the list

### HTTP API

The service exposes also an HTTP API to intereact with WIFI connections

- `/connect` connect to a WIFI connection. Expects a body in the format `{ "payload": "WIFI:T:WPA;S:your ssid;P:your password;H:false;;" }`
- `/status` return the connection status with format `{"status": "connected"}`
- `/listap` list the reachable APs in format `{"accessPoints": [ { "ssid": "example", "strength": 54 } ]}`


The connection string format is based on https://github.com/zxing/zxing/wiki/Barcode-Contents#wi-fi-network-config-android-ios-11