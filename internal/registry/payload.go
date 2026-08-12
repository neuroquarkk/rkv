package registry

type HeartbeatReq struct {
	Address string `json:"address"`
}

type MembersResp struct {
	Members []string `json:"members"`
}
