# espressod
Espresso Machine Control - HomeKit and Promtheus enabled

Usage of ./espressod:
  -db string
    	Database path (default "./db")
  -gpio
    	load gpio code (default true)
  -maxTemp float
    	Maximum temperature value (default 130)
  -metrics
    	Enable prometheus metrics (default true)
  -minTemp float
    	Minimum temperature value (default 10)
  -name string
    	Device name (default "Astoria Boiler Thermostat")
  -pin string
    	Homekit Pairing PIN (default "00102003")
  -promPort int
    	Port to reigster /metrics handler on (default 2112)
  -stepTemp float
    	Temperature setting step size (default 0.1)

# How to Use
`espressod` is a standalone golang daemon that creates a Thermostat accessory for use in HomeKit. It also exposes most of the metrics relevant to the operation of the system via a prometheus instrumentation URI. Starting the daemon with no options will setup the accessory and try to control the boiler via a solid state relay. A momentary switch changes the setpoint from min to max and back.

# Why I Wrote This

I purchased a second hand commerical espresso machine a few years ago, and it is far too much work to turn it on and then wait the 30+ minutes for it to warm up fully. I was also concerend with how much power it was drawing. So I in-lined a 40A 250V Solid State Relay which I hooked to a raspberry pi, and wired a PT100 temperature sensor into the boiler. I now have very precise and very easy to automate temperature control of the machine.
