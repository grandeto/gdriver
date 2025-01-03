# Gdriver

Gdriver is intended to automate the actions provided by [gdrive](https://github.com/grandeto/gdrive) lib by watching for changes a given local folder and automatically executing some preconfigured actions

*NOTE: Currently only `Upload file to directory (uploadFileToDir)` action is implemented

## Prerequisites

Go version 1.8 or higher.

Create an `.env` and set all the environment variables from `.env-example`

## Instalation

### Install binary

```bash
go install github.com/grandeto/gdriver@latest
```

The gdriver binary should now be available at `$GOPATH/bin/gdriver`

### Compile binary from source

```bash
git clone https://github.com/grandeto/gdriver && cd gdriver && go build .
```

## Authentication

You need to choose between two authentication methods:

### Service Account

For server to server communication, where user interaction is not a viable option, 
is it possible to use a service account, as described in this [Google document](https://developers.google.com/identity/protocols/OAuth2ServiceAccount).
If you want to use a service account, instead of being interactively prompted for
authentication, you need to set up `SERVICE_ACCOUNT_AUTH` environment variable to `true`
and `AUTH_SERVICE_ACCOUNT_FILE_NAME` to hold your Service Account file name 
e.g. `AUTH_SERVICE_ACCOUNT_FILE_NAME="gdrive-automated-asdf.json"`.
Then place your Service Account file inside your `GDRIVE_CONFIG_DIR` (deafults to `~/.gdrive`).

Note that anyone with access to this file will also have access to your google drive.

### Prompt

You need to set up `SERVICE_ACCOUNT_AUTH` environment variable to `false`

The first time gdriver is launched and takes an upload/sync action 
you will be prompted for a verification code.
The code is obtained by following the instructions printed inside 
`gdrive_auth_value.txt` in your home directory and authenticating with the 
google account for the drive you want access to.
This will create a token file inside your `GDRIVE_CONFIG_DIR` (deafults to `~/.gdrive`).

Note that anyone with access to this file will also have access to your google drive.

### Note

If you want to manage multiple drives you can set different environment variable `GDRIVE_CONFIG_DIR` for each client binary you build.
Example: `GDRIVE_CONFIG_DIR="/home/user/.gdrive-secondary"`

## Initialization

Run compiled binary as a daemon or through supervisor

Make sure you have set up an `.env` following the `.env-example` in e.g. `/home/<user>/.gdrive/.env`

### Ubuntu - Example systemd user config

- create `~/.config/systemd/user/gdriver.service`

```
[Unit]
Description=Gdriver
Wants=network-online.target
RequiresMountsFor=/home
After=network.target network-online.target

[Service]
EnvironmentFile=/home/<user>/.gdrive/.env
WorkingDirectory=/home/<user>/go/bin
ExecStart=/home/<user>/go/bin/gdriver
Restart=always

[Install]
WantedBy=default.target network-online.target
```

- start gdriver.service

```bash
systemctl --user daemon-reload
systemctl --user start gdriver
systemctl --user enable gdriver
systemctl --user status gdriver
```

### Mac OS - Example plist user config

- create `~/Library/LaunchAgents/com.gdriver.plist`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.gdriver.app</string>
	<key>WorkingDirectory</key>
	<string>path to dir containing .env</string>
	<key>Program</key>
	<string>path to griver executable</string>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/tmp/gdriver.log</string>
	<key>StandardErrorPath</key>
	<string>/tmp/gdriver.error.log</string>
</dict>
</plist>
```

- start com.gdriver.plist

	`launchctl load -w ~/Library/LaunchAgents/com.gdriver.plist`

## Other

Referrence to [gdrive](https://github.com/grandeto/gdrive) documentation
