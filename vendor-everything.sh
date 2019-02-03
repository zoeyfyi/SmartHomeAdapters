for SERVER in clientserver infoserver robotserver switchserver userserver
do
	cd $SERVER
    go mod vendor
    cd ..
done