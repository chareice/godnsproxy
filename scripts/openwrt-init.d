#!/bin/sh /etc/rc.common
# OpenWRT init script for DNS Proxy
# Reference: https://openwrt.org/docs/guide-developer/procd-init-scripts

START=99
STOP=10
USE_PROCD=1

BIN="/opt/godnsproxy/godnsproxy"
CONF="/etc/godnsproxy.conf"

start_service() {
    procd_open_instance
    procd_set_param command "$BIN" \
        -f "$(uci -q get godnsproxy.@config[0].domain_file || echo '/opt/godnsproxy/domains.txt')" \
        -p "$(uci -q get godnsproxy.@config[0].port || echo '5300')" \
        -c "$(uci -q get godnsproxy.@config[0].china_dns || echo '223.5.5.5')" \
        -t "$(uci -q get godnsproxy.@config[0].trust_dns || echo 'https://1.1.1.1/dns-query')"
    procd_set_param respawn
    procd_set_param stdout 1
    procd_set_param stderr 1
    procd_close_instance
}

reload_service() {
    stop
    start
}
