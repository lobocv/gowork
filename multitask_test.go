package work

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that multiple tasks run successfully
func TestMultiTask(t *testing.T) {
	ctx := context.Background()
	j := job{}
	names := []string{
		"calvin",
		"pk",
		"anthony",
		"abdul",
		"nirav",
	}
	for ii, name := range names {
		m := MultiTask{}
		p := param{name: name}
		m.AddTask(ctx, j.NoParamFunc)
		m.AddTask(ctx, func(ctx context.Context) error {
			return j.ParamFunc(ctx, &p)
		})
		errs := m.Run()

		require.Equal(t, j.NoParamFuncCalls, ii+1)
		require.Equal(t, name+" processed", p.name)
		require.Empty(t, errs)
	}
}

// Test that when some of the tasks error out, the other tasks still process and the overall error is a combination
// of all the failed task errors
func TestMultiTaskError(t *testing.T) {
	ctx := context.Background()
	j := job{}
	names := []string{
		"calvin",
		"pk",
		"anthony",
		"abdul",
		"nirav",
	}
	for ii, name := range names {
		m := MultiTask{}
		p := param{name: name}
		m.AddTask(ctx, j.NoParamFunc)
		m.AddTask(ctx, func(ctx context.Context) error {
			return j.ParamFunc(ctx, &p)
		})
		m.AddTask(ctx, j.ErrorFunc)
		m.AddTask(ctx, j.ErrorFunc)
		m.AddTask(ctx, j.ErrorFunc)
		errs := m.Run()

		require.Equal(t, j.NoParamFuncCalls, ii+1)
		require.Equal(t, name+" processed", p.name)
		for _, err := range errs {
			require.EqualError(t, err, "some error")
		}
	}
}

type job struct {
	NoParamFuncCalls int
}

func (j *job) NoParamFunc(_ context.Context) error {
	j.NoParamFuncCalls++
	return nil
}

type param struct {
	name string
}

func (j *job) ParamFunc(_ context.Context, p *param) error {
	p.name = p.name + " processed"
	return nil
}

func (j *job) ErrorFunc(_ context.Context) error {
	return fmt.Errorf("some error")
}
