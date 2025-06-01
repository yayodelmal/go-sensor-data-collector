# go-sensor-data-collector
Lightweight Go service that exposes two POST endpoints—one for a temperature sensor and one for a temperature/humidity sensor—to receive JSON readings with timestamps. All incoming data is stored persistently to build and maintain a historical record for each sensor.
