package branchio

import "encoding/json"

type DeepLinkProperties struct {
	Channel  string                 `json:"channel,omitempty"`
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

func (b *branchio) structToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	return mapData, nil
}

func (b *branchio) buildParameters(deepLinkProperties *DeepLinkProperties, data *DeepLinkData, customData map[string]interface{}) (map[string]interface{}, error) {
	dataMap, err := b.structToMap(data)
	if err != nil {
		return nil, err
	}

	customDataMap := b.mergeMaps(dataMap, customData)
	if deepLinkProperties == nil {
		return customDataMap, nil
	}

	deepLinkProperties.Data = customDataMap

	return b.structToMap(deepLinkProperties)
}
