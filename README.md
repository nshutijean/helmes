# Helmes

Helmes is a demo 12 factor app in Go that can be deployed on any cloud native compliant platform as a container.
Helmes as the name obviously implies is a messenger that Sends SMS messages to any supported carrier in Rwanda.

## Goals

- [x] Simple sms delivery.
- [x] Delivery notifications.
- [ ] Implement a simple store(badger) to record access records.
- [ ] Document authentication via [pomerium.io](https://www.pomerium.io/).
- [ ] Add application metrics.
- [ ] Build and deploy container imag.
- [ ] Add valid knative deployment manifests

## Environment variables
To function Helmes requires a couple of environment variables:

```bash
cat .template.env
PORT=8080
HELMES_SMS_APP_ID="fdi sms app id"
HELMES_SMS_APP_SECRET="fdi sms app password"
HELMES_SENDER_IDENTITY="fdi sms sender id"
HELMES_CALLBACK_URL="delivery(dlr) report callback url" # optional
```

Besides the port the other variables can be obtained by subscribing to https://www.fdibiz.com/ messaging API.

## Try helmes
To start Helmes on your laptop:

1. git clone this repository

2. source the environment variables
```bash
cp .template.env .env
# edit the file .env with variables and credentials the source the file
source .env

```

3. build the helmes binary
```bash
CGO_ENABLED=0 go build -o bin/helmes ./cmd/helmes
# view the output binary
ls bin
```

For convience your could install [task](https://taskfile.dev/) a make alternative then:
```
# it will build your binary and start the helmes server
task run 
```

4. check helmes version via `/api/version`

```bash
# source .env

curl localhost:$PORT/api/version
```
5. send an sms messsage via `/api/send`

````bash
# source .env

# replace with your 078xxxxxxx with your number
export PHONE="your phone"
export MESSAGE="your message"

cat example.json | jq --arg PHONE $PHONE '.recipient=$PHONE' | tee helmes.json
cat helmes.json | jq --arg MESSAGE $MESSAGE '.payload=$MESSAGE' | tee helmes.json
 ````

Finally send the payload as defined in `helmes.json`

```bash
curl -d "@helmes.json" -H "Content-Type: application/json" -X POST localhost:$PORT/api/send
```

There is an notifications endpoint `api/events/$ID/status` you could subscribe to to recieve
sms delivery notications. Therss is a helping script you could use to run an example:

```bash
./scripts/send.sh
```

6. You can build a docker image
```bash
task image
```
