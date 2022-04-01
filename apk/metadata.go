// Copyright (C) 2022 Amaan Qureshi (aq0527@pm.me)
//
//
// This file is a part of APK Updater.
//
//
// This project, APK Updater, is not to be redistributed or copied without
//
// the express permission of the copyright holder, Amaan Qureshi (amaanq).

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
