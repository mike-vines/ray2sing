package ray2sing

import (
	T "github.com/sagernet/sing-box/option"
)

func Hysteria2Singbox(hysteria2Url string) (*T.Outbound, error) {
	u, err := ParseUrl(hysteria2Url)
	if err != nil {
		return nil, err
	}
	decoded := u.Params
	var ObfsOpts *T.Hysteria2Obfs
	ObfsOpts = nil
	if obfs, ok := decoded["obfs"]; ok && obfs != "" {
		ObfsOpts = &T.Hysteria2Obfs{
			Type:     obfs,
			Password: decoded["obfs-password"],
		}
	}

	valECH, hasECH := decoded["ech"]
	hasECH = hasECH && (valECH != "0")
	var ECHOpts *T.OutboundECHOptions
	ECHOpts = nil
	if hasECH {
		ECHOpts = &T.OutboundECHOptions{
			Enabled: hasECH,
		}
	}

	SNI := decoded["sni"]
	if SNI == "" {
		SNI = decoded["hostname"]
	}
	turnRelay, err := u.GetRelayOptions()
	if err != nil {
		return nil, err
	}
	result := T.Outbound{
		Type: "hysteria2",
		Tag:  u.Name,
		Hysteria2Options: T.Hysteria2OutboundOptions{
			ServerOptions: u.GetServerOption(),
			Obfs:          ObfsOpts,
			Password:      u.Username,
			OutboundTLSOptionsContainer: T.OutboundTLSOptionsContainer{
				TLS: &T.OutboundTLSOptions{
					Enabled:    true,
					Insecure:   decoded["insecure"] == "1",
					DisableSNI: isIPOnly(SNI),
					ServerName: SNI,
					ECH:        ECHOpts,
				},
			},
			TurnRelay: turnRelay,
		},
	}

	return &result, nil
}
