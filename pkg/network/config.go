package network

// Config is the metrics server configuration
type Config struct {
	// NodeName is the name of the node when using clustering
	NodeName string
	// ContainerdAddr is the containerd address
	ContainerdAddr string
	// ContainerdNamespace is the containerd namespace to manage
	ContainerdNamespace string
	// CNIPath is the configured CNI path
	CNIPath string
	// NetworkLabel is the label used for autoconnecting containers to networks
	NetworkLabel string
	// GRPCAddress is the address of the GRPC server
	GRPCAddress string
	// DsUri is the datastore URI
	DsURI string
	// TLSCertificate is the certificate used for grpc communication
	TLSServerCertificate string
	// TLSKey is the key used for grpc communication
	TLSServerKey string
	// TLSInsecureSkipVerify disables certificate verification
	TLSInsecureSkipVerify bool
	// RedisURL is the Redis address for clustering
	RedisURL string
}