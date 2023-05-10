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

## Usage

### Clone this repo and change dir

~~~bash
git clone git@github.com:jorgefuertes/bcf-monitor.git
cd bcf-monitor
~~~

### Copy Makefile.dist

~~~bash
cp Makefile.dist Makefile
~~~

### Edit the Makefile to fit your needs

Take a look specifically to this var:

~~~make
SERVER = api.blockchainfue.com
~~~

### Configure the monitor

Copy `conf.dist.yaml`

~~~bash
cp conf/conf.dist.yaml conf/prod.yaml
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
