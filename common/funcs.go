package common

func CreateVmess(id string) *Vmess {
	return &Vmess{Id: id}
}

func CreateVless(id, flow string) *Vless {
	return &Vless{Id: id, Flow: flow}
}

func CreateTrojan(password string) *Trojan {
	return &Trojan{Password: password}
}

func CreateShadowsocks(password, method string) *Shadowsocks {
	return &Shadowsocks{Password: password, Method: method}
}

func CreateWireguard(publicKey string, peerIps []string) *Wireguard {
	return &Wireguard{PublicKey: publicKey, PeerIps: peerIps}
}

func CreateHysteria(auth string) *Hysteria {
	return &Hysteria{Auth: auth}
}

func CreateMtproto(secret, userAdTag string, maxTcpConns, maxUniqueIps uint32) *Mtproto {
	return &Mtproto{Secret: secret, UserAdTag: userAdTag, MaxTcpConns: maxTcpConns, MaxUniqueIps: maxUniqueIps}
}

func CreateProxies(vmess *Vmess, vless *Vless, trojan *Trojan, shadowsocks *Shadowsocks, wireguard *Wireguard, hysteria *Hysteria, mtproto *Mtproto) *Proxy {
	return &Proxy{Vmess: vmess, Vless: vless, Trojan: trojan, Shadowsocks: shadowsocks, Wireguard: wireguard, Hysteria: hysteria, Mtproto: mtproto}
}

func CreateUser(email string, proxies *Proxy, inbounds []string) *User {
	return &User{Email: email, Proxies: proxies, Inbounds: inbounds}
}
