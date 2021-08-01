// Package config defines the key names for config items
package config

// Config keys for overrides
var (
	OverrideConfigFile = "override.configfile"
)

// Config keys for the database
var (
	DatabaseName     = "database.name"
	DatabaseHost     = "database.host"
	DatabaseType     = "database.type"
	DatabasePort     = "database.port"
	DatabaseUser     = "database.user"
	DatabasePassword = "database.password"
	DatabaseKeyFile  = "database.keyfile"
)

// Config keys for mqtt
var (
	BrokerAddress        = "broker.address"
	BrokerPort           = "broker.port"
	BrokerSSL            = "broker.ssl"
	BrokerCAPath         = "broker.capath"
	BrokerPublicKeyPath  = "broker.pubkeypath"
	BrokerPrivateKeyPath = "broker.privkeypath"
)

// Logging configuration for logrus
var (
	LoggingLevel   = "logger.level"
	LoggingFormat  = "logger.formatter"
	LoggingFile    = "logger.filepath"
	LoggingUseFile = "logger.uselogfile"
)

// Config keys for webserver
var (
	WebServerAddress = "webserver.address"
	WebServerPort    = "webserver.port"
	WebServerCache   = "webserver.cache"
	WebServerFiles   = "webserver.files"
)
