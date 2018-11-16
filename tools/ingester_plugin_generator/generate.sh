for i in `ls ../../proto/*.proto` ; do
    dest=$(basename $i .proto)_plugin
    ./ingester_plugin_generator $i $dest 
    sed -i -e '/package main/a import "strings"' $dest.go
    sed -i -e '/package main/a import "fmt"' $dest.go
    sed -i -e '/package main/a import xbospb "github.com/gtfierro/xboswave/proto"' $dest.go
    sed -i -e '/package main/a import "github.com/gtfierro/xboswave/ingester/types"' $dest.go
done
