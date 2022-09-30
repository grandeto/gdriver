# Gdriver

Gdriver is intended to automate the actions provided by [gdrive](https://github.com/grandeto/gdrive) lib by watching for changes a given local folder and automatically executing some preconfigured action

*NOTE: Currently only `Upload file to directory (uploadFileToDir)` action is implemented

## Prerequisites

None, binaries are statically linked.

If you want to compile from source you need the [go toolchain](http://golang.org/doc/install).

Version 1.8 or higher.

You need to set all the environment variables present in `.env-example`

## Instalation

### With [Homebrew](http://brew.sh) on Mac

```
brew install gdrive
```

### Compile from source

```bash
go get github.com/grandeto/gdriver
```

The gdrive binary should now be available at `$GOPATH/bin/gdriver`

or

Download `gdriver` from one of the [links in the latest release](https://github.com/grandeto/gdriver/releases).

## Authentication

You need to choose between two authentication methods:

### Service Account (enabled by default)

For server to server communication, where user interaction is not a viable option, 
is it possible to use a service account, as described in this [Google document](https://developers.google.com/identity/protocols/OAuth2ServiceAccount).
If you want to use a service account, instead of being interactively prompted for
authentication, you need to set up `USE_SERVICE_ACCOUNT_AUTH` environment variable to `true`
and `AUTH_SERVICE_ACCOUNT_FILE_NAME` to hold your Service Account file name 
e.g. `AUTH_SERVICE_ACCOUNT_FILE_NAME="gdrive-automated-asdf.json"`.
Then place your Service Account file inside the `.gdrive` folder in your home directory.

### Prompt

You need to set up `USE_SERVICE_ACCOUNT_AUTH` environment variable to `false`

The first time gdriver is launched and takes an upload/sync action 
you will be prompted for a verification code.
The code is obtained by following the instructions printed inside 
`gdrive_auth_value.txt` in your home directory and authenticating with the 
google account for the drive you want access to.
This will create a token file inside the `.gdrive` folder in your home directory.

Note that anyone with access to this file will also have access to your google drive.

If you want to manage multiple drives you can set the environment variable `GDRIVE_CONFIG_DIR`.
Example: `GDRIVE_CONFIG_DIR="/home/user/.gdrive-secondary"`
You will be prompted for a new verification code if the folder does not exist.

## Other

Referrence to [gdrive](https://github.com/grandeto/gdrive) documentation