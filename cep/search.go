package cep

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"time"
)

var (
	handlers = []HandlerFunc{
		searchCorreios,
		searchRepublicaVirtual,
		searchPostmon,
		searchViaCEP,
	}
)

func Search(cep string) (CEP, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return SearchWithContext(ctx, cancel, cep)
}

func SearchWithContext(ctx context.Context, cancel context.CancelFunc, cep string) (CEP, error) {
	defer cancel()

	var err error
	cep, err = formatCEP(cep)
	if err != nil {
		return CEP{}, err
	}
	return searchCorreios(ctx, cep), nil
}

func SearchPublicBase(cep string) (CEP, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return SearchPublicBaseWithContext(ctx, cancel, cep)
}

func SearchPublicBaseWithContext(ctx context.Context, cancel context.CancelFunc, cep string) (CEP, error) {
	defer cancel()

	var err error
	cep, err = formatCEP(cep)
	if err != nil {
		return CEP{}, err
	}

	var (
		cepPublic    CEP
		correiosDone bool
		monitor      sync.WaitGroup
		result       = make(chan CEP, len(handlers))
	)

	for _, fn := range handlers {
		fn := fn
		monitor.Add(1)
		go func() {
			defer monitor.Done()
			result <- fn(ctx, cep)
		}()
	}

	go func() {
		monitor.Wait()
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return cepPublic, nil
		case r := <-result:
			if r.FromCorreios() {
				if r.Valid() {
					return r, nil
				}
				correiosDone = true
			}
			if !cepPublic.Valid() {
				cepPublic = r
			}
			if cepPublic.Valid() && correiosDone {
				return cepPublic, nil
			}
		case <-time.After(4 * time.Second):
			return cepPublic, nil
		}
	}
}

func formatCEP(cep string) (string, error) {
	cep = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cep, "")
	if len(cep) != 8 {
		return "", errors.New("CEP inválido.")
	}
	return cep, nil
}
