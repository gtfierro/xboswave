for container in xboswave-demo-setup-waved xboswave-demo-setup-wavemq xboswave-demo-setup-influxdb xboswave-demo-setup-ingester xbos-demo-driver-system-monitor ; do
    docker kill $container
    docker rm $container
done
