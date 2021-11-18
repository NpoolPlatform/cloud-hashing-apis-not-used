FROM jrei/systemd-centos:7

COPY CloudHashingApis.viper.yaml /etc/CloudHashingApis/CloudHashingApis.viper.yaml
COPY cloud-hashing-apis /usr/local/bin/cloud-hashing-apis

CMD ["cloud-hashing-apis", "run"]
