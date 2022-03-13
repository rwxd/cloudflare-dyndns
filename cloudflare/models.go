package cloudflare

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CloudFlareResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Result interface{} `json:"result"`
}

type CloudFlareTokenVerify struct {
	Result struct {
		Id     string `json:"id"`
		Status string `json:"status"`
	}
	Sucess bool `json:"success"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	Messages []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}
