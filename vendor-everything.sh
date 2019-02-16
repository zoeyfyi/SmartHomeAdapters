for SERVER in clientserver infoserver robotserver switchserver thermostatserver userserver
do
	cd $SERVER
    go mod vendor
    cd ..
done