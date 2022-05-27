package branchio

type DeepLinkProperties struct {
	Channel  string                 `json:"channel"`
	Feature  string                 `json:"feature,omitempty"`
	Stage    string                 `json:"stage,omitempty"`
	Campaign string                 `json:"campaign,omitempty"`
	Alias    string                 `json:"alias,omitempty"`
	Type     int                    `json:"type,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Identity string                 `json:"identity,omitempty"`
	Data     map[string]interface{} `json:"data"`
}

type DeepLinkData struct {
	CanonicalIdentifier string `json:"$canonical_identifier,omitempty"`
	OgTitle             string `json:"$og_title,omitempty"`
	OgDescription       string `json:"$og_description,omitempty"`
	OgImageUrl          string `json:"$og_image_url,omitempty"`
	DesktopUrl          string `json:"$desktop_url,omitempty"`
}

func (b *branchio) mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	mergedMaps := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			mergedMaps[k] = v
		}
	}
	return mergedMaps
}
