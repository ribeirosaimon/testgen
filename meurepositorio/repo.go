package meurepositorio

import (
	"context"
	"time"
)

type AnotherInterface interface {
	ComputeSum(ctx context.Context, a, b int) (int, error)
	IsActive(ctx context.Context, user string) bool
	FetchRandom(ctx context.Context, id int) (int, int)
	GetScores(ctx context.Context, ids ...string) ([]float64, error)
	GetTimestamps(ctx context.Context, tags []string) (map[string]time.Time, error)
}

type another struct {
}

func NewAnother() AnotherInterface {
	return another{}
}
func (a2 another) ComputeSum(ctx context.Context, a, b int) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (a2 another) IsActive(ctx context.Context, user string) bool {
	//TODO implement me
	panic("implement me")
}

func (a2 another) FetchRandom(ctx context.Context, id int) (int, int) {
	//TODO implement me
	panic("implement me")
}

func (a2 another) GetScores(ctx context.Context, ids ...string) ([]float64, error) {
	//TODO implement me
	panic("implement me")
}

func (a2 another) GetTimestamps(ctx context.Context, tags []string) (map[string]time.Time, error) {
	//TODO implement me
	panic("implement me")
}
