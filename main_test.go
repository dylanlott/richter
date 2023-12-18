package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestAnalyze(t *testing.T) {
	type args struct {
		state State
		rules []Rule
	}
	var rule = Rule{
		Condition: func(state State, card Card) bool {
			return !card.Tapped
		},
		Transformation: nil,
	}
	tests := []struct {
		name string
		args args
		want Analysis
	}{
		{
			name: "should analyze that biz can be tapped",
			args: args{
				state: State{Board: map[string]map[string][]Card{
					"foo": {
						"bar": []Card{
							{
								Name:   "biz",
								Tapped: false,
							},
						},
					},
				}},
				rules: []Rule{rule},
			},
			want: Analysis{
				ValidActions: []Action{{
					Target: Card{
						Name:   "biz",
						Tapped: false,
					},
					Rule:   rule,
					Player: "foo",
					Zone:   "bar",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := Analyze(tt.args.state, tt.args.rules)
			for idx, want := range tt.want.ValidActions {
				got := analysis.ValidActions[idx]
				if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(Rule{}, "Condition")); diff != "" {
					t.Errorf("%s", diff)
				}
			}
		})
	}
}
