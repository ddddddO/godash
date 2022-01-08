package datasource

import (
	"testing"
)

func TestPostgreSQL_Parse(t *testing.T) {
	tests := []struct {
		query           string
		wantParsedQuery *parsedQuery
		wantErr         error
	}{
		{
			query:   "delete * from test1",
			wantErr: errUndefinedType,
		},
	}

	for _, tt := range tests {
		p := NewPostgreSQL()

		gotErr := p.Parse(tt.query)
		if gotErr != tt.wantErr {
			t.Errorf("\ngotErr: \n%v\nwantErr: \n%v", gotErr, tt.wantErr)
		}
		if p.parsedQuery != tt.wantParsedQuery {
			t.Errorf("\ngotParsedQuery: \n%v\nwantParsedQuery: \n%v", p.parsedQuery, tt.wantParsedQuery)
		}
	}
}
