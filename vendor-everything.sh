for SERVER in clientserver infoserver robotserver switchserver userserver usecase
do
    echo $SERVER
    cd $SERVER
    go mod vendor
    cd ..
done
