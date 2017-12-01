go-cld2ch
===================

> **Very Important Note :fire:**
> 
> - go-cld2ch is in early beta. Made for my own needs
> - go-cld2ch is requires golang 1.8+
> - you can must set dbname in configs.yml to use cld2ch, 
> - table CollectD will be created in DB with dbname



How to install
-------------
 - Download [<i class="icon-upload"></i> latest release](https://github.com/riftbit/go-cld2ch/releases) archive for your plaform
 - Unrachive and go to the unarchived path
 - If you are on linux with systemd - run `install.sh`
 - Edit configs in /etc/cld2ch/configs.yml
 - Run `systemctl start cld2ch`
 - Check status `systemctl status cld2ch`


How to configure CollectD
-------------
 - Open configureation file `/etc/collectd/collectd.conf`
 - Set param `Interval` to `60`
 - Enable network plugin if disabled `LoadPlugin network` (delete `#` before)
 - Configure network plugin (examples below)
 - Restart CollectD service
```xml
<Plugin network>
  Server "ip_cld2ch" "port_cld2ch"
</Plugin>
```
For example:
```xml
<Plugin network>
  Server "127.0.0.1" "25826"
</Plugin>
```
