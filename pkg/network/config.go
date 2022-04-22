package network

// Config is the metrics server configuration
type Config struct {
	// CNIPath is the configured CNI path
	CNIPath string
	// NetworkLabel is the label used for autoconnecting containers to networks
	NetworkLabel string
	// GRPCAddress is the address of the GRPC server
	GRPCAddress string
	// DsUri is the datastore URI
	DsURI string
}