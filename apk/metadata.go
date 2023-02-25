/*
The GPLv3 License (GPLv3)

Copyright (c) 2023 Amaan Qureshi <amaanq12@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package apk

type MetaData struct {
	Context    string     `json:"@context"`
	Type       string     `json:"@type"`
	URL        string     `json:"url"`
	IsPartOf   IsPartOf   `json:"isPartOf"`
	MainEntity MainEntity `json:"mainEntity"`
}

type IsPartOf struct {
	Type            string          `json:"@type"`
	URL             string          `json:"url"`
	Publisher       Publisher       `json:"publisher"`
	PotentialAction PotentialAction `json:"potentialAction"`
}

type PotentialAction struct {
	Type       string `json:"@type"`
	Target     string `json:"target"`
	QueryInput string `json:"query-input"`
}

type Publisher struct {
	ID string `json:"@id"`
}

type MainEntity struct {
	Type                   string               `json:"@type"`
	Name                   string               `json:"name"`
	URL                    string               `json:"url"`
	Description            string               `json:"description"`
	Image                  string               `json:"image"`
	OperatingSystem        string               `json:"operatingSystem"`
	SoftwareVersion        string               `json:"softwareVersion"` // This is the one
	DatePublished          string               `json:"datePublished"`
	InteractionStatistic   InteractionStatistic `json:"interactionStatistic"`
	ApplicationCategory    string               `json:"applicationCategory"`
	ApplicationSubCategory string               `json:"applicationSubCategory"`
	Author                 Author               `json:"author"`
	Offers                 Offers               `json:"offers"`
	AggregateRating        AggregateRating      `json:"aggregateRating"`
	Screenshot             []Screenshot         `json:"screenshot"`
	InLanguage             []InLanguage         `json:"inLanguage"`
}

type AggregateRating struct {
	Type        string `json:"@type"`
	RatingValue string `json:"ratingValue"`
	RatingCount string `json:"ratingCount"`
	BestRating  string `json:"bestRating"`
	WorstRating string `json:"worstRating"`
}

type Author struct {
	Type             string           `json:"@type"`
	Name             string           `json:"name"`
	MainEntityOfPage MainEntityOfPage `json:"mainEntityOfPage"`
	URL              string           `json:"url"`
}

type MainEntityOfPage struct {
	Type string `json:"@type"`
	ID   string `json:"@id"`
}

type InLanguage struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

type InteractionStatistic struct {
	Type                 string `json:"@type"`
	InteractionType      string `json:"interactionType"`
	UserInteractionCount string `json:"userInteractionCount"`
}

type Offers struct {
	Type          string `json:"@type"`
	Price         string `json:"price"`
	PriceCurrency string `json:"priceCurrency"`
}

type Screenshot struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
}
