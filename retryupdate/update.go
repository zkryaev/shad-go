//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(
	c kvapi.Client,
	key string,
	updateFn func(oldValue *string) (newValue string, err error)) error {
LOOP:
	var req kvapi.SetRequest
	for {
		var respValue *string
		var respVersion uuid.UUID
		for {
			resp, err := c.Get(&kvapi.GetRequest{Key: key})
			var errAuth *kvapi.AuthError
			switch {
			case errors.Is(err, nil):
				respValue = &(resp.Value)
				respVersion = resp.Version
				goto UPDATE
			case errors.Is(err, kvapi.ErrKeyNotFound):
				respValue = nil
				respVersion = uuid.UUID{}
				goto UPDATE
			case errors.As(err, &errAuth):
				return err
			}
		}
	UPDATE:
		nValue, err := updateFn(respValue)
		if err != nil {
			return err
		}
		req.Key = key
		req.Value = nValue
		req.OldVersion = respVersion
		req.NewVersion = uuid.Must(uuid.NewV4())
		for {
			_, err := c.Set(&req)
			if err == nil {
				return nil
			}
			var errAuth *kvapi.AuthError
			var errConf *kvapi.ConflictError
			switch {
			case errors.Is(err, kvapi.ErrKeyNotFound):
				_, err = c.Set(&kvapi.SetRequest{
					Key:        key,
					Value:      "V0",
					OldVersion: uuid.UUID{},
					NewVersion: uuid.Must(uuid.NewV4()),
				})
				return err
			case errors.As(err, &errAuth):
				return err
			case errors.As(err, &errConf):
				if errConf.ExpectedVersion != req.NewVersion {
					goto LOOP
				}
				return nil
			}
		}
	}
}
