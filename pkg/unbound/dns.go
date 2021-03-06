package unbound

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/dns/pkg/models"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

type Configurator interface {
	SetupFiles(ctx context.Context) error
	MakeUnboundConf(settings models.Settings,
		hostnamesLines, ipsLines []string, username string,
		puid, pgid int) (err error)
	UseDNSInternally(IP net.IP)
	UseDNSSystemWide(ip net.IP, keepNameserver bool) error
	Start(ctx context.Context, verbosityDetailsLevel uint8) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
	WaitForUnbound(ctx context.Context) (err error)
	Version(ctx context.Context) (version string, err error)
	BuildBlocked(ctx context.Context, client *http.Client,
		blockMalicious, blockAds, blockSurveillance bool,
		blockedHostnames, blockedIPs, allowedHostnames []string) (
		hostnamesLines, ipsLines []string, errs []error)
}

type configurator struct {
	openFile      os.OpenFileFunc
	commander     command.Commander
	resolver      *net.Resolver
	dnscrypto     dnscrypto.DNSCrypto
	unboundEtcDir string
	unboundPath   string
	cacertsPath   string
}

func NewConfigurator(logger logging.Logger, openFile os.OpenFileFunc,
	dnscrypto dnscrypto.DNSCrypto, unboundEtcDir, unboundPath, cacertsPath string) Configurator {
	return &configurator{
		openFile:      openFile,
		commander:     command.NewCommander(),
		resolver:      net.DefaultResolver,
		dnscrypto:     dnscrypto,
		unboundEtcDir: unboundEtcDir,
		unboundPath:   unboundPath,
		cacertsPath:   cacertsPath,
	}
}
