package masjid_info

import "testing"

/*
name: "Centre Culturel Islamique de Laval Khalid Bin Walid"
address: "1330 Antonio, Laval QC"
coordinates:

	latitude: 45.547559
	longitude: -73.7568045

contributed_by: "@ccil_kbw"
*/
func TestGetMasjidInfoFromFile(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		folder   string
		expected MasjidInfo
	}{
		{
			name:   "Test GetMasjidInfoFromFile",
			folder: "Chomedey Laval QC/@ccil_kbw",
			expected: MasjidInfo{
				Name:          "Centre Culturel Islamique de Laval Khalid Bin Walid",
				Website:       "https://ccil-kbw.com/",
				Address:       "1330 Antonio, Laval QC",
				Coordinates:   Coordinates{Latitude: 45.547559, Longitude: -73.7568045},
				ContributedBy: "@ccil_kbw",
			},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := GetMasjidInfoFromFile(test.folder)
			if actual.Name != test.expected.Name {
				t.Errorf("Expected Name: %s, Got: %s", test.expected.Name, actual.Name)
			}
			if actual.Address != test.expected.Address {
				t.Errorf("Expected Address: %s, Got: %s", test.expected.Address, actual.Address)
			}
			if actual.Coordinates.Latitude != test.expected.Coordinates.Latitude {
				t.Errorf("Expected Latitude: %f, Got: %f", test.expected.Coordinates.Latitude, actual.Coordinates.Latitude)
			}
			if actual.Coordinates.Longitude != test.expected.Coordinates.Longitude {
				t.Errorf("Expected Longitude: %f, Got: %f", test.expected.Coordinates.Longitude, actual.Coordinates.Longitude)
			}
			if actual.ContributedBy != test.expected.ContributedBy {
				t.Errorf("Expected ContributedBy: %s, Got: %s", test.expected.ContributedBy, actual.ContributedBy)
			}
		})
	}
}
