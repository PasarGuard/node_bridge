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

func CreateProxies(vmess *Vmess, vless *Vless, trojan *Trojan, shadowsocks *Shadowsocks) *Proxy {
	return &Proxy{Vmess: vmess, Vless: vless, Trojan: trojan, Shadowsocks: shadowsocks}
}

func CreateUser(email string, proxies *Proxy, inbounds []string) *User {
	return &User{Email: email, Proxies: proxies, Inbounds: inbounds}
}
