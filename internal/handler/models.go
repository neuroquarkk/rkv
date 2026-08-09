package handler

type ValReq struct {
	Value string `json:"value"`
}

type GetResp struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

type InfoResp struct {
	Name string `json:"name"`
}
