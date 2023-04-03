package types

type DataToAccess struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

type ContactResponce struct {
	Response struct {
		Contacts []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			//FirstName   string `json:"first_name"`
			//LastName    string `json:"last_name"`
			EmailValues []struct {
				FieldCode string `json:"field_code"`
				Values    []struct {
					Value string `json:"value"`
					//EmailCode string `json:"field_code"`
				} `json:"values"`
			} `json:"custom_fields_values"`
		} `json:"contacts"`
	} `json:"_embedded"`
}

type ImportUniResponse struct {
	Result struct {
		Added   int `json:"added"`
		Updated int `json:"updated"`
		Deleted int `json:"deleted"`
		Invalid int `json:"invalid"`
		Log     []struct {
			Index string `json:"index"`
		} `json:"log"`
	} `json:"result"`
	Error string `json:"error"`
}

type CurrentAcc struct {
	Current int `json:"id"`
}
