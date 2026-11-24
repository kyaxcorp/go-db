

## Local Docker ports: 127.0.0.185 only (security)

`docker compose` publishes this project only on **127.0.0.185**, its own loopback address, so it can run next to
other projects on the same machine: open 127.0.0.185:9911. **Never publish a port on 0.0.0.0** (a bare `"8443:443"`
does that): it opens the service to every network the machine is on, and Docker bypasses the firewall.
The rule is at the top of `docker-compose.yml`. On macOS, run `sudo ifconfig lo0 alias 127.0.0.185 up` once after
each boot. After pulling this change, run `docker compose up -d` to move running containers.
