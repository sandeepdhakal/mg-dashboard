A dashboard for MG vehicles. Currently only provides a TUI.\
_Note_: So far I've tested this only on my personal car, and on a _Linux_ machine.

## Requirements

The dashboard connects to a MQTT server that provides the data. The [SAIC Python MQTT Gateway provides](https://github.com/SAIC-iSmart-API/saic-python-mqtt-gateway) provides exactly that.

Before that, the first step is to have a MQTT broker up and running. I installed [Eclipse Mosquitto](https://github.com/eclipse-mosquitto/mosquitto) on my machine with

```bash
sudo pacman -S mosquitto
```

Once installed, the mosquitto systemd service should also be started.

```bash
sudo systemctl start mosquitto.service
```

### SAIC account

You also need a SAIC account, which you can create in the MG iSmart app on your phone. The MQTT service will use your username/password to poll the data from MG's servers before publishing them.

## Configuring the saic mqtt gateway

`saic-python-mqtt-gateway` can read environment variables to configure the broker. In my case, I have put them in a .envrc file, which is then read by [direnv](https://direnv.net/) to make the environment variables available for the project.

Clone the repo from GitHub and install the dependencies as usual. In my case I used **uv**:

```bash
uv venv # for a virtual environment
uv add -r requirements.txt # this will also add the requirements to pyproject.toml
```

Now, inside the directory where I cloned the `saic-python-mqtt-gateway` project, I have the following `.envrc` file.

```bash
export VIRTUAL_ENV=".venv"
layout python python3

export SAIC_USER="<your-saic-username>"
export SAIC_PASSWORD="<your-password>"

export SAIC_REST_URI="https://gateway-mg-au.soimt.com/api.app/v1/"
export SAIC_REGION="au"
export MQTT_URI="tcp://localhost:1883"
export MQTT_USERNAME="mqtt_user"
export MQTT_PASSWORD="secret"
export MQTT_LOG_LEVEL="DEBUG"
export HA_DISCOVERY_ENABLED=false
```

The `SAIC_REST_URI`, `SAIC_REGION` are for AU/NZ only. Check the documentation [here](https://github.com/SAIC-iSmart-API/saic-python-mqtt-gateway) for configuration for other regions.

Finally we can start the service with:

```bash
python ./mqtt_gateway.py
```

## Running the dashboard

### Environment variables

The environment variables used by the MQTT gateway are also used by the dashboard. I have set these in my `~/.zshenv` file:

```bash
...

# MG dashboard
export SAIC_USER="<your_saic_username>"
export SAIC_BROKER_URI="tcp://localhost"
export SAIC_BROKER_PORT="1883"
export SAIC_MQTT_USER="mqtt_user"
export SAIC_MQTT_PASS="secret"

...
```

Now, we're ready to start our dashboard:

```bash
go run ./tui.go
```
