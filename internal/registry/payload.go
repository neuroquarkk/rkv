package registry

type MemberReq struct {
	Address string `json:"address"`
}

type MembersResp struct {
	Members []string `json:"members"`
}
