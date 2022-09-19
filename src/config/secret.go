package config

import "encoding/json"

const secretPlaceHolder = "*****"

type SecretString string

func (s SecretString) String() string {
	return secretPlaceHolder
}

func (s SecretString) Get() string {
	return string(s)
}

func (s SecretString) MarshalJSON() ([]byte, error) {
	type secretString SecretString
	ss := secretString(s)
	ss = secretPlaceHolder
	return json.Marshal((*secretString)(&ss))
}
