package main

import "time"

type message struct {
	Name      string   `json:"name"`
	Message   string   `json:"message"`
	When      jsonTime `json:"when"`
	AvatarURL string   `json:"avatarURL"`
}

type jsonTime struct {
	time.Time
}

func (j jsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.Format("2006-01-02 15:04:05") + `"`), nil
}

func (j *jsonTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	timeTime, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	if err != nil {
		return err
	}

	*j = jsonTime{timeTime}
	return err
}
