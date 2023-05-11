# BCFMonitor

BlockchainFUE infraestructure monitor.

## Intro

BCFMonitor is a simple service monitor that allows, without any complications, to have a simple email alert when one of the critical services is not responding.

It consists of an executable and a very straightforward configuration file in YAML format.

Installing it is as easy as copying both files to a server and starting the service, for example with a systemd unit, one is included as an example and a Makefile that allows you to install everything and start the service.

## Principles

Check if you are looking for at least two of this principles before to use this software:

- **KISS**: Keep it simple, stupid!
- **SRP**: Single Responsibility Principle.
- **IHNT**: I have no time!
- **FAIF**: Free as in freedom.

## Installation from binary

### Download the latest release for your architecture

Just go to the releases page:

- <https://github.com/jorgefuertes/bcf-monitor/releases>

If you can't find your arquitecture, please, raise an issue.

### Upload the executable to your server

~~~bash
tar xvzf bcf-monitor_linux-amd64_1.2.tar.gz
cd bcf-monitor_linux-amd64_1.2
scp bcf-monitor root@your.server.domain:/usr/local/bin/.
~~~

### Rename `example.yaml` to `bcf-monitor-prod.yaml`

~~~bash
mv example.yaml bcf-monitor-prod.yaml
~~~

### Edit the configuration and upload it

Edit the `yaml` and configure your _runners_, the _smtp_ and the _administrative contacts_.

At this time we are only supporting the monitorization of

- _ping_:
  - Using the system `ping` command. If we use ICMP then we need root privileges.
- _mongodb_
- _redis_
- _web applications_:
  - _GET request_ with or without aditional headers.
  - In addition it looks for a needle text in the _html_.

Upload to `/etc`:

~~~bash
scp bcf-monitor-prod.yml root@your.server.domain:/etc/.
~~~

### Try it

~~~bash
/usr/local/bin/bcf-monitor -c /etc/bcf-monitor-prod.yaml
~~~

### Service it

Now you can define a systemd unit, a daemontools svc or whatever you want to have it alwais running.

### Example systemd unit

Create this file in `/etc/systemd/system/.`:

~~~ini
[Unit]
Description=BCF Monitor
After=network.target auditd.service

[Service]
Type=simple
ExecStart=/usr/local/bin/bcf-monitor -c /etc/bcf-monitor-prod.yaml
Restart=always
User=root
WorkingDirectory=/usr/local/bin

[Install]
WantedBy=multi-user.target
Alias=bcf-monitor.service
~~~

Now you and load the unit and start the service:

~~~bash
systemctl daemon-reload
service bcf-monitor restart
~~~

Watch the log to see if all is working fine:

~~~bash
journalctl -efu bcf-monitor
~~~

## Installation from source

### Clone this repo and change dir

~~~bash
git clone git@github.com:jorgefuertes/bcf-monitor.git
cd bcf-monitor
~~~

### Makefile

Create a `.secrets` file like this one:

~~~ini
SERVER=api.blockchainfue.com
~~~

Take a look into the make file if you need to modify something or add any other architecture.

### Configure the monitor

Copy `example.yaml`

~~~bash
cp conf/example.yaml conf/prod.yaml
~~~

Edit `prod.yaml` and configure your _runners_, the _smtp_ and the _administrative contacts_.

At this time we are only supporting the monitorization of _mongodb_, _redis_ and _web applications_ via _GET request_ with or without aditional headers. In addition it looks for a needle text in the _html_.

#### Try it locally

You may need to raise some tunnels to be allowed to reach the remote services or you can replicate the services on your local machine.

Copy `prod.yaml` to `dev.yaml` and configure it properly. Remember to keep only you as admin in order to don't disturb anyone else.

~~~bash
make run
~~~

### Publish to your server

By default we are building a binary to a GNU/Linux OS, with amd64 architecture. You'll need to adjust the `Makefile` if it doesn't suit your server, but it's pretty easy.

~~~bash
make publish
~~~

## Improvements and issues

- Pull requests are welcome.
- New runners are welcome, always respecting the _KISS_ principle, please.
- Open an issue if you need help or you reach a bug.

## ToDo

- Alerting by telegram.
- Simple authenticated webpage to watch the status in real time.
- New runners:
  - MySQL.
  - Postgres.
  - Memcached.
  - Nginx.
  - SMTP, IMAP.
  - POST in addition to GET.
  - etc...

Feel you free to improve it and make a pull request.

## Author

©2023 Jorge Fuertes Alfranca

- AKA Queru
- <me@jorgefuertes.com>
- <https://jorgefuertes.com>
- <https://github.com/jorgefuertes>

## License

This software is licensed under GNU GPL v3.
You have a copy of this license in this repository, file: `gpl-3.0.txt`, or you can read it online here:

- <https://www.gnu.org/licenses/gpl-3.0-standalone.html>
