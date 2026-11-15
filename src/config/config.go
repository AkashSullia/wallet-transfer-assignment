package config

import (
	"os"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// var parsingError = errors.New("Error parsing")

type ServerConfig struct {
	APIServerConfig *ServerConfigAPIServer `json:"apiserver" validate:"required,dive"`
	PostgresConfig  *ServerConfigDB        `json:"db" validate:"required,dive"`
}

type ServerConfigAPIServer struct {
	Port string `json:"port" validate:"required"`
}

type ServerConfigDB struct {
	DSN                   string `json:"dsn"`
	MigrationFilePath     string `json:"migrationfilepath"`
	ConnMaxOpen           int
	ConnMaxIdleTimeMinute time.Duration
	ConnMaxIdleConns      int32
}

func GetServerConfig(jsonconfile string) (*ServerConfig, error) {
	rawjsonbytes, err := os.ReadFile(jsonconfile)
	if err != nil {
		log.Error().Err(err).Msg("Error reading file")
		return nil, err
	}

	var severConfig ServerConfig

	err = sonic.Unmarshal(rawjsonbytes, &severConfig)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling JSON config.")
		return nil, err
	}

	//msg, err := utils.ValidateStruct(severConfig)
	//if err != nil {
	//	log.Error().Err(err).Msg(msg)
	//	return nil, err
	//}

	return &severConfig, nil
}

func ConfigureZeroLog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
}
