# pi-wifi

Simple WIFI setup over bluetooth


The BLE server exposes one service `12342233-0000-1000-8000-00805f9b34fb` with two characteristics

1. `0x3344` that supports
   - read - return the connectivity status as enum with values connected, limited, unknown, disconnected
   - write - accept a UTF8 string in the format WIFI:T:WPA;S:<ssid>;P:<password>;H:false; (if ssid or password contains : or ; they must be backslashed eg \;
2. `0x4455` that support read and list the available APs the wifi device found. The response is in the format `interface;SSID;strength\n` a double \n indicates the end of the list