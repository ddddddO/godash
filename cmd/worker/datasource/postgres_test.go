package datasource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: test case 増やす
func TestPostgreSQL_Parse(t *testing.T) {
	tests := []struct {
		query           string
		wantParsedQuery *parsedQuery
		wantErr         error
	}{
		{
			query: "delete * from test1",
			wantParsedQuery: &parsedQuery{
				qType: deleteType,
				value: "delete * from test1",
			},
			wantErr: nil,
		},
	}

	opt := cmp.AllowUnexported(parsedQuery{})

	for _, tt := range tests {
		p := NewPostgreSQL()

		gotErr := p.Parse(tt.query)
		if gotErr != tt.wantErr {
			t.Errorf("\ngotErr: \n%v\nwantErr: \n%v", gotErr, tt.wantErr)
		}
		if diff := cmp.Diff(p.parsedQuery, tt.wantParsedQuery, opt); diff != "" {
			t.Errorf("parsedQuery is mismatch (-gotParsedQuery +wantParsedQuery):\n%s", diff)
		}
	}
}
