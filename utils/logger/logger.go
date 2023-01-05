package logger

import (
	// Go Native Packages
	"errors"
	"fmt"
	"net"
	"os"

	// Vendor Packages
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//ZSLogger Zap Sugared Logger
var ZSLogger *zap.SugaredLogger
var ZLogger *zap.Logger

func init() {
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	ip, err := getExternalIP()
	if err != nil {
		fmt.Println(err)
	}

	ZSLogger = logger.Sugar().With("host_name", hostName, "ip", ip)
	ZLogger = logger
}

func getExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
