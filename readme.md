# CloudServer VPN

Run a New Zealand based VPN on demand for only [1.5c per hour](https://voyager.nz/business/hosting/virtual-servers#balanced).

## Prerequisites

Some initial set up is required.

### Cloud Server

This system utilises [Voyager VPS Manager](https://voyager.nz/business/hosting/virtual-servers) to create servers to run WireGuard. You must sign up to use VPS Manager and [generate an API token](https://cloudserver.nz/account#api-tokens).

### Wireguard

You will need to have a [private/public key pair](https://www.wireguard.com/quickstart/#key-generation) for both the server and at least one peer.

You will also need to decide what [private network range](https://datatracker.ietf.org/doc/html/rfc1918#section-3) WireGuard should use. This must be unique to WireGuard and not conflict with your home private network range. _Most_ home networks use IPs in `192.168.0.0/23` or `10.0.0.0/23` ranges by default, so these should probably be avoided. Picking something completely random is good, e.g. `10.194.89.0/24`. The range only needs to be big enough to support your server and any peers. Two peers (and the server) use three IPs total, so `/30` would be sufficient.

## Usage

The app accepts three command line arguments, which have different configuration requirements. Configuration is read through environment variables.

### Create VPN

A VPN can be created by using `cloudserver-vpn --create`

| Key | Description | Mandatory |
|-----|-------------|-----------|
| CLOUDFLARE_APIKEY | [Scoped API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) to update a Cloudflare DNS record | N |
| CLOUDFLARE_ZONE | The name of the Cloudflare zone to update, e.g. example.com | N |
| CLOUDSERVER_APIKEY | [API token](https://cloudserver.nz/account#api-tokens) for cloudserver.nz | Y |
| CLOUDSERVER_PROJECT | The ID of the [Cloud Server project](https://cloudserver.nz/projects) where the server will be provisioned<sup>1</sup> | N |
| SERVER_NAME | The name for this server, must be [a valid RFC 3696 subdomain](https://datatracker.ietf.org/doc/html/rfc3696) | Y |
| WIREGUARD_ADDRESS | The IPv4 CIDR to use for the WireGuard interface, must include a big enough subnet to accomodate all peers, e.g. `10.194.89.1/24` | Y |
| WIREGUARD_LISTENPORT | The port for WireGuard to listen on, if not specified, `51820` will be used | N |
| WIREGUARD_PEER#\_ALLOWEDIPS | The IPv4 CIDR to allow connections for peer #<sup>2</sup>, e.g. `10.194.89.2/32` | Y |
| WIREGUARD_PEER#\_PUBLICKEY | The public key for the associated peer | Y |
| WIREGUARD_PRIVATEKEY | The private key for the WireGuard server | Y |

<sup>1</sup> If a project is not specified a new project called `VPNs` will be created. This project **must** only contain VPN servers as all servers will be removed when `--remove` is called.
<br/>
<sup>2</sup> You can add up to 255 peers as long as the pair of `ALLOWEDIPS` and `PUBLICKEY` are both specified. Peer prefixes range from `WIREGUARD_PEER0_...` to `WIREGUARD_PEER254_...`.</sub>

### Remove all VPNs

All VPNs can be removed by using `cloudserver-vpn --remove`

This command will remove **all** servers in the Cloud Server project.

| Key | Description | Mandatory |
|-----|-------------|-----------|
| CLOUDSERVER_APIKEY | [API token](https://cloudserver.nz/account#api-tokens) for cloudserver.nz | Y |
| CLOUDSERVER_PROJECT | The ID of the [Cloud Server project](https://cloudserver.nz/projects) where the server is provisioned | N |

### Remove a single VPN

A single VPN can be removed by using `cloudserver-vpn --remove <server id>`

| Key | Description | Mandatory |
|-----|-------------|-----------|
| CLOUDSERVER_APIKEY | [API token](https://cloudserver.nz/account#api-tokens) for cloudserver.nz | Y |

### Run as HTTP server

An HTTP server can be run by using `cloudserver-vpn --serve`

The HTTP server requires the same configuration as [Create VPN](#create-vpn) with an additional configuration for http port.

| Key | Description | Mandatory |
|-----|-------------|-----------|
| HTTP_PORT | Port to listen for HTTP connections on, if not specified `5252` will be used | N |
