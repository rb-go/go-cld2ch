#!/usr/bin/env bash
mkdir -p /etc/cld2ch
cp ./config.yml /etc/cld2ch/config.yml
cp ./cld2ch /etc/cld2ch/cld2ch
mv cld2ch.service /lib/systemd/system/.
chmod 755 /lib/systemd/system/cld2ch.service
chmod +x /etc/cld2ch/cld2ch
systemctl daemon-reload
systemctl enable cld2ch.service
echo "Set settings in /etc/cld2ch/config.yml"
echo "After that run service: [systemctl start cld2ch]"
#systemctl start cld2ch
#journalctl -f -u cld2ch