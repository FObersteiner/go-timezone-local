//go:build windows
// +build windows

package timezone

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"path"
	"time"
)

const winTZtoIANA = map[string]string{
	"AUS Central Standard Time":       "Australia/Darwin",
	"AUS Eastern Standard Time":       "Australia/Sydney",
	"Afghanistan Standard Time":       "Asia/Kabul",
	"Alaskan Standard Time":           "America/Anchorage",
	"Aleutian Standard Time":          "America/Adak",
	"Altai Standard Time":             "Asia/Barnaul",
	"Arab Standard Time":              "Asia/Riyadh",
	"Arabian Standard Time":           "Asia/Dubai",
	"Arabic Standard Time":            "Asia/Baghdad",
	"Argentina Standard Time":         "America/Buenos_Aires",
	"Astrakhan Standard Time":         "Europe/Astrakhan",
	"Atlantic Standard Time":          "America/Halifax",
	"Aus Central W. Standard Time":    "Australia/Eucla",
	"Azerbaijan Standard Time":        "Asia/Baku",
	"Azores Standard Time":            "Atlantic/Azores",
	"Bahia Standard Time":             "America/Bahia",
	"Bangladesh Standard Time":        "Asia/Dhaka",
	"Belarus Standard Time":           "Europe/Minsk",
	"Bougainville Standard Time":      "Pacific/Bougainville",
	"Canada Central Standard Time":    "America/Regina",
	"Cape Verde Standard Time":        "Atlantic/Cape_Verde",
	"Caucasus Standard Time":          "Asia/Yerevan",
	"Cen. Australia Standard Time":    "Australia/Adelaide",
	"Central America Standard Time":   "America/Guatemala",
	"Central Asia Standard Time":      "Asia/Almaty",
	"Central Brazilian Standard Time": "America/Cuiaba",
	"Central Europe Standard Time":    "Europe/Budapest",
	"Central European Standard Time":  "Europe/Warsaw",
	"Central Pacific Standard Time":   "Pacific/Guadalcanal",
	"Central Standard Time":           "America/Chicago",
	"Central Standard Time (Mexico)":  "America/Mexico_City",
	"Chatham Islands Standard Time":   "Pacific/Chatham",
	"China Standard Time":             "Asia/Shanghai",
	"Cuba Standard Time":              "America/Havana",
	"Dateline Standard Time":          "Etc/GMT+12",
	"E. Africa Standard Time":         "Africa/Nairobi",
	"E. Australia Standard Time":      "Australia/Brisbane",
	"E. Europe Standard Time":         "Europe/Chisinau",
	"E. South America Standard Time":  "America/Sao_Paulo",
	"Easter Island Standard Time":     "Pacific/Easter",
	"Eastern Standard Time":           "America/New_York",
	"Eastern Standard Time (Mexico)":  "America/Cancun",
	"Egypt Standard Time":             "Africa/Cairo",
	"Ekaterinburg Standard Time":      "Asia/Yekaterinburg",
	"FLE Standard Time":               "Europe/Kiev",
	"Fiji Standard Time":              "Pacific/Fiji",
	"GMT Standard Time":               "Europe/London",
	"GTB Standard Time":               "Europe/Bucharest",
	"Georgian Standard Time":          "Asia/Tbilisi",
	"Greenland Standard Time":         "America/Godthab",
	"Greenwich Standard Time":         "Atlantic/Reykjavik",
	"Haiti Standard Time":             "America/Port-au-Prince",
	"Hawaiian Standard Time":          "Pacific/Honolulu",
	"India Standard Time":             "Asia/Calcutta",
	"Iran Standard Time":              "Asia/Tehran",
	"Israel Standard Time":            "Asia/Jerusalem",
	"Jordan Standard Time":            "Asia/Amman",
	"Kaliningrad Standard Time":       "Europe/Kaliningrad",
	"Korea Standard Time":             "Asia/Seoul",
	"Libya Standard Time":             "Africa/Tripoli",
	"Line Islands Standard Time":      "Pacific/Kiritimati",
	"Lord Howe Standard Time":         "Australia/Lord_Howe",
	"Magadan Standard Time":           "Asia/Magadan",
	"Magallanes Standard Time":        "America/Punta_Arenas",
	"Marquesas Standard Time":         "Pacific/Marquesas",
	"Mauritius Standard Time":         "Indian/Mauritius",
	"Middle East Standard Time":       "Asia/Beirut",
	"Montevideo Standard Time":        "America/Montevideo",
	"Morocco Standard Time":           "Africa/Casablanca",
	"Mountain Standard Time":          "America/Denver",
	"Mountain Standard Time (Mexico)": "America/Chihuahua",
	"Myanmar Standard Time":           "Asia/Rangoon",
	"N. Central Asia Standard Time":   "Asia/Novosibirsk",
	"Namibia Standard Time":           "Africa/Windhoek",
	"Nepal Standard Time":             "Asia/Katmandu",
	"New Zealand Standard Time":       "Pacific/Auckland",
	"Newfoundland Standard Time":      "America/St_Johns",
	"Norfolk Standard Time":           "Pacific/Norfolk",
	"North Asia East Standard Time":   "Asia/Irkutsk",
	"North Asia Standard Time":        "Asia/Krasnoyarsk",
	"North Korea Standard Time":       "Asia/Pyongyang",
	"Omsk Standard Time":              "Asia/Omsk",
	"Pacific SA Standard Time":        "America/Santiago",
	"Pacific Standard Time":           "America/Los_Angeles",
	"Pacific Standard Time (Mexico)":  "America/Tijuana",
	"Pakistan Standard Time":          "Asia/Karachi",
	"Paraguay Standard Time":          "America/Asuncion",
	"Qyzylorda Standard Time":         "Asia/Qyzylorda",
	"Romance Standard Time":           "Europe/Paris",
	"Russia Time Zone 10":             "Asia/Srednekolymsk",
	"Russia Time Zone 11":             "Asia/Kamchatka",
	"Russia Time Zone 3":              "Europe/Samara",
	"Russian Standard Time":           "Europe/Moscow",
	"SA Eastern Standard Time":        "America/Cayenne",
	"SA Pacific Standard Time":        "America/Bogota",
	"SA Western Standard Time":        "America/La_Paz",
	"SE Asia Standard Time":           "Asia/Bangkok",
	"Saint Pierre Standard Time":      "America/Miquelon",
	"Sakhalin Standard Time":          "Asia/Sakhalin",
	"Samoa Standard Time":             "Pacific/Apia",
	"Sao Tome Standard Time":          "Africa/Sao_Tome",
	"Saratov Standard Time":           "Europe/Saratov",
	"Singapore Standard Time":         "Asia/Singapore",
	"South Africa Standard Time":      "Africa/Johannesburg",
	"South Sudan Standard Time":       "Africa/Juba",
	"Sri Lanka Standard Time":         "Asia/Colombo",
	"Sudan Standard Time":             "Africa/Khartoum",
	"Syria Standard Time":             "Asia/Damascus",
	"Taipei Standard Time":            "Asia/Taipei",
	"Tasmania Standard Time":          "Australia/Hobart",
	"Tocantins Standard Time":         "America/Araguaina",
	"Tokyo Standard Time":             "Asia/Tokyo",
	"Tomsk Standard Time":             "Asia/Tomsk",
	"Tonga Standard Time":             "Pacific/Tongatapu",
	"Transbaikal Standard Time":       "Asia/Chita",
	"Turkey Standard Time":            "Europe/Istanbul",
	"Turks And Caicos Standard Time":  "America/Grand_Turk",
	"US Eastern Standard Time":        "America/Indianapolis",
	"US Mountain Standard Time":       "America/Phoenix",
	"UTC":                             "UTC",
	"UTC+12":                          "Etc/GMT-12",
	"UTC+13":                          "Etc/GMT-13",
	"UTC-02":                          "Etc/GMT+2",
	"UTC-08":                          "Etc/GMT+8",
	"UTC-09":                          "Etc/GMT+9",
	"UTC-11":                          "Etc/GMT+11",
	"Ulaanbaatar Standard Time":       "Asia/Ulaanbaatar",
	"Venezuela Standard Time":         "America/Caracas",
	"Vladivostok Standard Time":       "Asia/Vladivostok",
	"Volgograd Standard Time":         "Europe/Volgograd",
	"W. Australia Standard Time":      "Australia/Perth",
	"W. Central Africa Standard Time": "Africa/Lagos",
	"W. Europe Standard Time":         "Europe/Berlin",
	"W. Mongolia Standard Time":       "Asia/Hovd",
	"West Asia Standard Time":         "Asia/Tashkent",
	"West Bank Standard Time":         "Asia/Hebron",
	"West Pacific Standard Time":      "Pacific/Port_Moresby",
	"Yakutsk Standard Time":           "Asia/Yakutsk",
	"Yukon Standard Time":             "America/Whitehorse",
}

func getKey(key, name string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, key, registry.QUERY_VALUE)
	defer k.Close()
	if err != nil {
		return "", err
	}
	val, _, err := k.GetStringValue(name)
	if err != nil {
		return "", err
	}
	return val, err
}

// LocalTZ will run `/etc/localtime` and get the timezone from the resulting value `/usr/share/zoneinfo/America/New_York`
func LocalTZ() (string, error) {
	winTZName, err := getKey(`SYSTEM\CurrentControlSet\Control\TimeZoneInformation`, "TimeZoneKeyName")
	if IANATZname, ok := winTZtoIANA[winTZName]; ok {
		return IANATZname, nil
	}
	return "", fmt.Errorf("could not find IANA time zone name for %v", winTZName)
}
