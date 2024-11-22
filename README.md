# vera-volume-manager

This is a simple volume manager for VeraCrypt volumes. It is written in Go and uses the `veracrypt` command line tool to mount and unmount volumes.

## Installation

To install the program, you can use the following command:

```bash
git clone https://github.com/kamuridesu/vera-volume-manager.git
cd vera-volume-manager
go build -ldflags='-s -w -extldflags "-static"' -o veramanager
sudo mv veramanager /usr/local/bin
chmod +x /usr/local/bin/veramanager
```

## Usage

### Config file format

The program uses a Yaml file to store the configuration. The file can be located at any path with no default. The file should have the following format:

```yaml
veracrypt_path: path to the veracrypt executable

bitwarden:
  url: url to the bitwarden server running with bw serve
  password_base64: base64 encoded master password
  credential_name: name of the credential to use

volume:
  folder: path to the folder where the volumes are stored
  name: name of the volume
  mount_point: path to the folder where the volume will be mounted
  size: size of the volume in MB

default_structure:
  - name: name of the folder
    # works with subdirs
    children:
      - name: name of the folder
```

You can refer to the `config.yaml` file for an example.
