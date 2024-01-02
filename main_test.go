package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var mayBeTapped = Rule{
	Name: "only untapped cards may be tapped",
	Condition: func(state State, card Card) bool {
		return !card.Tapped
	},
	Transformation: func(state State, target Card) State {
		var updated State = state
		for zone, cards := range updated.Board[target.Owner] {
			for idx, c := range cards {
				if c.ID == target.ID {
					updated.Board[target.Owner][zone][idx].Tapped = true
				}
			}
		}
		return updated
	},
}

func TestApply(t *testing.T) {
	type args struct {
		state   State
		actions []Action
	}
	testCases := []struct {
		desc    string
		args    args
		want    State
		wantErr error
	}{
		{
			desc: "should eval condition as true and transform state",
			args: args{
				state: State{
					map[string]map[string][]Card{
						"player_foo": {
							"zone_bar": []Card{
								{
									ID:     "card_biz_001",
									Name:   "card_biz",
									Owner:  "player_foo",
									Tapped: false,
								},
							},
						},
					},
				},
				actions: []Action{{
					Rule:     mayBeTapped,
					TargetID: "card_biz_001",
					Player:   "player_foo",
					Zone:     "zone_bar",
					Card: Card{
						ID:     "card_biz_001",
						Name:   "card_biz",
						Owner:  "player_foo",
						Tapped: false,
					},
				}},
			},
			want: State{
				map[string]map[string][]Card{
					"player_foo": {
						"zone_bar": []Card{
							{
								ID:     "card_biz_001",
								Name:   "card_biz",
								Owner:  "player_foo",
								Tapped: true,
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got, err := Apply(tC.args.state, tC.args.actions); err != nil && tC.wantErr == nil {
				t.Errorf("got %+v - wanted %+v", got, tC.wantErr)
			} else {
				if diff := cmp.Diff(got, tC.want); diff != "" {
					t.Errorf("got: %+v\n - diff: %s", got, diff)
				}
			}
		})
	}
}

func TestAnalyze(t *testing.T) {
	type args struct {
		state State
		rules []Rule
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
								ID:     "biz_123",
								Name:   "biz",
								Tapped: false,
								Owner:  "player_foo",
							},
						},
					},
				}},
				rules: []Rule{mayBeTapped},
			},
			want: Analysis{
				ValidActions: []Action{{
					TargetID: "biz_123",
					Card: Card{
						Name:   "biz",
						Tapped: false,
						ID:     "biz_123",
						Owner:  "player_foo",
					},
					Rule:   mayBeTapped,
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
				if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(Rule{}, "Condition", "Transformation")); diff != "" {
					t.Errorf("%s", diff)
				}
			}
		})
	}
}
