package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"github.com/whatsauth/wa"
)

func PostBalasan(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)
	link := "https://medium.com/@prisyahaura15/tugas-gis-7-be15a493b663#id_token=eyJhbGciOiJSUzI1NiIsImtpZCI6ImY4MzNlOGE3ZmUzZmU0Yjg3ODk0ODIxOWExNjg0YWZhMzczY2E4NmYiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIyMTYyOTYwMzU4MzQtazFrNnFlMDYwczJ0cDJhMmphbTRsamRjbXMwMHN0dGcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiIyMTYyOTYwMzU4MzQtazFrNnFlMDYwczJ0cDJhMmphbTRsamRjbXMwMHN0dGcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDIwMDUyNjgzMTYyMDMxOTIzMzciLCJlbWFpbCI6InByaXN5YWhhdXJhMTVAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5iZiI6MTY5OTg0Mzc3MywibmFtZSI6IlByaXN5YSBIYXVyYSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NJNlpsYVI4R0toVGYyVlhQWGRvakxDZ0h2RW9ncjdLRnB0ek5CbjRQLU5jbHc9czk2LWMiLCJnaXZlbl9uYW1lIjoiUHJpc3lhIiwiZmFtaWx5X25hbWUiOiJIYXVyYSIsImxvY2FsZSI6ImlkIiwiaWF0IjoxNjk5ODQ0MDczLCJleHAiOjE2OTk4NDc2NzMsImp0aSI6ImY5ZTM0OWRhZmNlMThlMmQxMzMxM2UxMWE4ZTBmNmE2MjFkYWFhMmIifQ.ELt-Io8qhQIGgxyaavzwyjdMqfaCeLqXWp1Yz4uf9ry-RWWL-3oLgUhy1dqDLLur8OgHd6Gvam7omQ3uCPgUpgzv9HOmumP9Kj72jwUsFhgUDsboT2jyWoaYsiiT2t6NyaAMi6pQ3y7kP6SiEhDu0bb5uJBSuDV9dRKRqk5dKSVhaXoS57rno2-Mjtvdnm6XS2Pr9bXLZq90TkkPGTxf_ZyNQPT1qEjShM0ml--oSDY1-CtleskQOkqjDlfsgC2-GTi8ETC8L7CJfvfeThyKI6DXGSThz3YOZGV9Bul6UPjOfI6jk0lLnDCW5dC9vcr3DQooQbZItpapQFEt3QZizQ"
	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		if msg.Message == "loc" || msg.Message == "Loc" || msg.Message == "lokasi" || msg.LiveLoc {
			location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
			if err != nil {
				// Handle the error (e.g., log it) and set a default location name
				location = "Unknown Location"
			}

			reply := fmt.Sprintf("Hai hai haiii kamu pasti lagi di %s \n Koordinatenya : %s - %s\n Cara Penggunaan WhatsAuth Ada di link dibawah ini"+
				"yaa kak %s\n", location,
				strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)), link)
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: reply,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		} else {
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: "Hai Hai Haiii kamuuuu " + msg.Alias_name + "\nPrisyanya lagi gaadaa tauk \n haloo aku Prisya salam kenall yaaaa \n Cara penggunaan WhatsAuth ada di link berikut ini ya kak...\n" + link,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		}
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}

func ReverseGeocode(latitude, longitude float64) (string, error) {
	// OSM Nominatim API endpoint
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", latitude, longitude)

	// Make a GET request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	// Extract the place name from the response
	displayName, ok := result["display_name"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract display_name from the API response")
	}

	return displayName, nil
}

func Liveloc(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)

	// Reverse geocode to get the place name
	location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
	if err != nil {
		// Handle the error (e.g., log it) and set a default location name
		location = "Unknown Location"
	}

	reply := fmt.Sprintf("Hai hai haiii kamu pasti lagi di %s \n Koordinatenya : %s - %s\n", location,
		strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)))

	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		dt := &wa.TextMessage{
			To:       msg.Phone_number,
			IsGroup:  false,
			Messages: reply,
		}
		resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://cloud.wa.my.id/api/send/message/text")
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}
