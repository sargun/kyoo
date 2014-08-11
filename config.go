import "encoding/json"
import "flag"


var String* config
func init()
{
	config = flag.String("config", "", "Configuration file location")
	flag.Parse()
	if *config == "" {
		panic("Config file location not passed")
	}
}