package client_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vx416/bitfinex-api-go/pkg/models/event"
	"github.com/vx416/bitfinex-api-go/pkg/mux/client"
)

func TestSubsLimitReached(t *testing.T) {
	cases := map[string]struct {
		limit    int
		expected bool
		subs     []event.Subscribe
	}{
		"limit unreached": {
			limit:    20,
			expected: false,
			subs:     []event.Subscribe{{Event: "foo"}},
		},
		"limit reached": {
			limit:    2,
			expected: true,
			subs: []event.Subscribe{
				{Event: "foo"},
				{Event: "bar"},
			},
		},
		"limit unreached #2": {
			limit:    0,
			expected: false,
			subs: []event.Subscribe{
				{Event: "foo"},
				{Event: "bar"},
			},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			c := client.New().WithSubsLimit(v.limit)
			for _, e := range v.subs {
				c.AddSub(e)
			}

			got := c.SubsLimitReached()
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestSubAdded(t *testing.T) {
	cases := map[string]struct {
		expected bool
		subs     []event.Subscribe
		pld      event.Subscribe
	}{
		"not added": {
			expected: false,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "bar"},
		},
		"added": {
			expected: true,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "foo"},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			c := client.New()
			for _, e := range v.subs {
				c.AddSub(e)
			}

			got := c.SubAdded(v.pld)
			assert.Equal(t, v.expected, got)
		})
	}
}

func TestRemoveSub(t *testing.T) {
	cases := map[string]struct {
		expected bool
		subs     []event.Subscribe
		pld      event.Subscribe
	}{
		"removing existing sub": {
			expected: false,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "foo"},
		},
		"removing unexisting sub": {
			expected: true,
			subs:     []event.Subscribe{{Event: "foo"}},
			pld:      event.Subscribe{Event: "bar"},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			c := client.New()
			for _, e := range v.subs {
				c.AddSub(e)
			}

			c.RemoveSub(v.pld)
			got := c.SubAdded(v.subs[0])
			assert.Equal(t, v.expected, got)
		})
	}
}

type byEvent []event.Subscribe

func (x byEvent) Len() int {
	return len(x)
}

func (x byEvent) Less(i, j int) bool {
	return x[i].Event < x[j].Event
}

func (x byEvent) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func TestGetAllSubs(t *testing.T) {
	cases := map[string]struct {
		expected []event.Subscribe
		subs     []event.Subscribe
	}{
		"get all subs": {
			expected: []event.Subscribe{{Event: "bar"}, {Event: "foo"}},
			subs:     []event.Subscribe{{Event: "bar"}, {Event: "foo"}},
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			c := client.New()
			for _, e := range v.subs {
				c.AddSub(e)
			}

			got := c.GetAllSubs()
			sort.Sort(byEvent(got))
			assert.Equal(t, v.expected, got)
		})
	}
}
